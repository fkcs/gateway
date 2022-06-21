package service

import (
	"github.com/fkcs/gateway/internal/context"
	"github.com/fkcs/gateway/internal/domain/adaptor"
	"github.com/fkcs/gateway/internal/domain/service"
	"github.com/fkcs/gateway/internal/interfaces/command"
	"github.com/fkcs/gateway/internal/interfaces/dto"
)

type CommonAppImpl struct {
	ctx     *context.Ctx
	log     adaptor.LogDomainInterface
	server  adaptor.RegisterServerInterface
	cluster adaptor.RegisterClusterInterface
}

func NewCommonAppImpl(ctx *context.Ctx) *CommonAppImpl {
	return &CommonAppImpl{
		ctx:     ctx,
		log:     new(service.LogDomainImpl),
		server:  service.NewRegisterDomainImpl(ctx),
		cluster: service.NewRegisterClusterDomainImpl(ctx),
	}
}

func (x *CommonAppImpl) SetLogLevel(level string) dto.ErrorCode {
	return x.log.SetLevel(level)
}

func (x *CommonAppImpl) GetLogLevel() dto.ErrorCode {
	return x.log.GetLevel()
}

func (x *CommonAppImpl) GetServer(addr string) dto.ErrorCode {
	return x.server.GetServer(addr)
}

func (x *CommonAppImpl) GetServers() dto.ErrorCode {
	return x.server.GetServers()
}

func (x *CommonAppImpl) AddServer(info *command.ServerInfo) dto.ErrorCode {
	return x.server.AddServer(info)
}

func (x *CommonAppImpl) DeleteServer(info *command.ServerInfo) dto.ErrorCode {
	return x.server.DeleteServer(info)
}

func (x *CommonAppImpl) GetCluster(cluster string) dto.ErrorCode {
	return x.cluster.GetClusterInfo(cluster)
}
