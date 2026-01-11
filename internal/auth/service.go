package auth

import (
	"context"
	"errors"

	"go-backend-project/internal/user"

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
}

type service struct {
	userRepo   user.Repository
	jwtService JWTService
}

func NewService(userRepo user.Repository, jwtService JWTService) Service {
	return &service{
		userRepo:   userRepo,
		jwtService: jwtService,
	}
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

	return &AuthResponse{
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
		ExpiresIn:    tokenPair.ExpiresIn,
		User:         existingUser,
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

	return &AuthResponse{
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
		ExpiresIn:    tokenPair.ExpiresIn,
		User:         newUser,
	}, nil
}

func (s *service) RefreshToken(ctx context.Context, req RefreshTokenRequest) (*AuthResponse, error) {
	tokenPair, err := s.jwtService.RefreshAccessToken(req.RefreshToken)
	if err != nil {
		return nil, err
	}

	claims, _ := s.jwtService.ValidateToken(tokenPair.AccessToken)
	existingUser, err := s.userRepo.GetUserByEmail(ctx, claims.Email)
	if err != nil {
		return nil, errors.New("user not found")
	}

	if !existingUser.IsActive {
		return nil, errors.New("user account is deactivated")
	}

	return &AuthResponse{
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
		ExpiresIn:    tokenPair.ExpiresIn,
		User:         existingUser,
	}, nil
}
