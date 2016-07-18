package cache

import (
	"errors"

	"github.com/garyburd/redigo/redis"
)

type Cache struct {
	address string
}

func (c *Cache) Set(key string, value string) error {
	conn, err := redis.Dial("tcp", c.address)

	if err != nil {
		return err
	}

	defer conn.Close()
	if _, setErr := conn.Do("SET", key, value); setErr != nil {
		return setErr
	}
	return nil
}

func (c *Cache) Get(key string) ([]byte, error) {
	var s []byte
	conn, err := redis.Dial("tcp", c.address)

	if err != nil {
		return s, err
	}

	defer conn.Close()

	result, getErr := conn.Do("GET", key)
	if getErr != nil {
		return s, getErr
	}

	b, ok := result.([]byte)

	if !ok {
		return s, errors.New("Failed to parse value as a string")
	}

	return b, nil

}

func NewCache(address string) Cache {
	return Cache{address: address}
}
