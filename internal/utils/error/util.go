package error

import (
	"github.com/fkcs/gateway/internal/interfaces/dto"
	"github.com/fkcs/gateway/internal/utils/types"
	"net/http"
)

func MakeOkRsp(data interface{}) dto.ErrorCode {
	return dto.ErrorCode{
		Data:      data,
		Code:      http.StatusOK,
		ErrorCode: types.OK,
	}
}

func MakeCustomErrorRequest(httpCode int, errorCode string) dto.ErrorCode {
	return dto.ErrorCode{
		Code:      httpCode,
		ErrorCode: errorCode,
	}
}

func MakeBadRequest(errorCode string) dto.ErrorCode {
	return dto.ErrorCode{
		Code:      http.StatusBadRequest,
		ErrorCode: errorCode,
	}
}

// Unauthorized generates a 401 error.
func MakeUnauthorized(errorCode string) dto.ErrorCode {
	return dto.ErrorCode{
		Code:      http.StatusUnauthorized,
		ErrorCode: errorCode,
	}
}

// Forbidden generates a 403 error.
func MakeForbidden(errorCode string) dto.ErrorCode {
	return dto.ErrorCode{
		Code:      http.StatusForbidden,
		ErrorCode: errorCode,
	}
}

// NotFound generates a 404 error.
func MakeNotFound(errorCode string) dto.ErrorCode {
	return dto.ErrorCode{
		Code:      http.StatusNotFound,
		ErrorCode: errorCode,
	}
}

// Timeout generates a 408 error.
func MakeTimeout(errorCode string) dto.ErrorCode {
	return dto.ErrorCode{
		Code:      http.StatusRequestTimeout,
		ErrorCode: errorCode,
	}
}

// Timeout generates a 408 error.
func MakeLimited(errorCode string) dto.ErrorCode {
	return dto.ErrorCode{
		Code:      http.StatusTooManyRequests,
		ErrorCode: errorCode,
	}
}

// InternalServerError generates a 500 error.
func MakeInternalServerError(errorCode string) dto.ErrorCode {
	return dto.ErrorCode{
		Code:      http.StatusInternalServerError,
		ErrorCode: errorCode,
	}
}

/*
func MakeOkRsp(req interface{}) (int, []byte) {
	code := dto.ErrorCode{
		Data:      req,
		Code:      http.StatusOK,
		ErrorCode: types.OK,
	}
	data, _ := json.Marshal(&code)
	return http.StatusOK, data
}

func MakeCustomErrorRequest(httpCode int, errorCode string) dto.ErrorCode {
	return dto.ErrorCode{
		Code:      httpCode,
		ErrorCode: errorCode,
	}
}

func MakeBadRequest(errorCode string) (int, []byte) {
	code := dto.ErrorCode{
		Code:      http.StatusBadRequest,
		ErrorCode: errorCode,
	}
	data, _ := json.Marshal(&code)
	return http.StatusBadRequest, data
}

// Unauthorized generates a 401 error.
func MakeUnauthorized(errorCode string) (int, []byte) {
	code := dto.ErrorCode{
		Code:      http.StatusUnauthorized,
		ErrorCode: errorCode,
	}
	data, _ := json.Marshal(&code)
	return http.StatusUnauthorized, data
}

// Forbidden generates a 403 error.
func MakeForbidden(errorCode string) (int, []byte) {
	code := dto.ErrorCode{
		Code:      http.StatusForbidden,
		ErrorCode: errorCode,
	}
	data, _ := json.Marshal(&code)
	return http.StatusForbidden, data
}

// NotFound generates a 404 error.
func MakeNotFound(errorCode string) (int, []byte) {
	code := dto.ErrorCode{
		Code:      http.StatusNotFound,
		ErrorCode: errorCode,
	}
	data, _ := json.Marshal(&code)
	return http.StatusNotFound, data
}

// Timeout generates a 408 error.
func MakeTimeout(errorCode string) (int, []byte) {
	code := dto.ErrorCode{
		Code:      http.StatusRequestTimeout,
		ErrorCode: errorCode,
	}
	data, _ := json.Marshal(&code)
	return http.StatusRequestTimeout, data
}

// Timeout generates a 408 error.
func MakeLimited(errorCode string) (int, []byte) {
	code := dto.ErrorCode{
		Code:      http.StatusTooManyRequests,
		ErrorCode: errorCode,
	}
	data, _ := json.Marshal(&code)
	return http.StatusTooManyRequests, data
}

// InternalServerError generates a 500 error.
func MakeInternalServerError(errorCode string) (int, []byte) {
	code := dto.ErrorCode{
		Code:      http.StatusInternalServerError,
		ErrorCode: errorCode,
	}
	data, _ := json.Marshal(&code)
	return http.StatusInternalServerError, data
}*/
