package service

import (
	"context"
	"github.com/rusystem/cache"
	"github.com/rusystem/notes-app/internal/config"
	"github.com/rusystem/notes-app/internal/domain"
	"github.com/rusystem/notes-app/internal/repository"
	logs "github.com/rusystem/notes-log/pkg/domain"
)

type Authorization interface {
	CreateUser(ctx context.Context, user domain.User) (int, error)
	GetSession(ctx context.Context, token string) (int, error)
	SignIn(ctx context.Context, inp domain.SignInInput) (domain.Cookie, error)
	Logout(ctx context.Context, token string) error
}

type Note interface {
	Create(ctx context.Context, userId int, note domain.Note) (int, error)
	GetByID(ctx context.Context, userId, id int) (domain.Note, error)
	GetAll(ctx context.Context, userId int) ([]domain.Note, error)
	Delete(ctx context.Context, userId, id int) error
	Update(ctx context.Context, userId, id int, newNote domain.UpdateNote) error
}

type LogsClient interface {
	LogRequest(ctx context.Context, req logs.LogItem) error
}

type Service struct {
	Authorization
	Note
	LogsClient
}

func NewService(cfg *config.Config, cache *cache.Cache, repos *repository.Repository, logsClient LogsClient) *Service {
	return &Service{
		Authorization: NewAuthService(cfg, repos.Authorization, repos.Session, logsClient),
		Note:          NewNoteService(cfg, cache, repos.Note, logsClient),
	}
}
