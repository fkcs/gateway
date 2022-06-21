package filter

import (
	"github.com/fkcs/gateway/internal/context"
	"github.com/fkcs/gateway/internal/domain/proxy/filter/adapter"
	"github.com/fkcs/gateway/internal/infrastructure/logger"
	"github.com/fkcs/gateway/internal/interfaces/dto"
	"github.com/fkcs/gateway/internal/utils/types"
	"fmt"
	"net/http"
)

type MetricInfo struct {
	Metric string
	Apis   []string
}

type MetricThreshold struct {
	Module string
	Metric *MetricInfo
	adapter.BaseFilter
}

func NewMetricThreshold() *MetricThreshold {
	return &MetricThreshold{
		Metric: &MetricInfo{
			Apis: make([]string, 0),
		},
	}
}

func (x *MetricThreshold) Name() string {
	return types.FilterFlowThresh
}

func (x *MetricThreshold) Init(args map[string]interface{}) error {
	module, ok := args["module"].(string)
	if !ok {
		return fmt.Errorf("module is invalid")
	}
	x.Module = module
	metric, ok := args["metric"].(string)
	if !ok {
		return fmt.Errorf("metric is invalid")
	}
	x.Metric.Metric = metric
	if paths, ok := args["api"].([]interface{}); ok {
		for _, path := range paths {
			x.Metric.Apis = append(x.Metric.Apis, path.(string))
		}
	}
	logger.Logger().Infof("[%v] metric:%v, event_api:%v", x.Name(), x.Metric.Metric, x.Metric.Apis)
	return nil
}

func (x *MetricThreshold) Pre(ctx *context.Ctx, req *http.Request) dto.ErrorCode {
	isEvent := false
	path := req.URL.Path
	for _, v := range x.Metric.Apis {
		if v == path {
			isEvent = true
		}
	}
	if !isEvent {
		return x.BaseFilter.Pre(ctx, req)
	}
	// TODO 需要分布式锁，避免同时读写引发数据竞争，导致脏数据
	// TODO 首先校验是否超过阈值，redis存储
	// TODO 如果没有超过则计数+1，并且去重
	return x.BaseFilter.Pre(ctx, req)
}
