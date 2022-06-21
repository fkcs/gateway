package service

import (
	"github.com/fkcs/gateway/internal/infrastructure/logger"
	"github.com/fkcs/gateway/internal/interfaces/dto"
)

type LogDomainImpl struct {
}

func (x *LogDomainImpl) SetLevel(level string) dto.ErrorCode {
	return logger.SetLogLevel(level)
}

func (x *LogDomainImpl) GetLevel() dto.ErrorCode {
	return logger.LoadLogLevel()
}
