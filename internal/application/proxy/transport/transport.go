// 承接proxy与后端服务之间的纽带
package transport

import (
	"github.com/fkcs/gateway/internal/context"
	adaptor2 "github.com/fkcs/gateway/internal/domain/adaptor"
	"github.com/fkcs/gateway/internal/domain/entity"
	"github.com/fkcs/gateway/internal/domain/service"
	"github.com/fkcs/gateway/internal/domain/vo"
	"github.com/fkcs/gateway/internal/infrastructure/logger"
	"github.com/fkcs/gateway/internal/infrastructure/repository"
	"github.com/fkcs/gateway/internal/interfaces/dto"
	"github.com/fkcs/gateway/internal/utils/config"
	errord "github.com/fkcs/gateway/internal/utils/error"
	"github.com/fkcs/gateway/internal/utils/types"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/valyala/fasthttp"
	"net/http"
	"strings"
	"sync"
)

type Transport struct {
	mux *sync.RWMutex
	*repository.NodeExporter
	*repository.ClusterStore
	cluster adaptor2.RegisterClusterInterface
	server  adaptor2.RegisterServerInterface
}

func NewTransport(ctx *context.Ctx) *Transport {
	transport := &Transport{
		mux:     &sync.RWMutex{},
		server:  service.NewRegisterDomainImpl(ctx),
		cluster: service.NewRegisterClusterDomainImpl(ctx),
	}
	transport.init(ctx)
	return transport
}

func (x *Transport) init(ctx *context.Ctx) {
	prometheusAddr := fmt.Sprintf("%v:%v", ctx.Config.GateWay.Prometheus.Ip, ctx.Config.GateWay.Prometheus.Port)
	x.NodeExporter = repository.NewNodeExporter(prometheusAddr)
}

// 服务节点探活，心跳机制
func (x *Transport) getPrefixPath(cluster string, ctx *context.Ctx) string {
	for _, v := range ctx.CfgInfo.GateWay.Routes {
		if v.Uri == cluster {
			return v.Path
		}
	}
	return ""
}

// 发送请求之前进行预过滤
func (x *Transport) DoPreFilters(r *gin.Context, ctx *context.Ctx) dto.ErrorCode {
	path := r.Request.URL.Path
	routeInfo := x.getRouteInfo(ctx.GateWay, path)
	x.mux.Lock()
	defer x.mux.Unlock()
	clusterName := fmt.Sprintf("%s/%s", types.ClusterPrefix, routeInfo.Uri)
	rsp := x.cluster.GetClusterInfo(clusterName)
	if rsp.Code >= http.StatusBadRequest {
		return rsp
	}
	cluster := rsp.Data.(*entity.ClusterEntity)
	for _, filter := range cluster.Filters {
		statusCode := filter.Pre(ctx, r.Request)
		if statusCode.Code != fasthttp.StatusOK {
			logger.Logger().Errorf("[%v] %v", filter.Name(), statusCode)
			return statusCode
		}
	}
	return errord.MakeOkRsp(nil)
}

// 分发http请求到服务后端
func (x *Transport) DestNode(c *gin.Context, ctx *context.Ctx) string {
	clusterName, destNode := x.selectNode(c, ctx)
	if len(destNode) == 0 {
		logger.Logger().Errorf("clusterName:%v, have no health server!", clusterName)
		errord.InternalServerError(types.NoHealthService, c)
	}
	logger.Logger().Infof("node:%v, is selected!", destNode)
	return fmt.Sprintf("http://%s", destNode)
}

// 同步任务选择健康的服务后端
func (x *Transport) selectNode(req *gin.Context, ctx *context.Ctx) (string, string) {
	routeInfo := x.getRouteInfo(ctx.GateWay, req.Request.URL.Path)
	logger.Logger().Debugf("Select Cluster (%v)", routeInfo.Uri)
	x.mux.Lock()
	defer x.mux.Unlock()
	clusterName := fmt.Sprintf("%s/%s", types.ClusterPrefix, routeInfo.Uri)
	rsp := x.cluster.GetClusterInfo(clusterName)
	if rsp.Code >= http.StatusBadRequest {
		return clusterName, ""
	}
	cluster := rsp.Data.(*entity.ClusterEntity)
	rsp = x.server.GetBindInfo(routeInfo.Uri)
	if rsp.Code >= http.StatusBadRequest {
		return clusterName, ""
	}
	servers := rsp.Data.([]*entity.ServiceEntity)
	healthServers := x.getHealthServer(routeInfo.Limits, servers, ctx.CfgInfo.GateWay.BlackList, ctx.GateWay)
	if cluster.Lb == nil {
		logger.Logger().Errorf("%v, lb is nil", clusterName)
		return clusterName, ""
	}
	return clusterName, cluster.Lb.Select(healthServers)
}

func (x *Transport) getRouteInfo(cfg config.GateWay, path string) config.Route {
	for _, v := range cfg.Routes {
		if strings.HasPrefix(path, v.Path) {
			return v
		}
	}
	return config.Route{}
}

func (x *Transport) getHealthServer(taskLimit config.Limit, meta []*entity.ServiceEntity, blackList []string, way config.GateWay) []*vo.LbVO {
	blackLists := strings.Join(blackList, ",")
	servers := make([]*vo.LbVO, 0)
	for _, v := range meta {
		if v.Health {
			if !x.checkNodeLimit(v.Addr, taskLimit, way.Prometheus) {
				logger.Logger().Errorf("%v , resource is not matched limit", v.Addr)
				continue
			}
			if !strings.HasSuffix(blackLists, v.Addr) {
				servers = append(servers, &vo.LbVO{
					Addr:   v.Addr,
					Weight: v.Weight,
				})
			} else {
				logger.Logger().Debugf("[%v] in black_list", v.Addr)
			}
		}
		logger.Logger().Debugf("[%v] is %v", v.Addr, v.CurStatus)
	}
	return servers
}

func (x *Transport) checkNodeLimit(node string, limit config.Limit, prometheus config.Prometheus) bool {
	if prometheus.Switch == types.PrometheusSwitchOff {
		return true
	}
	curNodeResource, err := x.CurNodeLeftResource(node, prometheus.Mount)
	if err != nil {
		return false
	}
	if curNodeResource.DiskLeft < limit.Disk ||
		curNodeResource.MemoryAvail < limit.Memory ||
		curNodeResource.CpuIdle < limit.Cpu {
		logger.Logger().Errorf("%v, {%v} is not matched limit {%v}", node, curNodeResource, limit)
		return false
	}
	return true
}

func (x *Transport) Close() {

}
