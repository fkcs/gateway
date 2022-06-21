package factory

import (
	"github.com/fkcs/gateway/internal/domain/proxy/filter/adapter"
	"github.com/fkcs/gateway/internal/domain/proxy/filter/filter"
	"github.com/fkcs/gateway/internal/utils/types"
)

func NewFilterFactory(name string) adapter.Filter {
	switch name {
	case types.FilterOAuthValid:
		return filter.NewHeaderValidation()
	case types.FilterRateLimiting:
		return filter.NewRateLimitingFilter()
	case types.FilterLeaseValid:
		return filter.NewLicenseValidity()
	case types.FilterFlowThresh:
		return filter.NewMetricThreshold()
	}
	return nil
}
