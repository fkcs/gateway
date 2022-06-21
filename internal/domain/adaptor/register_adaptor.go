package adaptor

import (
	"github.com/fkcs/gateway/internal/context"
	"github.com/fkcs/gateway/internal/interfaces/command"
	"github.com/fkcs/gateway/internal/interfaces/dto"
)

type RegisterServerInterface interface {
	Init(ctx *context.Ctx)
	AddServer(info *command.ServerInfo) dto.ErrorCode
	DeleteServer(info *command.ServerInfo) dto.ErrorCode
	GetServers() dto.ErrorCode
	GetServer(addr string) dto.ErrorCode
	UpServer(info *command.ServerInfo) dto.ErrorCode
	DownServer(info *command.ServerInfo, ctx *context.Ctx) dto.ErrorCode
	GetBindInfo(url string) dto.ErrorCode
}

type RegisterClusterInterface interface {
	Init(ctx *context.Ctx)
	GetClusterInfo(name string) dto.ErrorCode
	UpdateCluster(ctx *context.Ctx) dto.ErrorCode
}
