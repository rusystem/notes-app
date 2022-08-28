package psql

import (
	"context"
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

func (r *AuthRepository) CreateUser(ctx context.Context, user domain.User) (int, error) {
	query := fmt.Sprintf("INSERT INTO %s (name, username, password_hash) VALUES ($1, $2, $3) RETURNING id", usersTable)

	var id int
	row := r.db.QueryRowContext(ctx, query, user.Name, user.Username, user.Password)
	if err := row.Scan(&id); err != nil {
		return 0, err
	}

	return id, nil
}

func (r *AuthRepository) GetUser(ctx context.Context, username, password string) (domain.User, error) {
	query := fmt.Sprintf("SELECT id from %s WHERE username=$1 AND password_hash=$2", usersTable)

	var user domain.User
	err := r.db.GetContext(ctx, &user, query, username, password)

	return user, err
}

func (r *AuthRepository) CreateToken(ctx context.Context, token domain.RefreshSession) error {
	query := fmt.Sprintf("INSERT INTO %s (user_id, token, expires_at) values ($1, $2, $3)",
		refreshTokensTable)
	_, err := r.db.ExecContext(ctx, query, token.UserID, token.Token, token.ExpiresAt)

	return err
}

func (r *AuthRepository) GetToken(ctx context.Context, token string) (domain.RefreshSession, error) {
	query := fmt.Sprintf("SELECT id, user_id, token, expires_at FROM %s WHERE token = $1",
		refreshTokensTable)

	var t domain.RefreshSession
	if err := r.db.GetContext(ctx, &t, query, token); err != nil {
		return t, err
	}

	query = fmt.Sprintf("DELETE FROM %s WHERE user_id=$1", refreshTokensTable)
	_, err := r.db.ExecContext(ctx, query, t.UserID)

	return t, err
}
