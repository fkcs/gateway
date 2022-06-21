package filter

import (
	"github.com/fkcs/gateway/internal/context"
	"github.com/fkcs/gateway/internal/domain/proxy/filter/adapter"
	"github.com/fkcs/gateway/internal/infrastructure/logger"
	"github.com/fkcs/gateway/internal/interfaces/dto"
	"github.com/fkcs/gateway/internal/utils/types"
	"net/http"
)

type LicenseValidity struct {
	adapter.BaseFilter
}

func NewLicenseValidity() *LicenseValidity {
	return &LicenseValidity{}
}

func (x *LicenseValidity) Name() string {
	return types.FilterLeaseValid
}

func (x *LicenseValidity) Init(args map[string]interface{}) error {
	logger.Logger().Infof("[%v] init", x.Name())
	return nil
}

func (x *LicenseValidity) Pre(ctx *context.Ctx, req *http.Request) dto.ErrorCode {
	return x.BaseFilter.Pre(ctx, req)
}
