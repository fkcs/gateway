package lb

import (
	"github.com/fkcs/gateway/internal/domain/proxy/filter/adapter"
	"github.com/fkcs/gateway/internal/domain/vo"
	"github.com/fkcs/gateway/internal/utils/types"
	"fmt"
	"reflect"
)

func init() {
	Name := "demo"
	adapter.RegisterEngine(types.EngineTypeLb,
		Name,
		reflect.TypeOf((*DemoBalance)(nil)))
}

type DemoBalance struct {
}

func NewRandBalance() adapter.LoadBalance {
	lb := DemoBalance{}
	return lb
}

func (rb DemoBalance) Init() {
	fmt.Println("demo init success")
}

func (rb DemoBalance) Select(servers []*vo.LbVO) string {
	fmt.Println("run demo balance")
	return ""
}
