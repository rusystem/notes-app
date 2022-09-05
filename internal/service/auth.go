package service

import (
	"context"
	"crypto/sha1"
	"fmt"
	"github.com/rusystem/notes-app/internal/config"
	"github.com/rusystem/notes-app/internal/domain"
	"github.com/rusystem/notes-app/internal/repository"
	logs "github.com/rusystem/notes-log/pkg/domain"
	"github.com/sirupsen/logrus"
	"math/rand"
	"time"
)

type AuthService struct {
	cfg        *config.Config
	repo       repository.Authorization
	session    repository.Session
	logsClient LogsClient
}

func NewAuthService(cfg *config.Config, repo repository.Authorization, session repository.Session, logsClient LogsClient) *AuthService {
	return &AuthService{cfg, repo, session, logsClient}
}

func (s *AuthService) CreateUser(ctx context.Context, user domain.User) (int, error) {
	user.Password = generatePasswordHash(s.cfg, user.Password)
	id, err := s.repo.CreateUser(ctx, user)
	if err != nil {
		return 0, err
	}

	if err := s.logsClient.LogRequest(ctx, logs.LogItem{
		Entity:    logs.ENTITY_USER,
		Action:    logs.ACTION_REGISTER,
		EntityID:  int64(id),
		Timestamp: time.Now(),
	}); err != nil {
		logrus.WithFields(logrus.Fields{
			"method": "auth.SignUp",
		}).Error("failed to send log request:", err)
	}

	return id, nil
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

	if err := s.logsClient.LogRequest(ctx, logs.LogItem{
		Entity:    logs.ENTITY_USER,
		Action:    logs.ACTION_LOGIN,
		EntityID:  int64(user.Id),
		Timestamp: time.Now(),
	}); err != nil {
		logrus.WithFields(logrus.Fields{
			"method": "auth.SignIn",
		}).Error("failed to send log request:", err)
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
