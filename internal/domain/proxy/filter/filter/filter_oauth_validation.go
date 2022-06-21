package filter

import (
	"github.com/fkcs/gateway/internal/context"
	"github.com/fkcs/gateway/internal/domain/proxy/filter/adapter"
	"github.com/fkcs/gateway/internal/infrastructure/logger"
	"github.com/fkcs/gateway/internal/interfaces/dto"
	errord "github.com/fkcs/gateway/internal/utils/error"
	"github.com/fkcs/gateway/internal/utils/types"
	"net/http"
)

type HeaderValidation struct {
	WhitePath []string
	adapter.BaseFilter
}

func NewHeaderValidation() *HeaderValidation {
	return &HeaderValidation{
		WhitePath: make([]string, 0),
	}
}

func (x *HeaderValidation) Name() string {
	return types.FilterOAuthValid
}

func (x *HeaderValidation) Init(args map[string]interface{}) error {
	if len(args) == 0 {
		return nil
	}
	x.WhitePath = make([]string, 0)
	if paths, ok := args["white_path"].([]interface{}); ok {
		for _, path := range paths {
			x.WhitePath = append(x.WhitePath, path.(string))
		}
	}
	logger.Logger().Debugf("[%s] init white_path=%v", x.Name(), x.WhitePath)
	return nil
}

func (x *HeaderValidation) Pre(ctx *context.Ctx, req *http.Request) dto.ErrorCode {
	if x.isWhitePath(req.URL.Path) {
		return errord.MakeOkRsp(nil)
	}
	if x.isWebsocket(&req.Header) {
		return x.checkWsSessionValid(ctx, &req.Header)
	} else {
		return x.checkHttpSessionValid(ctx, &req.Header)
	}
}

func (x *HeaderValidation) isWebsocket(header *http.Header) bool {
	return string(header.Get(types.HeaderUpgrade)) == "websocket"
}

func (x *HeaderValidation) isDebugMode(header *http.Header) bool {
	debugMode := header.Get(types.HeaderNlpDebugMode)
	if len(debugMode) == 0 {
		return false
	}
	if string(debugMode) == "true" {
		return true
	}
	return false
}

func (x *HeaderValidation) isWhitePath(path string) bool {
	for _, v := range x.WhitePath {
		if v == path {
			return true
		}
	}
	return false
}

// 检查HTTP
func (x *HeaderValidation) checkHttpSessionValid(ctx *context.Ctx, header *http.Header) dto.ErrorCode {
	if x.isDebugMode(header) {
		return errord.MakeOkRsp(nil)
	}
	sessionIDBytes := header.Get(types.HeaderSessionId)
	return x.isExpired(ctx, string(sessionIDBytes))
}

// 校验Websocket
func (x *HeaderValidation) checkWsSessionValid(ctx *context.Ctx, header *http.Header) dto.ErrorCode {
	if x.isDebugMode(header) {
		return errord.MakeOkRsp(nil)
	}
	sessionIDBytes := header.Get(types.HeaderSecWsProtocol)
	return x.isExpired(ctx, string(sessionIDBytes))
}

// 校验登录是否过期
func (x *HeaderValidation) isExpired(ctx *context.Ctx, sessionID string) dto.ErrorCode {
	if len(sessionID) == 0 {
		return errord.MakeUnauthorized(types.NoSessionIdErrCode)
	}
	return errord.MakeOkRsp(nil)
}
