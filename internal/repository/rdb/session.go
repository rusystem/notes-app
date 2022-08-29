package rdb

import (
	"fmt"
	"github.com/gin-contrib/sessions/redis"
	"github.com/rusystem/notes-app/pkg/database"
)

type SessionRepository struct {
	store redis.Store
}

func NewSessionRepository(rdb *database.RedisConnectionInfo) *SessionRepository {
	store, _ := redis.NewStore(rdb.Size, rdb.Network, fmt.Sprintf("localhost:%d", rdb.Port), rdb.Password, []byte(rdb.Key))

	return &SessionRepository{store: store}
}

func (r *SessionRepository) Save() error {
	return nil
}

func (r *SessionRepository) Set(key interface{}) error {
	return nil
}

func (r *SessionRepository) Delete(key interface{}) error {
	return nil
}

func (r *SessionRepository) Get(key interface{}) (int, error) {
	return 0, nil
}
