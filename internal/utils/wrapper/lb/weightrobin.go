package lb

import (
	"github.com/fkcs/gateway/internal/domain/proxy/filter/adapter"
	"github.com/fkcs/gateway/internal/domain/vo"
	"github.com/fkcs/gateway/internal/utils/types"
	"reflect"
)

func init() {
	Name := types.WeightRobinLoadBalance
	adapter.RegisterEngine(types.EngineTypeLb,
		Name,
		reflect.TypeOf((*WeightRobin)(nil)))
}

type WeightRobin struct {
	opts map[uint64]*weightRobin
}

type weightRobin struct {
	effectiveWeight int64
	currentWeight   int64
}

func NewWeightRobin() adapter.LoadBalance {
	return &WeightRobin{
		opts: make(map[uint64]*weightRobin, 1024),
	}
}

func (w *WeightRobin) Init() {
	w.opts = make(map[uint64]*weightRobin, 1024)
}

func (w *WeightRobin) Select(servers []*vo.LbVO) string {
	var total int64
	var best uint64
	l := len(servers)
	if 0 >= l {
		return ""
	}

	for i := l - 1; i >= 0; i-- {
		svr := servers[i]
		if _, ok := w.opts[uint64(i)]; !ok {
			w.opts[uint64(i)] = &weightRobin{
				effectiveWeight: svr.Weight,
			}
		}

		wt := w.opts[uint64(i)]
		wt.currentWeight += wt.effectiveWeight
		total += wt.effectiveWeight

		if wt.effectiveWeight < svr.Weight {
			wt.effectiveWeight++
		}

		if best == 0 || w.opts[uint64(best)] == nil || wt.currentWeight > w.opts[best].currentWeight {
			best = uint64(i)
		}
	}

	if best == 0 {
		return ""
	}

	w.opts[best].currentWeight -= total

	return servers[best].Addr
}
