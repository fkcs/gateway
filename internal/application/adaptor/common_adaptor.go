package adaptor

import (
	"github.com/fkcs/gateway/internal/interfaces/command"
	"github.com/fkcs/gateway/internal/interfaces/dto"
)

type CommonInterface interface {
	SetLogLevel(level string) dto.ErrorCode
	GetLogLevel() dto.ErrorCode
	GetServer(addr string) dto.ErrorCode
	GetServers() dto.ErrorCode
	AddServer(info *command.ServerInfo) dto.ErrorCode
	DeleteServer(info *command.ServerInfo) dto.ErrorCode
	GetCluster(cluster string) dto.ErrorCode
}
