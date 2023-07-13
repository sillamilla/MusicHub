package Redis_storage

import (
	"context"
	"github.com/redis/go-redis/v9"
	"time"
)

type Storage interface {
	UpsertSession(ctx context.Context, id string, session string) error
	Logout(ctx context.Context, id string) error
	GetSession(ctx context.Context, id string) (string, error)
}

type redisDB struct {
	re *redis.Client
}

func New(re *redis.Client) Storage {
	return &redisDB{
		re: re,
	}
}

func (db *redisDB) UpsertSession(ctx context.Context, id string, session string) error {
	key := "sessions:" + id
	err := db.re.Set(ctx, key, session, 5*time.Hour).Err()
	if err != nil {
		return err
	}

	return nil
}

func (db *redisDB) Logout(ctx context.Context, id string) error {
	key := "sessions:" + id
	err := db.re.Del(ctx, key).Err()
	if err != nil {
		return err
	}

	return nil
}

func (db *redisDB) GetSession(ctx context.Context, id string) (string, error) {
	key := "sessions:" + id
	session, err := db.re.Get(ctx, key).Result()
	if err != nil {
		return "", err
	}

	return session, nil

}
