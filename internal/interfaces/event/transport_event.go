package transport

import (
	"github.com/fkcs/gateway/internal/application/adaptor"
	"github.com/fkcs/gateway/internal/application/event"
	"github.com/fkcs/gateway/internal/context"
	"github.com/fkcs/gateway/internal/domain/vo"
	"github.com/fkcs/gateway/internal/infrastructure/cache"
	"github.com/fkcs/gateway/internal/infrastructure/logger"
	"github.com/fkcs/gateway/internal/interfaces/command"
	"github.com/fkcs/gateway/internal/utils/common"
	"github.com/fkcs/gateway/internal/utils/types"
	"encoding/json"
	"fmt"
	"go.etcd.io/etcd/v3/clientv3"
	"go.etcd.io/etcd/v3/mvcc/mvccpb"
	"net"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
)

type TransportEvent struct {
	ctx     *context.Ctx
	mux     *sync.RWMutex
	watchCh chan *types.Evt
	checkCh chan *command.ServerInfo
	event   adaptor.TransportInterface
}

func NewTransportEvent(ctx *context.Ctx) *TransportEvent {
	return &TransportEvent{
		ctx:     ctx,
		mux:     new(sync.RWMutex),
		watchCh: make(chan *types.Evt),
		checkCh: make(chan *command.ServerInfo, 128),
		event:   event.NewTransportEventImpl(ctx),
	}
}

func (x *TransportEvent) Init() {
	x.initServers()
	common.BatchGoSafe(func() {
		x.readyToReceiveWatchEvent(x.ctx)
	}, func() {
		x.doWatch(x.ctx)
	}, func() {
		x.doCheck(x.ctx)
	}, func() {
		x.watchClusterCfg(x.ctx)
	})
}

func (x *TransportEvent) initServers() {
	rsp := x.event.GetServers()
	if rsp.Code >= http.StatusBadRequest {
		panic(fmt.Errorf("init fail")) // TODO 待优化优雅退出
	}
	var serversInfo = rsp.Data.([]vo.ServerVo)
	for _, v := range serversInfo {
		server := &command.ServerInfo{
			Cluster: v.Cluster,
			Addr:    v.Addr,
		}
		x.event.AddServer(server)
		x.checkCh <- server
	}
}

func (x *TransportEvent) watchClusterCfg(ctx *context.Ctx) {
	for {
		select {
		case <-ctx.Inotify:
			x.event.UpdateCluster(ctx)
		case <-ctx.Ctx.Done():
			logger.Logger().Infof("finish to watch cluster config")
			return
		}
	}
}

// 监控事件 key: /nlp/gateway/{src}/{type}/{key}    value: {value}
func (x *TransportEvent) doWatch(ctx *context.Ctx) {
	rspCh := ctx.EtcdClient.Client.Watch(ctx.Ctx, types.WatchEventPrefix, clientv3.WithPrefix())
	for v := range rspCh {
		for _, ev := range v.Events {
			switch ev.Type {
			case mvccpb.PUT:
				keys := strings.Split(string(ev.Kv.Key), "/")
				logger.Logger().Debugf("watch: [%v] [%v]", string(ev.Kv.Key), string(ev.Kv.Value))
				if len(keys) < 6 {
					continue
				}
				typeVal, err := strconv.Atoi(keys[4])
				if err != nil {
					continue
				}
				x.watchCh <- &types.Evt{
					Src:   types.EvtSrc(keys[3]),
					Type:  types.EvtType(typeVal),
					Key:   string(ev.Kv.Key),
					Value: string(ev.Kv.Value),
				}
			}
		}
	}
}

// 实现cluster/server/bind之间的增删改查
func (x *TransportEvent) readyToReceiveWatchEvent(ctx *context.Ctx) {
	for {
		select {
		case event := <-x.watchCh:
			logger.Logger().Debugf("[EVENT] watch %v", event)
			if event.Src == types.EventSrcCluster {
				x.doClusterEvent(event)
			} else if event.Src == types.EventSrcServer {
				x.doServeEvent(event)
			} else if event.Src == types.EventSrcBind {
				x.doBindEvent(event)
			}
		case <-ctx.Ctx.Done():
			logger.Logger().Infof("finish to receive watch event!")
			return
		}
	}
}

// 起个线程实时监控是否有server产生，如果有判断是否存在，不存在则起线程进行探活
func (x *TransportEvent) doCheck(ctx *context.Ctx) {
	for {
		select {
		case addr := <-x.checkCh:
			go x.check(addr, ctx)
		case <-ctx.Ctx.Done():
			logger.Logger().Infof("finish to heartbeat")
			return
		}
	}
}

func (x *TransportEvent) ping(addr string, timeout time.Duration) bool {
	var (
		err  error
		conn net.Conn
	)
	time.Sleep(time.Second)
	logger.Logger().Debugf("ping %v", addr)
	if conn, err = net.DialTimeout("tcp", addr, timeout); err != nil {
		logger.Logger().Errorf("%v,%v", addr, err)
		return false
	}
	defer conn.Close()
	return true
}

func (x *TransportEvent) check(info *command.ServerInfo, ctx *context.Ctx) {
	logger.Logger().Debugf("check heartBeat: %v,%v", info.Cluster, info.Addr)
	for {
		serverInfo, _, err := cache.ServersCache.Get(info.Addr)
		if err != nil {
			logger.Logger().Warnf("%s had been deleted!", info)
			return
		}
		timeOut := time.Duration(ctx.CfgInfo.GateWay.HealthCheck.Timeout) * time.Second
		if x.ping(info.Addr, timeOut) {
			x.event.UpServer(info)
		} else {
			rsp := x.event.DownServer(info, ctx)
			if rsp.Code >= http.StatusBadRequest || !rsp.Data.(bool) {
				return
			}
		}
		logger.Logger().Debugf("check server info! %v", serverInfo)
		time.Sleep(time.Duration(ctx.CfgInfo.GateWay.HealthCheck.Internal) * time.Second)
	}
}

// Server事件处理增删改查
func (x *TransportEvent) doServeEvent(evt *types.Evt) {
	logger.Logger().Debugf("receive server event (%v)", evt)
	var svr command.ServerInfo
	if err := json.Unmarshal([]byte(evt.Value), &svr); err != nil {
		logger.Logger().Errorf("invalid event (%s)! %v", evt.Value, err)
		return
	}
	if len(svr.Addr) == 0 || len(svr.Cluster) == 0 {
		logger.Logger().Errorf("event(%s) hadn't cluster!", evt.Value)
		return
	}
	if evt.Type == types.EventTypeNew {
		x.addServer(&svr)
	} else if evt.Type == types.EventTypeDelete {
		x.event.DelServer(&svr)
	} else if evt.Type == types.EventTypeUpdate {
		x.event.UpdateServer(&svr)
	}
}

// 增加服务后端IP
func (x *TransportEvent) addServer(info *command.ServerInfo) {
	rsp := x.event.AddServer(info)
	logger.Logger().Infof("add server %v", rsp)
	if rsp.Code >= http.StatusBadRequest {
		return
	}
	isFirst := rsp.Data.(bool)
	if isFirst {
		x.checkCh <- &command.ServerInfo{
			Cluster: info.Cluster,
			Addr:    info.Addr,
		}
	}
}

// Cluster事件处理增删改查
func (x *TransportEvent) doClusterEvent(evt *types.Evt) {
	logger.Logger().Infof("receive cluster %v", evt)
}

// Bind事件处理增删改查
func (x *TransportEvent) doBindEvent(evt *types.Evt) {
	logger.Logger().Infof("bind event(%v)", evt)
}
