package repository

import (
	"context"
	"github.com/go-redis/redis/v9"
	"github.com/jmoiron/sqlx"
	"github.com/rusystem/notes-app/internal/domain"
	"github.com/rusystem/notes-app/internal/repository/psql"
	"github.com/rusystem/notes-app/internal/repository/rdb"
	"time"
)

type Session interface {
	Set(ctx context.Context, token string, userId int, ttl time.Duration) error
	Delete(ctx context.Context, token string) error
	Get(ctx context.Context, token string) (int, error)
}

type Authorization interface {
	CreateUser(ctx context.Context, user domain.User) (int, error)
	GetUser(ctx context.Context, username, password string) (domain.User, error)
}

type Note interface {
	Create(ctx context.Context, userId int, note domain.Note) (int, error)
	GetByID(ctx context.Context, userId, id int) (domain.Note, error)
	GetAll(ctx context.Context, userId int) ([]domain.Note, error)
	Delete(ctx context.Context, userId, id int) error
	Update(ctx context.Context, userId, id int, newNote domain.UpdateNote) error
}

type Repository struct {
	Session
	Authorization
	Note
}

func NewRepository(db *sqlx.DB, rdbClient *redis.Client) *Repository {
	return &Repository{
		Session:       rdb.NewSessionRepository(rdbClient),
		Authorization: psql.NewAuthRepository(db),
		Note:          psql.NewNoteRepository(db),
	}
}
