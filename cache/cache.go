package cache

import (
	"errors"
	"time"

	"github.com/garyburd/redigo/redis"
)

type Cache struct {
	pool *redis.Pool
}

func (c *Cache) Set(key, value string) error {
	conn := c.pool.Get()
	defer conn.Close()

	conn.Send("PUBLISH", "ingress-router", "updated")
	if _, setErr := conn.Do("SET", key, value); setErr != nil {
		return setErr
	}
	return nil
}

func (c *Cache) Subscribe(channel string, action func()) error {
	conn := c.pool.Get()
	defer conn.Close()

	_, err := conn.Do("SUBSCRIBE", "ingress-router")

	if err != nil {
		return err
	}

	for {
		if _, receiveErr := conn.Receive(); receiveErr != nil {
			return receiveErr
		}

		action()
	}
}

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
