package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type TokenType string

const (
	AccessToken  TokenType = "access"
	RefreshToken TokenType = "refresh"
)

type Claims struct {
	UserID    uuid.UUID `json:"user_id"`
	Email     string    `json:"email"`
	Role      string    `json:"role"`
	TokenType TokenType `json:"token_type"`
	JTI       *string   `json:"jti"`
	jwt.RegisteredClaims
}

type TokenPair struct {
	AccessToken      string    `json:"access_token"`
	RefreshToken     string    `json:"refresh_token"`
	ExpiresIn        int64     `json:"expires_in"`
	RefreshJTI       string    `json:"refresh_jti"`
	RefreshExpiresAt time.Time `json:"refresh_expires_at"`
}

type RefreshTokenResult struct {
	Token     string
	JTI       string
	ExpiresAt time.Time
}

type JWTService interface {
	GenerateTokenPair(userID uuid.UUID, email, role string) (*TokenPair, error)
	GenerateAccessToken(userID uuid.UUID, email, role string) (string, error)
	GenerateRefreshToken(userID uuid.UUID, email, role string) (string, error)
	ValidateToken(tokenString string) (*Claims, error)
	RefreshAccessToken(refreshToken string) (*TokenPair, error)
}

type jwtService struct {
	secretKey                 []byte
	accessTokenExpirationMin  int
	refreshTokenExpirationDay int
}

func NewJWTService(secret string, accessTokenExpirationMin, refreshTokenExpirationDay int) JWTService {
	if accessTokenExpirationMin <= 0 {
		accessTokenExpirationMin = 15
	}
	if refreshTokenExpirationDay <= 0 {
		refreshTokenExpirationDay = 7
	}

	return &jwtService{
		secretKey:                 []byte(secret),
		accessTokenExpirationMin:  accessTokenExpirationMin,
		refreshTokenExpirationDay: refreshTokenExpirationDay,
	}
}

func (s *jwtService) GenerateAccessToken(userID uuid.UUID, email, role string) (string, error) {
	claims := &Claims{
		UserID:    userID,
		Email:     email,
		Role:      role,
		TokenType: AccessToken,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * time.Duration(s.accessTokenExpirationMin))),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.secretKey)
}

func (s *jwtService) GenerateRefreshToken(userID uuid.UUID, email, role string) (string, error) {
	result, err := s.generateRefreshTokenWithDetails(userID, email, role)
	if err != nil {
		return "", err
	}
	return result.Token, nil
}

func (s *jwtService) generateRefreshTokenWithDetails(userID uuid.UUID, email, role string) (*RefreshTokenResult, error) {
	jti := uuid.NewString()
	expiresAt := time.Now().Add(time.Hour * 24 * time.Duration(s.refreshTokenExpirationDay))

	claims := &Claims{
		UserID:    userID,
		Email:     email,
		Role:      role,
		TokenType: RefreshToken,
		JTI:       &jti,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(s.secretKey)
	if err != nil {
		return nil, err
	}

	return &RefreshTokenResult{
		Token:     signedToken,
		JTI:       jti,
		ExpiresAt: expiresAt,
	}, nil
}

func (s *jwtService) GenerateTokenPair(userID uuid.UUID, email, role string) (*TokenPair, error) {
	accessToken, err := s.GenerateAccessToken(userID, email, role)
	if err != nil {
		return nil, err
	}

	refreshResult, err := s.generateRefreshTokenWithDetails(userID, email, role)
	if err != nil {
		return nil, err
	}

	return &TokenPair{
		AccessToken:      accessToken,
		RefreshToken:     refreshResult.Token,
		ExpiresIn:        int64(s.accessTokenExpirationMin * 60),
		RefreshJTI:       refreshResult.JTI,
		RefreshExpiresAt: refreshResult.ExpiresAt,
	}, nil
}

func (s *jwtService) ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return s.secretKey, nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}

func (s *jwtService) RefreshAccessToken(refreshToken string) (*TokenPair, error) {
	claims, err := s.ValidateToken(refreshToken)
	if err != nil {
		return nil, errors.New("invalid refresh token")
	}

	if claims.TokenType != RefreshToken {
		return nil, errors.New("invalid token type: expected refresh token")
	}

	return s.GenerateTokenPair(claims.UserID, claims.Email, claims.Role)
}
