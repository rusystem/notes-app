package repository

import (
	"context"
	"github.com/jmoiron/sqlx"
	"github.com/rusystem/notes-app/internal/domain"
	"github.com/rusystem/notes-app/internal/repository/psql"
)

type Authorization interface {
	CreateUser(ctx context.Context, user domain.User) (int, error)
	GetUser(ctx context.Context, username, password string) (domain.User, error)
	GetToken(ctx context.Context, token string) (domain.RefreshSession, error)
	CreateToken(ctx context.Context, token domain.RefreshSession) error
}

type Note interface {
	Create(ctx context.Context, userId int, note domain.Note) (int, error)
	GetByID(ctx context.Context, userId, id int) (domain.Note, error)
	GetAll(ctx context.Context, userId int) ([]domain.Note, error)
	Delete(ctx context.Context, userId, id int) error
	Update(ctx context.Context, userId, id int, newNote domain.UpdateNote) error
}

type Repository struct {
	Authorization
	Note
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Authorization: psql.NewAuthRepository(db),
		Note:          psql.NewNoteRepository(db),
	}
}
