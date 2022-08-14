package repository

import (
	"github.com/jmoiron/sqlx"
	"github.com/rusystem/notes-app/internal/domain"
	"github.com/rusystem/notes-app/internal/repository/psql"
)

type Authorization interface {
	CreateUser(user domain.User) (int, error)
	GetUser(username, password string) (domain.User, error)
}

type Note interface {
	Create(userId int, note domain.Note) (int, error)
	GetByID(userId, id int) (domain.Note, error)
	GetAll(userId int) ([]domain.Note, error)
	Delete(userId, id int) error
	Update(userId, id int, newNote domain.UpdateNote) error
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
