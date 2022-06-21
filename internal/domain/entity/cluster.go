package entity

import (
	"github.com/fkcs/gateway/internal/domain/proxy/filter/adapter"
	"github.com/fkcs/gateway/internal/utils/config"
)

type ClusterEntity struct {
	Name        string
	MaxQPS      int64
	RateLimit   int64
	LoadBalance string
	Lb          adapter.LoadBalance
	FilterCfg   map[string]config.Filter
	Filters     map[string]adapter.Filter
}

type ServiceEntity struct {
	Cluster     string
	Addr        string
	Health      bool
	Weight      int64
	PreStatus   string
	CurStatus   string
	HealthCheck int
}
