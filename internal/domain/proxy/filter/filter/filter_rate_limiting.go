package filter

import (
	"github.com/fkcs/gateway/internal/context"
	"github.com/fkcs/gateway/internal/domain/proxy/filter/adapter"
	"github.com/fkcs/gateway/internal/infrastructure/logger"
	"github.com/fkcs/gateway/internal/interfaces/dto"
	error2 "github.com/fkcs/gateway/internal/utils/error"
	"github.com/fkcs/gateway/internal/utils/types"
	"fmt"
	"github.com/juju/ratelimit"
	"net/http"
	"time"
)

type RateLimitingFilter struct {
	duration int64
	capacity int64
	limiter  *ratelimit.Bucket
	adapter.BaseFilter
}

func NewRateLimitingFilter() adapter.Filter {
	return &RateLimitingFilter{
		limiter: ratelimit.NewBucket(types.DefaultDuration/time.Duration(types.DefaultQPS), types.DefaultQPS),
	}
}

func (f *RateLimitingFilter) Init(args map[string]interface{}) error {
	if len(args) == 0 {
		return nil
	}
	logger.Logger().Infof("[ARGS] %v", args)
	capacity, ok := args["capacity"].(int)
	if !ok {
		return fmt.Errorf("capacity is not int")
	}
	f.capacity = int64(capacity)
	duration, ok := args["duration"].(int)
	if !ok {
		return fmt.Errorf("durarion is not int")
	}
	f.duration = int64(duration)
	logger.Logger().Infof("[%s] init capacity=%d,duration=%d", f.Name(), capacity, duration)
	f.limiter = ratelimit.NewBucket(time.Duration(duration)*time.Second/time.Duration(capacity), int64(capacity))
	return nil
}

func (f *RateLimitingFilter) Name() string {
	return types.FilterRateLimiting
}

func (f *RateLimitingFilter) Pre(ctx *context.Ctx, req *http.Request) dto.ErrorCode {
	if !f.do(ctx) {
		return error2.MakeLimited(types.RateLimitErr)
	}
	return f.BaseFilter.Pre(ctx, req)
}

// 令牌桶实现，此处用interface方式实现，依赖注入方式
func (f *RateLimitingFilter) do(ctx *context.Ctx) bool {
	logger.Logger().Debugf("[capacity:%d] -> [limit:%d]", f.capacity, f.limiter.Available())
	return f.limiter.TakeAvailable(1) != 0
}
