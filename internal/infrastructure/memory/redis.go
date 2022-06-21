package memory

import (
	"github.com/fkcs/gateway/internal/infrastructure/logger"
	"github.com/go-redis/redis"
	"strings"
)

type RedisClient struct {
	*redis.ClusterClient
}

func NewRedisClient(redisAddr, redisPassword string) *RedisClient {
	redisIPs := strings.Split(redisAddr, ",")
	redisAddrs := make([]string, 0, len(redisIPs))
	for _, ip := range redisIPs {
		redisAddrs = append(redisAddrs, ip)
	}
	return &RedisClient{
		ClusterClient: redis.NewClusterClient(&redis.ClusterOptions{
			Addrs:    redisAddrs,
			Password: redisPassword,
		}),
	}
}

func (x *RedisClient) MustConnect() {
	if _, err := x.Ping().Result(); err != nil {
		logger.Logger().Errorf("failed to create redis,%v", err)
		panic(err)
	}
}
