package rdb

import (
	"context"
	"github.com/go-redis/redis/v9"
	"strconv"
	"time"
)

type SessionRepository struct {
	rdbClient *redis.Client
}

func NewSessionRepository(rdbClient *redis.Client) *SessionRepository {
	return &SessionRepository{rdbClient: rdbClient}
}

func (r *SessionRepository) Set(ctx context.Context, token string, userId int, ttl time.Duration) error {
	return r.rdbClient.Set(ctx, token, userId, ttl).Err()
}

func (r *SessionRepository) Delete(ctx context.Context, token string) error {
	return r.rdbClient.Del(ctx, token).Err()
}

func (r *SessionRepository) Get(ctx context.Context, token string) (int, error) {
	result, err := r.rdbClient.Get(ctx, token).Result()
	if err != nil {
		return 0, err
	}

	userId, err := strconv.Atoi(result)
	if err != nil {
		return 0, err
	}

	return userId, nil
}
