package service

import (
	"github.com/rusystem/cache"
	"github.com/rusystem/notes-app/internal/config"
	"github.com/rusystem/notes-app/internal/domain"
	"github.com/rusystem/notes-app/internal/repository"
)

type Authorization interface {
	CreateUser(user domain.User) (int, error)
	GenerateToken(username, password string) (string, error)
	ParseToken(token string) (int, error)
}

type Note interface {
	Create(userId int, note domain.Note) (int, error)
	GetByID(userId, id int) (domain.Note, error)
	GetAll(userId int) ([]domain.Note, error)
	Delete(userId, id int) error
	Update(userId, id int, newNote domain.UpdateNote) error
}

type Service struct {
	Authorization
	Note
}

func NewService(cfg *config.Config, c *cache.Cache, repos *repository.Repository) *Service {
	return &Service{
		Authorization: NewAuthService(cfg, repos.Authorization),
		Note:          NewNoteService(cfg, c, repos.Note),
	}
}
