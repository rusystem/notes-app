package service

import (
	"context"
	"crypto/sha1"
	"fmt"
	"github.com/rusystem/notes-app/internal/config"
	"github.com/rusystem/notes-app/internal/domain"
	"github.com/rusystem/notes-app/internal/repository"
	"math/rand"
	"time"
)

type AuthService struct {
	cfg     *config.Config
	repo    repository.Authorization
	session repository.Session
}

func NewAuthService(cfg *config.Config, repo repository.Authorization, session repository.Session) *AuthService {
	return &AuthService{cfg, repo, session}
}

func (s *AuthService) CreateUser(ctx context.Context, user domain.User) (int, error) {
	user.Password = generatePasswordHash(s.cfg, user.Password)
	return s.repo.CreateUser(ctx, user)
}

func (s *AuthService) SignIn(ctx context.Context, inp domain.SignInInput) (domain.Cookie, error) {
	user, err := s.repo.GetUser(ctx, inp.Username, generatePasswordHash(s.cfg, inp.Password))
	if err != nil {
		return domain.Cookie{}, err
	}

	token, err := newToken()
	if err != nil {
		return domain.Cookie{}, err
	}

	err = s.session.Set(ctx, token, user.Id, s.cfg.Auth.SessionTTL)
	if err != nil {
		return domain.Cookie{}, err
	}

	return domain.Cookie{Name: domain.AuthCookie, Token: token, MaxAge: int(s.cfg.Auth.SessionTTL.Seconds())}, nil
}

func (s *AuthService) GetSession(ctx context.Context, token string) (int, error) {
	return s.session.Get(ctx, token)
}

func (s *AuthService) Logout(ctx context.Context, token string) error {
	return s.session.Delete(ctx, token)
}

func newToken() (string, error) {
	b := make([]byte, 32)

	s := rand.NewSource(time.Now().Unix())
	r := rand.New(s)

	if _, err := r.Read(b); err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", b), nil
}

func generatePasswordHash(cfg *config.Config, password string) string {
	hash := sha1.New()
	hash.Write([]byte(password))

	return fmt.Sprintf("%x", hash.Sum([]byte(cfg.Key.Salt)))
}
