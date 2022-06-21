package lb

import (
	"github.com/fkcs/gateway/internal/domain/proxy/filter/adapter"
	"github.com/fkcs/gateway/internal/utils/types"
	"reflect"
)

func LoadLBEng(name string) (adapter.LoadBalance, bool) {
	reflectType := adapter.GetEngine(types.EngineTypeLb, name)
	if reflectType == nil {
		return nil, false
	}
	return reflect.New(reflectType.Elem()).Interface().(adapter.LoadBalance), true
}
