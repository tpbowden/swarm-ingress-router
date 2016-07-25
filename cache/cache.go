package cache

import (
	"errors"
	"time"

	"github.com/garyburd/redigo/redis"
)

// Cache holds the connection pool for an instance of the cache
type Cache struct {
	pool *redis.Pool
}

// Set sets a key to a value
func (c *Cache) Set(key, value string) error {
	conn := c.pool.Get()
	defer conn.Close()

	conn.Send("PUBLISH", "ingress-router", "updated")
	if _, err := conn.Do("SET", key, value); err != nil {
		return err
	}
	return nil
}

// Subscribe subscribes to a key and calls a function when messages are received
func (c *Cache) Subscribe(channel string, action func()) error {
	conn := c.pool.Get()
	defer conn.Close()

	if _, err := conn.Do("SUBSCRIBE", "ingress-router"); err != nil {
		return err
	}

	for {
		if _, err := conn.Receive(); err != nil {
			return err
		}

		action()
	}
}

// Get retrieves a key from redis
func (c *Cache) Get(key string) ([]byte, error) {
	var s []byte

	conn := c.pool.Get()
	defer conn.Close()

	result, err := conn.Do("GET", key)
	if err != nil {
		return s, err
	}

	b, ok := result.([]byte)
	if !ok {
		return s, errors.New("Failed to parse value as a string")
	}

	return b, nil
}

// NewCache returns a new instance of the cache with a prepared connection pool
func NewCache(address string) Cache {
	pool := &redis.Pool{
		MaxIdle:     2,
		IdleTimeout: 300 * time.Second,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", address)
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}

	return Cache{pool: pool}
}
