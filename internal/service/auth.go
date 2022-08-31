package service

import (
	"context"
	"crypto/sha1"
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/rusystem/notes-app/internal/config"
	"github.com/rusystem/notes-app/internal/domain"
	"github.com/rusystem/notes-app/internal/repository"
	"math/rand"
	"time"
)

type tokenClaims struct {
	jwt.StandardClaims
	UserId int `json:"user_id"`
}

type AuthService struct {
	cfg  *config.Config
	repo repository.Authorization
}

func NewAuthService(cfg *config.Config, repo repository.Authorization) *AuthService {
	return &AuthService{cfg, repo}
}

func (s *AuthService) CreateUser(ctx context.Context, user domain.User) (int, error) {
	user.Password = generatePasswordHash(s.cfg, user.Password)
	return s.repo.CreateUser(ctx, user)
}

func (s *AuthService) SignIn(ctx context.Context, input domain.SignInInput) (string, string, error) {
	user, err := s.repo.GetUser(ctx, input.Username, generatePasswordHash(s.cfg, input.Password))
	if err != nil {
		return "", "", err
	}

	return s.GenerateTokens(ctx, user.Id)
}

func (s *AuthService) GenerateTokens(ctx context.Context, userId int) (string, string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &tokenClaims{
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(s.cfg.Auth.TokenTTL).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		userId,
	})

	accessToken, err := token.SignedString([]byte(s.cfg.Key.SigningKey))
	if err != nil {
		return "", "", err
	}

	refreshToken, err := newRefreshToken()
	if err != nil {
		return "", "", err
	}

	if err := s.repo.CreateToken(ctx, domain.RefreshSession{
		UserID:    userId,
		Token:     refreshToken,
		ExpiresAt: time.Now().Add(s.cfg.Auth.RefreshTokenTTL),
	}); err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func (s *AuthService) ParseToken(accessToken string) (int, error) {
	token, err := jwt.ParseWithClaims(accessToken, &tokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}

		return []byte(s.cfg.Key.SigningKey), nil
	})
	if err != nil {
		return 0, err
	}

	claims, ok := token.Claims.(*tokenClaims)
	if !ok {
		return 0, errors.New("token claims are not of type *tokenClaims")
	}

	return claims.UserId, nil
}

func generatePasswordHash(cfg *config.Config, password string) string {
	hash := sha1.New()
	hash.Write([]byte(password))

	return fmt.Sprintf("%x", hash.Sum([]byte(cfg.Key.Salt)))
}

func newRefreshToken() (string, error) {
	b := make([]byte, 32)

	s := rand.NewSource(time.Now().Unix())
	r := rand.New(s)

	if _, err := r.Read(b); err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", b), nil
}

func (s *AuthService) RefreshTokens(ctx context.Context, refreshToken string) (string, string, error) {
	session, err := s.repo.GetToken(ctx, refreshToken)
	if err != nil {
		return "", "", err
	}

	if session.ExpiresAt.Unix() < time.Now().Unix() {
		return "", "", errors.New("refresh token expired")
	}

	return s.GenerateTokens(ctx, session.UserID)
}
