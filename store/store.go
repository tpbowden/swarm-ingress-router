package store

import (
	"fmt"
	"time"

	"github.com/garyburd/redigo/redis"

	"github.com/tpbowden/swarm-ingress-router/types"
)

type Store interface {
	Set(string, string)
}

type redisStore struct {
	pool   *redis.Pool
	config types.Configuration
}

func (s redisStore) Set(key, value string) {
	conn := s.pool.Get()
	defer conn.Close()

	if _, err := conn.Do("SET", key, value); err != nil {
		panic(fmt.Sprintf("Failed to store key '%s' as '%s': %v", key, value, err))
	}
}

func NewStore(config types.Configuration) Store {
	pool := &redis.Pool{
		MaxIdle:     2,
		IdleTimeout: 300 * time.Second,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", config.Redis)
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}
	return Store(redisStore{config: config, pool: pool})
}
