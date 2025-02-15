package cache

import (
	"github.com/redis/go-redis/v9"
)

type Storage struct {
	db *redis.Client
}

func InitRedis(adr string) *Storage {
	db := redis.NewClient(&redis.Options{
		Addr: adr,
	})

	cacheRepo := Storage{
		db: db,
	}

	return &cacheRepo
}
