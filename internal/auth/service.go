package auth

import (
	"context"
	"errors"
	"fmt"

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

type AuthResponse struct {
	Token string     `json:"token"`
	User  *user.User `json:"user"`
}

type Service interface {
	Login(ctx context.Context, req LoginRequest) (*AuthResponse, error)
	Register(ctx context.Context, req RegisterRequest) (*AuthResponse, error)
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

	token, err := s.jwtService.GenerateToken(existingUser.ID, existingUser.Email, existingUser.Role)
	if err != nil {
		return nil, errors.New("failed to generate token")
	}

	return &AuthResponse{
		Token: token,
		User:  existingUser,
	}, nil
}

func (s *service) Register(ctx context.Context, req RegisterRequest) (*AuthResponse, error) {
	existing, _ := s.userRepo.GetUserByEmail(ctx, req.Email)
	if existing != nil {
		return nil, errors.New("email already exists")
	}

	fmt.Println(bcrypt.DefaultCost, "bcrypt.DefaultCost..")

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

	token, err := s.jwtService.GenerateToken(newUser.ID, newUser.Email, newUser.Role)
	if err != nil {
		return nil, errors.New("failed to generate token")
	}

	return &AuthResponse{
		Token: token,
		User:  newUser,
	}, nil
}
