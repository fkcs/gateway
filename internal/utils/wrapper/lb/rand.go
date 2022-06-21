package lb

import (
	"github.com/fkcs/gateway/internal/domain/proxy/filter/adapter"
	"github.com/fkcs/gateway/internal/domain/vo"
	"github.com/fkcs/gateway/internal/utils/types"
	"github.com/valyala/fastrand"
	"reflect"
)

func init() {
	Name := types.RandLoadBalance
	adapter.RegisterEngine(types.EngineTypeLb,
		Name,
		reflect.TypeOf((*RandBalance)(nil)))
}

type RandBalance struct {
}

func NewRandBalance() adapter.LoadBalance {
	lb := RandBalance{}
	return lb
}

func (rb RandBalance) Init() {

}

func (rb RandBalance) Select(servers []*vo.LbVO) string {
	l := len(servers)
	if 0 >= l {
		return ""
	}
	server := servers[fastrand.Uint32n(uint32(l))]
	return server.Addr
}
