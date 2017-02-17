package types

import (
	"time"
)

type Configuration struct {
	Redis        string        `default:"localhost:6379"`
	Bind         string        `default:"0.0.0.0"`
	PollInterval time.Duration `default:"10s" split_words:"true"`
}

type Startable struct {
	Start func(Configuration)
}
