package repository

import (
	"github.com/fkcs/gateway/internal/infrastructure/logger"
	"github.com/fkcs/gateway/internal/interfaces/command"
	"github.com/fkcs/gateway/internal/utils/types"
	"encoding/json"
	"github.com/go-redis/redis"
)

type ClusterStore struct {
	redis *redis.ClusterClient
}

func NewClusterStore(redis *redis.ClusterClient) *ClusterStore {
	return &ClusterStore{
		redis: redis,
	}
}

// 增加服务地址信息
func (x *ClusterStore) AddServer(info *command.ServerInfo) error {
	logger.Logger().Infof("[Redis] add %v", info)
	serverInfo, err := json.Marshal(info)
	if err != nil {
		logger.Logger().Errorf("failed to unmarshal!%v", err)
		return err
	}
	if _, err := x.redis.SAdd(types.RedisServerKey, string(serverInfo)).Result(); err != nil {
		logger.Logger().Errorf("failed to add server_info!%v", err)
		return err
	}
	return nil
}

// 获取服务地址信息
func (x *ClusterStore) GetServers() ([]string, error) {
	serversInfo, err := x.redis.SMembers(types.RedisServerKey).Result()
	if err != nil {
		return nil, err
	}
	logger.Logger().Infof("[Redis] get %v", serversInfo)
	return serversInfo, nil
}

// 删除节点信息
func (x *ClusterStore) DelServer(info *command.ServerInfo) error {
	logger.Logger().Infof("[Redis] delete %v", info)
	serverInfo, err := json.Marshal(info)
	if err != nil {
		logger.Logger().Errorf("failed to unmarshal!%v", err)
		return err
	}
	if err := x.redis.SRem(types.RedisServerKey, serverInfo).Err(); err != nil {
		logger.Logger().Errorf("failed to delete server(%v)!%v", serverInfo, err)
		return err
	}
	return nil
}
