package store

import (
	"testing"

	"github.com/garyburd/redigo/redis"

	"github.com/tpbowden/swarm-ingress-router/types"
)

var config = types.Configuration{Redis: "localhost:6379"}
var subject = NewStore(config)

func TestSettingAKeyInTheStore(t *testing.T) {
	conn, _ := redis.Dial("tcp", "localhost:6379")
	defer conn.Close()

	conn.Do("DELETE", "ingress:foo")

	subject.Set("foo", "bar")

	result, _ := conn.Do("GET", "ingress:foo")

	if string(result.([]byte)) != "bar" {
		t.Error("Failed to set then retrieve a key from the store")
	}
}
