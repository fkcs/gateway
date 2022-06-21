package context

import (
	"github.com/fkcs/gateway/internal/infrastructure/memory"
	"github.com/fkcs/gateway/internal/infrastructure/watch"
	"github.com/fkcs/gateway/internal/utils/common"
	"github.com/fkcs/gateway/internal/utils/config"
	"context"
	"database/sql"
	"fmt"
	"go/types"
	"gorm.io/gorm"
	"sync"
)

// 上下文信息，无状态设计
type Ctx struct {
	*config.Config
	*gorm.DB
	SqlDB *sql.DB
	*watch.EtcdClient
	Redis   *memory.RedisClient
	Ctx     context.Context
	cancel  context.CancelFunc
	Pools   *sync.Pool // TODO 分布式锁
	Panic   *common.OnceChan
	Inotify chan bool
}

func NewCtx(repo *gorm.DB, sqlDB *sql.DB, etcdClient *watch.EtcdClient, redis *memory.RedisClient, cfg *config.Config) *Ctx {
	ctxInfo, cancel := context.WithCancel(context.Background())
	ctx := &Ctx{
		Ctx:        ctxInfo,
		cancel:     cancel,
		DB:         repo,
		SqlDB:      sqlDB,
		EtcdClient: etcdClient,
		Redis:      redis,
		Config:     cfg,
		Pools: &sync.Pool{
			New: func() interface{} {
				return new(types.Nil)
			},
		},
		Panic:   common.NewOnceChan(),
		Inotify: make(chan bool, 1),
	}
	return ctx
}

func (x *Ctx) Host() string {
	return fmt.Sprintf("%s:%d", x.Addr, x.Port)
}

func (x *Ctx) Close() {
	x.Redis.Close()
	x.EtcdClient.Client.Close()
	x.SqlDB.Close()
	x.cancel()
}
