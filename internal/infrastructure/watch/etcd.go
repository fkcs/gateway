package watch

import (
	"github.com/fkcs/gateway/internal/infrastructure/logger"
	"context"
	etcd "go.etcd.io/etcd/v3/clientv3"
	"go.etcd.io/etcd/v3/mvcc/mvccpb"
	"time"
)

type EtcdClient struct {
	Auth   etcd.Config
	Client *etcd.Client
}

func NewEtcdClient(addr string) *EtcdClient {
	return &EtcdClient{
		Auth: etcd.Config{
			Endpoints:   []string{addr},
			DialTimeout: 5 * time.Second,
		},
	}
}

func (x *EtcdClient) TryConnect() error {
	client, err := etcd.New(x.Auth)
	if err != nil {
		return err
	}
	x.Client = client
	return nil
}

func (x *EtcdClient) MustConnect() {
	client, err := etcd.New(x.Auth)
	if err != nil {
		panic(err)
	}
	x.Client = client
}

func (x *EtcdClient) Get(key string) *mvccpb.KeyValue {
	rsp, err := x.Client.Get(context.Background(), key)
	if err != nil {
		logger.Logger().Errorf("%v", err)
		return nil
	}
	if len(rsp.Kvs) == 0 {
		return nil
	}
	return rsp.Kvs[0]
}
