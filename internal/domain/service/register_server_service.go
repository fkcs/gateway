package service

import (
	"github.com/fkcs/gateway/internal/context"
	"github.com/fkcs/gateway/internal/domain/entity"
	"github.com/fkcs/gateway/internal/domain/repository"
	"github.com/fkcs/gateway/internal/domain/vo"
	"github.com/fkcs/gateway/internal/infrastructure/cache"
	"github.com/fkcs/gateway/internal/infrastructure/logger"
	repository2 "github.com/fkcs/gateway/internal/infrastructure/repository"
	"github.com/fkcs/gateway/internal/interfaces/command"
	"github.com/fkcs/gateway/internal/interfaces/dto"
	errord "github.com/fkcs/gateway/internal/utils/error"
	"github.com/fkcs/gateway/internal/utils/types"
	"encoding/json"
)

type RegisterServerDomainImpl struct {
	ctx     *context.Ctx
	repo    repository.RegisterRepoInterface
	server  cache.Cache
	cluster cache.Cache
	bind    cache.Cache
}

func NewRegisterDomainImpl(ctx *context.Ctx) *RegisterServerDomainImpl {
	registerImpl := &RegisterServerDomainImpl{
		ctx:     ctx,
		repo:    repository2.NewClusterStore(ctx.Redis.ClusterClient),
		server:  cache.ServersCache,
		cluster: cache.ClustersCache,
		bind:    cache.BindsClusterCache,
	}
	return registerImpl
}

func (x *RegisterServerDomainImpl) Init(ctx *context.Ctx) {

}

func (x *RegisterServerDomainImpl) initServers(ctx *context.Ctx) {
	serversInfo, err := x.repo.GetServers()
	if err != nil {
		panic(err) // TODO 待优化优雅退出
	}
	var serverInfo command.ServerInfo
	for _, v := range serversInfo {
		if err := json.Unmarshal([]byte(v), &serverInfo); err != nil {
			logger.Logger().Errorf("failed to unmarshal ServerInfo!%v", err)
			continue
		}
		x.AddServer(&serverInfo)
	}
}

func (x *RegisterServerDomainImpl) AddServer(info *command.ServerInfo) dto.ErrorCode {
	logger.Logger().Debugf("add server (%+v)", info)
	var serverMeta *entity.ServiceEntity
	serverInfo, _, ok := x.server.Get(info.Addr)
	if ok != nil {
		serverMeta = &entity.ServiceEntity{
			Cluster:     info.Cluster,
			Addr:        info.Addr,
			Health:      true,
			HealthCheck: 0,
			CurStatus:   types.ServerUpStatus,
			PreStatus:   types.ServerUpStatus,
		}

		if err := x.repo.AddServer(info); err != nil {
			return errord.MakeInternalServerError(types.RedisInternal)
		}
		x.bindClusterServer(serverMeta)
		return errord.MakeOkRsp(true)
	} else {
		serverMeta = serverInfo.(*entity.ServiceEntity)
		serverMeta.Health = true
		serverMeta.HealthCheck = 1
		serverMeta.CurStatus = types.ServerUpStatus
		x.updateBindClusterServer(serverMeta)
		logger.Logger().Warnf("(%+v) had bean existed, please don't add again!%v",
			info, serverMeta)

	}
	return errord.MakeOkRsp(false)
}

// 绑定节点和后端服务
func (x *RegisterServerDomainImpl) bindClusterServer(meta *entity.ServiceEntity) {
	logger.Logger().Infof("bind server (%+v)", meta)
	var servers = make([]*entity.ServiceEntity, 0)
	servers = append(servers, meta)
	binds, _, err := x.bind.Get(meta.Cluster)
	if err == nil {
		bindServers := binds.([]*entity.ServiceEntity)
		for _, server := range bindServers {
			if server.Cluster == meta.Cluster && server.Addr == meta.Addr {
				logger.Logger().Warnf("%v had been bind", meta)
				return
			}
		}
		servers = append(servers, bindServers...)
	}
	x.server.Put(meta.Addr, meta, 0)
	x.bind.Put(meta.Cluster, servers, 0)
}

// 更新绑定关系
func (x *RegisterServerDomainImpl) updateBindClusterServer(meta *entity.ServiceEntity) {
	logger.Logger().Debugf("update bind server (%v)", meta)
	x.server.Put(meta.Addr, meta, 0)
	if binds, _, ok := x.bind.Get(meta.Cluster); ok != nil {
		servers := binds.([]*entity.ServiceEntity)
		for k, v := range servers {
			if v.Cluster == meta.Cluster && v.Addr == meta.Addr {
				servers[k].CurStatus = meta.CurStatus
				servers[k].PreStatus = meta.PreStatus
				servers[k].Health = meta.Health
				servers[k].HealthCheck = meta.HealthCheck
			}
		}
		x.bind.Put(meta.Cluster, servers, 0)
	}
}

