package adaptor

import (
	"github.com/fkcs/gateway/internal/context"
	"github.com/fkcs/gateway/internal/interfaces/command"
	"github.com/fkcs/gateway/internal/interfaces/dto"
)

type TransportInterface interface {
	AddServer(info *command.ServerInfo) dto.ErrorCode
	DelServer(info *command.ServerInfo) dto.ErrorCode
	GetServers() dto.ErrorCode
	UpdateServer(info *command.ServerInfo) dto.ErrorCode
	BindClusterServer(info *command.ServerInfo) dto.ErrorCode
	UpdateCluster(ctx *context.Ctx)
	UpServer(info *command.ServerInfo) dto.ErrorCode
	DownServer(info *command.ServerInfo, ctx *context.Ctx) dto.ErrorCode
	GetCluster(name string) dto.ErrorCode
	GetBind(cluster string) dto.ErrorCode
}
