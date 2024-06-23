package main

import (
	"fmt"
	"os"
)

type Env struct {
	RedisUrl string
	Port     string
}

func NewEnv() (Env, error) {
	redis, ok := os.LookupEnv("REDIS_URL")
	if !ok {
		return Env{}, fmt.Errorf("REDIS_URL not found")
	}

	port, ok := os.LookupEnv("PORT")
	if !ok {
		port = "8080"
	}

	return Env{
		RedisUrl: redis,
		Port:     port,
	}, nil
}

func (e *Env) Addr() string {
	return fmt.Sprintf(":%s", e.Port)
}
