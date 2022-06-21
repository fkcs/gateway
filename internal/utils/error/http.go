package error

import (
	"github.com/fkcs/gateway/internal/interfaces/dto"
	"github.com/fkcs/gateway/internal/utils/types"
	"github.com/gin-gonic/gin"
	"net/http"
	"sync"
)

var mux *sync.RWMutex

func init() {
	mux = new(sync.RWMutex)
}

func OkRsp(ctx *gin.Context, data interface{}) {
	res := dto.Response{
		Code:    http.StatusOK,
		Message: "success",
		Data:    data,
		Status:  http.StatusText(http.StatusOK),
	}
	ctx.JSON(http.StatusOK, res)
}

func NullOkRsp(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, dto.Response{
		Code:    http.StatusOK,
		Message: "success",
		Status:  http.StatusText(http.StatusOK),
	})
}

func message(errorCode string, ctx *gin.Context) string {
	mux.RLock()
	defer mux.RUnlock()
	languageType := ctx.GetHeader("X-Package-Language")
	if msg, ok := errorCodeMap[errorCode]; ok {
		return language(languageType, msg)
	} else {
		return language(languageType, errorCodeMap[types.ErrCodeNotExist])
	}
}

func language(languageType string, packageLanguage *types.PackageLanguage) string {
	switch languageType {
	case types.PackageLanguageCN:
		return packageLanguage.MsgCN
	case types.PackageLanguageEN:
		return packageLanguage.MsgEN
	case types.PackageLanguageTW:
		return packageLanguage.MsgTW
	default:
		return packageLanguage.MsgCN
	}
}

func Response(code dto.ErrorCode, ctx *gin.Context) {
	if code.Code >= http.StatusBadRequest {
		CustomErrorRequest(code, ctx)
	} else {
		OkRsp(ctx, code.Data)
	}
}

func CustomErrorRequest(code dto.ErrorCode, ctx *gin.Context) {
	ctx.JSON(code.Code, dto.ResponseError{
		Code:      code.Code,
		ErrorCode: code.ErrorCode,
		Message:   message(code.ErrorCode, ctx),
	})
	ctx.Abort()
	return
}

func BadRequest(errorCode string, ctx *gin.Context) {
	ctx.JSON(http.StatusBadRequest, dto.ResponseError{
		Code:      http.StatusBadRequest,
		ErrorCode: errorCode,
		Message:   message(errorCode, ctx),
	})
	ctx.Abort()
	return
}

// Unauthorized generates a 401 error.
func Unauthorized(errorCode string, ctx *gin.Context) {
	ctx.JSON(http.StatusUnauthorized, dto.ResponseError{
		Code:      http.StatusUnauthorized,
		ErrorCode: errorCode,
		Message:   message(errorCode, ctx),
	})
	ctx.Abort()
}

// Forbidden generates a 403 error.
func Forbidden(errorCode string, ctx *gin.Context) {
	ctx.JSON(http.StatusForbidden, dto.ResponseError{
		Code:      http.StatusForbidden,
		ErrorCode: errorCode,
		Message:   message(errorCode, ctx),
	})
	ctx.Abort()
}

// NotFound generates a 404 error.
func NotFound(errorCode string, ctx *gin.Context) {
	ctx.JSON(http.StatusNotFound, dto.ResponseError{
		Code:      http.StatusNotFound,
		ErrorCode: errorCode,
		Message:   message(errorCode, ctx),
	})
	ctx.Abort()
}

// Timeout generates a 408 error.
func Timeout(errorCode string, ctx *gin.Context) {
	ctx.JSON(http.StatusRequestTimeout, dto.ResponseError{
		Code:      http.StatusRequestTimeout,
		ErrorCode: errorCode,
		Message:   message(errorCode, ctx),
	})
	ctx.Abort()
}

// Timeout generates a 408 error.
func Limited(errorCode string, ctx *gin.Context) {
	ctx.JSON(http.StatusTooManyRequests, dto.ResponseError{
		Code:      http.StatusTooManyRequests,
		ErrorCode: errorCode,
		Message:   message(errorCode, ctx),
	})
	ctx.Abort()
}

// InternalServerError generates a 500 error.
func InternalServerError(errorCode string, ctx *gin.Context) {
	ctx.JSON(http.StatusInternalServerError, dto.ResponseError{
		Code:      http.StatusInternalServerError,
		ErrorCode: errorCode,
		Message:   message(errorCode, ctx),
	})
	ctx.Abort()
}
