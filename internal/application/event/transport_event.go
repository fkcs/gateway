package event

import (
	"github.com/fkcs/gateway/internal/context"
	"github.com/fkcs/gateway/internal/domain/adaptor"
	"github.com/fkcs/gateway/internal/domain/service"
	"github.com/fkcs/gateway/internal/interfaces/command"
	"github.com/fkcs/gateway/internal/interfaces/dto"
	errord "github.com/fkcs/gateway/internal/utils/error"
)

type TransportEventImpl struct {
	server  adaptor.RegisterServerInterface
	cluster adaptor.RegisterClusterInterface
}

func NewTransportEventImpl(ctx *context.Ctx) *TransportEventImpl {
	impl := &TransportEventImpl{
		server:  service.NewRegisterDomainImpl(ctx),
		cluster: service.NewRegisterClusterDomainImpl(ctx),
	}
	impl.server.Init(ctx)
	impl.cluster.Init(ctx)
	return impl
}

func (x *TransportEventImpl) AddServer(info *command.ServerInfo) dto.ErrorCode {
	return x.server.AddServer(info)
}

func (x *TransportEventImpl) DelServer(info *command.ServerInfo) dto.ErrorCode {
	return x.server.DeleteServer(info)
}

func (x *TransportEventImpl) GetServers() dto.ErrorCode {
	return x.server.GetServers()
}

func (x *TransportEventImpl) UpdateServer(info *command.ServerInfo) dto.ErrorCode {
	return errord.MakeOkRsp(nil)
}

func (x *TransportEventImpl) BindClusterServer(info *command.ServerInfo) dto.ErrorCode {
	return errord.MakeOkRsp(nil)
}

func (x *TransportEventImpl) UpdateCluster(ctx *context.Ctx) {
	x.cluster.UpdateCluster(ctx)
}

func (x *TransportEventImpl) UpServer(info *command.ServerInfo) dto.ErrorCode {
	return x.server.UpServer(info)
}

func (x *TransportEventImpl) DownServer(info *command.ServerInfo, ctx *context.Ctx) dto.ErrorCode {
	return x.server.DownServer(info, ctx)
}

func (x *TransportEventImpl) GetCluster(name string) dto.ErrorCode {
	return x.cluster.GetClusterInfo(name)
}

func (x *TransportEventImpl) GetBind(cluster string) dto.ErrorCode {
	return x.server.GetBindInfo(cluster)
}
