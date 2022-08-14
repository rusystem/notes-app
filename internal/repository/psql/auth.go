package psql

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/rusystem/notes-app/internal/domain"
)

type AuthRepository struct {
	db *sqlx.DB
}

func NewAuthRepository(db *sqlx.DB) *AuthRepository {
	return &AuthRepository{db: db}
}

func (r *AuthRepository) CreateUser(user domain.User) (int, error) {
	query := fmt.Sprintf("INSERT INTO %s (name, username, password_hash) VALUES ($1, $2, $3) RETURNING id", usersTable)

	var id int
	row := r.db.QueryRow(query, user.Name, user.Username, user.Password)
	if err := row.Scan(&id); err != nil {
		return 0, err
	}

	return id, nil
}

func (r *AuthRepository) GetUser(username, password string) (domain.User, error) {
	query := fmt.Sprintf("SELECT id from %s WHERE username=$1 AND password_hash=$2", usersTable)

	var user domain.User
	err := r.db.Get(&user, query, username, password)

	return user, err
}
