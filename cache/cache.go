package cache

import (
	"errors"
	"log"
	"time"

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
	conn.Send("PUBLISH", "ingress-router", "updated")
	if _, setErr := conn.Do("SET", key, value); setErr != nil {
		return setErr
	}
	return nil
}

func (c *Cache) Subscribe(channel string, action func()) {
	conn, err := redis.Dial("tcp", c.address)
	if err != nil {
		log.Print("Failed to connect to redis, retrying in 5 seconds")
		time.Sleep(5 * time.Second)
		c.Subscribe(channel, action)
		return
	}

	defer conn.Close()
	conn.Do("SUBSCRIBE", "ingress-router")

	for {
		_, err := conn.Receive()
		if err != nil {
			log.Print("Failed to connect to redis, retrying in 5 seconds")
			time.Sleep(5 * time.Second)
			c.Subscribe(channel, action)
			return
		}

		action()
	}
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
