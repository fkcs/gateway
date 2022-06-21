package adapter

import (
	"fmt"
	"reflect"
	"sync"
)

var globalMux *sync.Mutex
var dpiEngine map[string]map[string]reflect.Type

func init() {
	globalMux = new(sync.Mutex)
	dpiEngine = make(map[string]map[string]reflect.Type)
}

func RegisterEngine(engType, name string, filter reflect.Type) {
	globalMux.Lock()
	defer globalMux.Unlock()
	fmt.Printf("%v,%v\n", engType, name)
	if engine, ok := dpiEngine[engType]; ok {
		if _, ok := engine[name]; ok {
			fmt.Printf("%v is existed\n", name)
			return
		}
		dpiEngine[engType][name] = filter
	} else {
		dpiEngine[engType] = make(map[string]reflect.Type)
		dpiEngine[engType][name] = filter
	}
}

func GetEngine(engType, name string) reflect.Type {
	globalMux.Lock()
	defer globalMux.Unlock()
	if engine, ok := dpiEngine[engType]; ok {
		return engine[name]
	}
	return nil
}
