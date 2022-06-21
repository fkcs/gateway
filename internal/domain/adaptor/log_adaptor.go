package adaptor

import "github.com/fkcs/gateway/internal/interfaces/dto"

type LogDomainInterface interface {
	SetLevel(level string) dto.ErrorCode
	GetLevel() dto.ErrorCode
}
