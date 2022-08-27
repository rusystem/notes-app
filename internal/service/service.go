package service

import (
	"context"
	"github.com/rusystem/cache"
	"github.com/rusystem/notes-app/internal/config"
	"github.com/rusystem/notes-app/internal/domain"
	"github.com/rusystem/notes-app/internal/repository"
)

type Authorization interface {
	CreateUser(ctx context.Context, user domain.User) (int, error)
	GenerateToken(ctx context.Context, username, password string) (string, error)
	ParseToken(token string) (int, error)
}

type Note interface {
	Create(ctx context.Context, userId int, note domain.Note) (int, error)
	GetByID(ctx context.Context, userId, id int) (domain.Note, error)
	GetAll(ctx context.Context, userId int) ([]domain.Note, error)
	Delete(ctx context.Context, userId, id int) error
	Update(ctx context.Context, userId, id int, newNote domain.UpdateNote) error
}

type Service struct {
	Authorization
	Note
}

func NewService(cfg *config.Config, cache *cache.Cache, repos *repository.Repository) *Service {
	return &Service{
		Authorization: NewAuthService(cfg, repos.Authorization),
		Note:          NewNoteService(cfg, cache, repos.Note),
	}
}
