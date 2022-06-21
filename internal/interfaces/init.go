package interfaces

import (
	context2 "github.com/fkcs/gateway/internal/context"
	"github.com/fkcs/gateway/internal/infrastructure/logger"
	"github.com/fkcs/gateway/internal/utils/common"
	"github.com/fkcs/gateway/internal/utils/config"
	"github.com/fkcs/gateway/internal/utils/path"
	"context"
	"fmt"
	"github.com/howeyc/fsnotify"
	"go.etcd.io/etcd/v3/clientv3"
	"io/ioutil"
	"plugin"
	"runtime/debug"
	"strings"
	"time"
)

type INotify struct {
	ctx *context2.Ctx
}

func NewINotify(ctx *context2.Ctx) *INotify {
	inotify := &INotify{
		ctx: ctx,
	}
	return inotify
}

// TODO 入参通过...options形式实现
func (x *INotify) Init(cfgFile string) {
	x.heartBeat()
	common.BatchGoSafe(func() {
		x.initPlugin()
	}, func() {
		x.inotifyPlugin()
	}, func() {
		x.inotifyConfig(cfgFile)
	})
	logger.Logger().Infof("success to init proxy")
}

func (x *INotify) initPlugin() {
	pluginPath := path.AbsPath("/plugin")
	fileInfoList, err := ioutil.ReadDir(pluginPath)
	if !x.ctx.Panic.IsNil(err) {
		logger.Logger().Errorf("%v", err)
		return
	}
	for i := range fileInfoList {
		if fileInfoList[i].IsDir() {
			logger.Logger().Warnf("%v, is dir", fileInfoList[i].Name())
			continue
		}
		pluginName := fileInfoList[i].Name()
		pluginFile := fmt.Sprintf("%v/%v", pluginPath, pluginName)
		logger.Logger().Infof("[Plugin] %v", pluginFile)
		_, err = plugin.Open(pluginFile)
		if !x.ctx.Panic.IsNil(err) {
			logger.Logger().Errorf("%v", err)
		}
	}
}

// 监听插件是否更新
func (x *INotify) inotifyPlugin() {

}

func (x *INotify) heartBeat() {
	leaseRsp, err := x.ctx.EtcdClient.Client.Grant(context.Background(), 60)
	if !x.ctx.Panic.IsNil(err) {
		logger.Logger().Errorf("%v", err)
		return
	}
	leaseID := leaseRsp.ID
	keepAliveContext, cancelFunc := context.WithCancel(context.Background())
	keepAliveChan, err := x.ctx.EtcdClient.Client.KeepAlive(keepAliveContext, leaseID)
	if !x.ctx.Panic.IsNil(err) {
		cancelFunc()
		return
	} else {
		x.lisKeepAlive(keepAliveContext, keepAliveChan)
	}
	key := fmt.Sprintf("/api/heartBeat/%v", x.ctx.Host())
	_, err = x.ctx.EtcdClient.Client.Put(context.Background(), key, "", clientv3.WithLease(leaseID))
	if !x.ctx.Panic.IsNil(err) {
		logger.Logger().Errorf("%v", err)
		return
	}
	return
}

func (x *INotify) lisKeepAlive(context context.Context, keepaliveRes <-chan *clientv3.LeaseKeepAliveResponse) {
	go func() {
		done := context.Done()
		for {
			select {
			case ret := <-keepaliveRes:
				if ret != nil {
					logger.Logger().Debugf("success to keep alive! %v", time.Now())
				}
			case <-done:
				return
			}
		}
	}()
}

func (x *INotify) restartInotify(cfgFile string) {
	go x.inotifyConfig(cfgFile)
}

func (x *INotify) inotifyConfig(cfgFile string) {
	defer func() {
		if err := recover(); err != nil {
			x.restartInotify(cfgFile)
			logger.Logger().Errorf("recover %v", err)
			logger.Logger().Errorf("%s", string(debug.Stack()))
		}
	}()

	watcher, err := fsnotify.NewWatcher()
	if !x.ctx.Panic.IsNil(err) {
		logger.Logger().Errorf("%v", err)
		return
	}
	done := make(chan bool, 1)
	go func() {
		defer func() {
			if err := recover(); err != nil {
				if err := recover(); err != nil {
					logger.Logger().Errorf("recover %v", err)
					logger.Logger().Errorf("%s", string(debug.Stack()))
				}
			}
		}()
		for {
			select {
			case ev, ok := <-watcher.Event:
				if !ok {
					logger.Logger().Infof("watch is close")
					return
				}
				logger.Logger().Debugf("inotify: %+v",ev)
				if ev.IsCreate() && (ev.Name == cfgFile || strings.Contains(ev.Name,"data")) {
					go func() {
						x.ctx.GateWay = config.NewCfgInfo(cfgFile).GateWay
						logger.Logger().Debugf("%v", x.ctx.GateWay)
						x.ctx.Inotify <- true

					}()
				}
			case err := <-watcher.Error:
				logger.Logger().Errorf("%v", err)
			case <-x.ctx.Ctx.Done():
				done <- true
				logger.Logger().Infof("inotify task is done!")
				return
			}
		}
	}()
	if !strings.HasPrefix(cfgFile, "/") {
		absPath := path.AbsPath(cfgFile)
		cfgFile = absPath
	}
	pathInfo := strings.Split(cfgFile, "/")
	cfgPath := strings.Join(pathInfo[:len(pathInfo)-1], "/")
	logger.Logger().Infof("inotify [%v]", cfgPath)
	if err = watcher.Watch(cfgPath); !x.ctx.Panic.IsNil(err) {
		logger.Logger().Errorf("%v", err)
		return
	}
	<-done
	watcher.Close()
	logger.Logger().Infof("close watcher")
}
