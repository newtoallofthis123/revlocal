package main

import (
	"context"
	"encoding/json"

	"github.com/newtoallofthis123/ranhash"
	"github.com/redis/go-redis/v9"
)

type Store struct {
	client *redis.Client
	ctx    context.Context
}

type Connection struct {
	Iden      string `json:"iden"`
	Token     string `json:"token"`
	CreatedAt string `json:"created_at"`
}

type ConnectionRequest struct {
	Token string `json:"token"`
}

func NewConnection(req ConnectionRequest) Connection {
	return Connection{
		Iden:      ranhash.GenerateRandomString(10),
		Token:     req.Token,
		CreatedAt: "2021-09-21",
	}
}

func NewStore(env Env) *Store {
	opt, err := redis.ParseURL(env.RedisUrl)
	if err != nil {
		panic(err)
	}

	client := redis.NewClient(opt)

	return &Store{
		client: client,
		ctx:    context.Background(),
	}
}

func ConvertToJson(c Connection) (string, error) {
	val, err := json.Marshal(c)
	if err != nil {
		return "", err
	}

	return string(val), nil
}

func (s *Store) SetConnection(c Connection) error {
	val, err := ConvertToJson(c)
	if err != nil {
		return err
	}

	err = s.client.Set(s.ctx, c.Iden, val, 0).Err()
	if err != nil {
		return err
	}

	return nil
}

func (s *Store) GetConnection(iden string) (Connection, error) {
	val, err := s.client.Get(s.ctx, iden).Result()
	if err != nil {
		return Connection{}, err
	}

	var c Connection
	err = json.Unmarshal([]byte(val), &c)
	if err != nil {
		return Connection{}, err
	}

	return c, nil
}

func (s *Store) DeleteConnection(iden string) error {
	err := s.client.Del(s.ctx, iden).Err()
	if err != nil {
		return err
	}

	return nil
}