func (x *RegisterServerDomainImpl) DeleteServer(info *command.ServerInfo) dto.ErrorCode {
	logger.Logger().Infof("delete server (%v)", info)
	if server, _, ok := x.server.Get(info.Addr); ok == nil {
		bindServer := server.(*entity.ServiceEntity)
		if info.Cluster != bindServer.Cluster {
			return errord.MakeBadRequest(types.InvalidParam)
		}
		binds, _, err := x.bind.Get(bindServer.Cluster)
		if err != nil {
			logger.Logger().Errorf("%v no binds", info)
			return errord.MakeBadRequest(types.InvalidParam)
		}
		bindServers := binds.([]*entity.ServiceEntity)
		for k, v := range bindServers {
			if v.Addr == info.Addr && v.Cluster == info.Cluster {
				bindServers = append(bindServers[0:k], bindServers[k+1:]...)
			}
		}
		if err := x.repo.DelServer(info); err != nil {
			return errord.MakeInternalServerError(types.RedisInternal)
		}
		x.bind.Put(info.Cluster, bindServers, 0)
		x.server.Delete(info.Addr)
	} else {
		logger.Logger().Warnf("no such host (%v)", info)
	}
	return errord.MakeOkRsp(nil)
}

func (x *RegisterServerDomainImpl) GetServers() dto.ErrorCode {
	servers, err := x.repo.GetServers()
	if err != nil {
		logger.Logger().Errorf("server:%v,err:%v", servers, err)
		return errord.MakeInternalServerError(types.RedisInternal)
	}
	var serversInfo []vo.ServerVo
	var serverInfo vo.ServerVo
	for _, v := range servers {
		if err := json.Unmarshal([]byte(v), &serverInfo); err != nil {
			logger.Logger().Errorf("%v", err)
			return errord.MakeBadRequest(types.UnmarshalErr)
		}
		serversInfo = append(serversInfo, serverInfo)
	}
	return errord.MakeOkRsp(serversInfo)
}

func (x RegisterServerDomainImpl) GetServer(addr string) dto.ErrorCode {
	server, _, err := x.server.Get(addr)
	if err != nil {
		logger.Logger().Errorf("addr:%s,err:%v", addr, err)
		return errord.MakeBadRequest(types.InvalidParam)
	}
	return errord.MakeOkRsp(server)
}

func (x *RegisterServerDomainImpl) GetBindInfo(url string) dto.ErrorCode {
	binds, _, err := x.bind.Get(url)
	if err != nil {
		logger.Logger().Errorf("%v", err)
		return errord.MakeBadRequest(types.InvalidParam)
	}
	return errord.MakeOkRsp(binds)
}

func (x *RegisterServerDomainImpl) UpServer(info *command.ServerInfo) dto.ErrorCode {
	server, _, ok := x.server.Get(info.Addr)
	if ok != nil {
		logger.Logger().Errorf("no such server %v", server)
		return errord.MakeBadRequest(types.InvalidParam)
	}
	serverMeta := server.(*entity.ServiceEntity)
	curStatus := serverMeta.CurStatus
	serverMeta.PreStatus = curStatus
	serverMeta.Health = true
	serverMeta.CurStatus = types.ServerUpStatus
	serverMeta.HealthCheck = 1
	x.updateBindClusterServer(serverMeta)
	return errord.MakeOkRsp(nil)
}

func (x *RegisterServerDomainImpl) DownServer(info *command.ServerInfo, ctx *context.Ctx) dto.ErrorCode {
	server, _, ok := x.server.Get(info.Addr)
	if ok != nil {
		logger.Logger().Errorf("no such server %v", server)
		return errord.MakeBadRequest(types.InvalidParam)
	}
	serverMeta := server.(*entity.ServiceEntity)
	serverMeta.Health = false
	curStatus := serverMeta.CurStatus
	serverMeta.PreStatus = curStatus
	serverMeta.CurStatus = types.ServerDownStatus
	serverMeta.HealthCheck++
	if serverMeta.HealthCheck >= ctx.GateWay.HealthCheck.ServiceAfterDel {
		x.DeleteServer(info)
		logger.Logger().Errorf("%v is DOWN, please delete it", info)
		return errord.MakeOkRsp(false)
	}
	x.updateBindClusterServer(serverMeta)
	return errord.MakeOkRsp(true)
}
