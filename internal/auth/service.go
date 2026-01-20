package auth

import (
	"context"
	"errors"
	"fmt"
	"time"

	"go-backend-project/internal/redis"
	"go-backend-project/internal/user"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type RegisterRequest struct {
	FirstName string `json:"first_name" binding:"required"`
	LastName  string `json:"last_name" binding:"required"`
	Email     string `json:"email" binding:"required,email"`
	Password  string `json:"password" binding:"required,min=6"`
	Phone     string `json:"phone"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type AuthResponse struct {
	AccessToken  string     `json:"access_token"`
	RefreshToken string     `json:"refresh_token"`
	ExpiresIn    int64      `json:"expires_in"`
	User         *user.User `json:"user"`
}

type Service interface {
	Login(ctx context.Context, req LoginRequest) (*AuthResponse, error)
	Register(ctx context.Context, req RegisterRequest) (*AuthResponse, error)
	RefreshToken(ctx context.Context, req RefreshTokenRequest) (*AuthResponse, error)
	GetMe(ctx context.Context, userID string) (*AuthResponse, error)
	Logout(ctx context.Context, refreshToken string) error
}

type service struct {
	userRepo     user.Repository
	jwtService   JWTService
	redisService *redis.RedisService
}

func NewService(userRepo user.Repository, jwtService JWTService, redisService *redis.RedisService) Service {
	return &service{
		userRepo:     userRepo,
		jwtService:   jwtService,
		redisService: redisService,
	}
}

func (s *service) setRefreshToken(ctx context.Context, userID, tokenID string, expiresIn time.Duration) error {
	key := fmt.Sprintf("refresh:%s", tokenID)
	return s.redisService.RedisClient.Set(ctx, key, userID, expiresIn).Err()
}

func (s *service) deleteRefreshToken(ctx context.Context, tokenID string) error {
	key := fmt.Sprintf("refresh:%s", tokenID)
	return s.redisService.RedisClient.Del(ctx, key).Err()
}

func (s *service) getRefreshTokenCache(ctx context.Context, tokenID string) (string, error) {
	key := fmt.Sprintf("refresh:%s", tokenID)
	return s.redisService.RedisClient.Get(ctx, key).Result()
}

func (s *service) Login(ctx context.Context, req LoginRequest) (*AuthResponse, error) {
	existingUser, err := s.userRepo.GetUserByEmail(ctx, req.Email)
	if err != nil {
		return nil, errors.New("invalid email or password")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(existingUser.Password), []byte(req.Password)); err != nil {
		return nil, errors.New("invalid email or password")
	}

	if !existingUser.IsActive {
		return nil, errors.New("user account is deactivated")
	}

	tokenPair, err := s.jwtService.GenerateTokenPair(existingUser.ID, existingUser.Email, existingUser.Role)
	if err != nil {
		return nil, errors.New("failed to generate tokens")
	}

	if err := s.setRefreshToken(ctx, existingUser.ID.String(), tokenPair.RefreshJTI, time.Until(tokenPair.RefreshExpiresAt)); err != nil {
		return nil, errors.New("failed to store session")
	}

	return &AuthResponse{
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
		ExpiresIn:    tokenPair.ExpiresIn,
		User:         existingUser,
	}, nil
}

func (s *service) GetMe(ctx context.Context, userId string) (*AuthResponse, error) {
	userUUID, err := uuid.Parse(userId)
	if err != nil {
		return nil, errors.New("invalid userId")
	}

	existingUser, err := s.userRepo.GetUserByID(ctx, userUUID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	if !existingUser.IsActive {
		return nil, errors.New("user account is deactivated")
	}

	return &AuthResponse{
		User: existingUser,
	}, nil
}

func (s *service) Register(ctx context.Context, req RegisterRequest) (*AuthResponse, error) {
	existing, _ := s.userRepo.GetUserByEmail(ctx, req.Email)
	if existing != nil {
		return nil, errors.New("email already exists")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.New("failed to hash password")
	}

	newUser := &user.User{
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Email:     req.Email,
		Password:  string(hashedPassword),
		Phone:     req.Phone,
		IsActive:  true,
		Role:      "customer",
	}

	if err := s.userRepo.Create(ctx, newUser); err != nil {
		return nil, errors.New("failed to create user")
	}

	tokenPair, err := s.jwtService.GenerateTokenPair(newUser.ID, newUser.Email, newUser.Role)
	if err != nil {
		return nil, errors.New("failed to generate tokens")
	}

	if err := s.setRefreshToken(ctx, newUser.ID.String(), tokenPair.RefreshJTI, time.Until(tokenPair.RefreshExpiresAt)); err != nil {
		return nil, errors.New("failed to store session")
	}

	return &AuthResponse{
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
		ExpiresIn:    tokenPair.ExpiresIn,
		User:         newUser,
	}, nil
}

func (s *service) Logout(ctx context.Context, refreshToken string) error {
	claims, err := s.jwtService.ValidateToken(refreshToken)
	if err != nil {
		return nil
	}

	if claims.JTI == nil {
		return errors.New("invalid refresh token: missing jti")
	}

	key := fmt.Sprintf("refresh:%s", *claims.JTI)
	return s.redisService.RedisClient.Del(ctx, key).Err()
}

func (s *service) RefreshToken(ctx context.Context, req RefreshTokenRequest) (*AuthResponse, error) {
	claims, err := s.jwtService.ValidateToken(req.RefreshToken)
	if err != nil {
		return nil, errors.New("invalid refresh token")
	}

	if claims.TokenType != RefreshToken {
		return nil, errors.New("invalid token type")
	}

	if claims.JTI == nil {
		return nil, errors.New("invalid refresh token: missing jti")
	}

	userID, err := s.getRefreshTokenCache(ctx, *claims.JTI)
	if err != nil {
		return nil, errors.New("refresh token revoked")
	}

	if userID != claims.UserID.String() {
		return nil, errors.New("token user mismatch")
	}

	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return nil, errors.New("invalid user id")
	}

	user, err := s.userRepo.GetUserByID(ctx, userUUID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	if !user.IsActive {
		return nil, errors.New("user account is deactivated")
	}

	newPair, err := s.jwtService.GenerateTokenPair(user.ID, user.Email, user.Role)
	if err != nil {
		return nil, err
	}

	if err := s.setRefreshToken(ctx, user.ID.String(), newPair.RefreshJTI, time.Until(newPair.RefreshExpiresAt)); err != nil {
		return nil, errors.New("failed to store session")
	}

	_ = s.deleteRefreshToken(ctx, *claims.JTI)

	return &AuthResponse{
		AccessToken:  newPair.AccessToken,
		RefreshToken: newPair.RefreshToken,
		ExpiresIn:    newPair.ExpiresIn,
		User:         user,
	}, nil
}
