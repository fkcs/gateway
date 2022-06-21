package adapter

import (
	"github.com/fkcs/gateway/internal/context"
	"github.com/fkcs/gateway/internal/interfaces/dto"
	error2 "github.com/fkcs/gateway/internal/utils/error"
	"github.com/valyala/fasthttp"
	"net/http"
)

// 依赖注入实现
type Filter interface {
	Name() string
	Init(args map[string]interface{}) error
	Pre(ctx *context.Ctx, req *http.Request) dto.ErrorCode
	Post(ctx *context.Ctx) (statusCode int, err error)
}

type BaseFilter struct{}

func (f BaseFilter) Init(args map[string]interface{}) error {
	return nil
}

// Pre execute before proxy
func (f BaseFilter) Pre(ctx *context.Ctx, req *http.Request) dto.ErrorCode {
	return error2.MakeOkRsp(nil)
}

// Post execute after proxy
func (f BaseFilter) Post(ctx *context.Ctx) (statusCode int, err error) {
	return fasthttp.StatusOK, nil
}
