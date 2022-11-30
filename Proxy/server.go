package main

import (
	"RemoteCmd/Common"
	"fmt"
	"time"
)
import "github.com/go-redis/redis"

type Server struct {
	Common.Server
}

func (s *Server) getRedisDb() *redis.Client {
	addr := fmt.Sprintf(s.Ip, 6379)
	rdb := redis.NewClient(&redis.Options{
		Addr:       addr,
		Password:   "",
		DB:         0,
		MaxConnAge: 10 * time.Second,
	})
	return rdb
}
