package lb

import (
	"github.com/fkcs/gateway/internal/domain/proxy/filter/adapter"
	"github.com/fkcs/gateway/internal/domain/vo"
	"github.com/fkcs/gateway/internal/utils/types"
	"reflect"
	"sync/atomic"
)

func init() {
	Name := types.RoundRobinLoadBalance
	adapter.RegisterEngine(types.EngineTypeLb,
		Name,
		reflect.TypeOf((*RoundRobin)(nil)))
}

type RoundRobin struct {
	ops *uint64
}

func NewRoundRobin() adapter.LoadBalance {
	var ops uint64
	ops = 0
	return &RoundRobin{
		ops: &ops,
	}
}

func (x *RoundRobin) Init() {
	var ops uint64
	ops = 0
	x.ops = &ops
}

func (x *RoundRobin) Select(servers []*vo.LbVO) string {
	l := uint64(len(servers))
	if 0 >= l {
		return ""
	}
	target := servers[int(atomic.AddUint64(x.ops, 1)%l)]
	return target.Addr
}
