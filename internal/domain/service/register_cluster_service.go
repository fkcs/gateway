package service

import (
	"github.com/fkcs/gateway/internal/context"
	"github.com/fkcs/gateway/internal/domain/entity"
	"github.com/fkcs/gateway/internal/domain/proxy/filter/adapter"
	"github.com/fkcs/gateway/internal/domain/proxy/filter/factory"
	"github.com/fkcs/gateway/internal/domain/repository"
	"github.com/fkcs/gateway/internal/infrastructure/cache"
	"github.com/fkcs/gateway/internal/infrastructure/logger"
	repository2 "github.com/fkcs/gateway/internal/infrastructure/repository"
	"github.com/fkcs/gateway/internal/interfaces/dto"
	"github.com/fkcs/gateway/internal/utils/config"
	errord "github.com/fkcs/gateway/internal/utils/error"
	"github.com/fkcs/gateway/internal/utils/types"
	loadbalance "github.com/fkcs/gateway/internal/utils/wrapper/lb"
	"fmt"
)

type RegisterClusterDomainImpl struct {
	ctx     *context.Ctx
	repo    repository.RegisterRepoInterface
	cluster cache.Cache
}

func NewRegisterClusterDomainImpl(ctx *context.Ctx) *RegisterClusterDomainImpl {
	registerImpl := &RegisterClusterDomainImpl{
		ctx:     ctx,
		repo:    repository2.NewClusterStore(ctx.Redis.ClusterClient),
		cluster: cache.ClustersCache,
	}
	//registerImpl.Init(ctx)
	return registerImpl
}

func (x *RegisterClusterDomainImpl) Init(ctx *context.Ctx) {
	x.initCluster(ctx)
}

func (x RegisterClusterDomainImpl) initCluster(ctx *context.Ctx) dto.ErrorCode {
	for _, v := range ctx.GateWay.Routes {
		meta := &entity.ClusterEntity{
			Name:        v.Uri,
			LoadBalance: v.LoadBalance,
			Lb:          x.createLB(v.LoadBalance),
			Filters:     make(map[string]adapter.Filter),
			FilterCfg:   make(map[string]config.Filter),
		}
		meta.Lb.Init()
		for _, filter := range v.Filters {
			filterEng := factory.NewFilterFactory(filter.Name)
			if err := filterEng.Init(filter.Args); err != nil {
				logger.Logger().Errorf("%v", err)
				continue
			}
			meta.FilterCfg[filter.Name] = filter
			meta.Filters[filter.Name] = filterEng
		}
		key := fmt.Sprintf("%s/%s", types.ClusterPrefix, v.Uri)
		x.cluster.Put(key, meta, 0)
	}
	return errord.MakeOkRsp(nil)
}

func (x *RegisterClusterDomainImpl) createLB(loadBalance string) adapter.LoadBalance {
	if lb, ok := loadbalance.LoadLBEng(loadBalance); ok {
		return lb
	} else {
		lb, _ := loadbalance.LoadLBEng(types.RandLoadBalance)
		logger.Logger().Warnf("invalid lb{%v}, but use rand lb", loadBalance)
		return lb
	}
}

func (x *RegisterClusterDomainImpl) GetClusterInfo(name string) dto.ErrorCode {
	cluster, _, ok := x.cluster.Get(name)
	if ok != nil {
		logger.Logger().Errorf("no such cluster %v", name)
		return errord.MakeBadRequest(types.InvalidParam)
	}
	return errord.MakeOkRsp(cluster.(*entity.ClusterEntity))
}

func (x *RegisterClusterDomainImpl) UpdateCluster(ctx *context.Ctx) dto.ErrorCode {
	return x.initCluster(ctx)
}
