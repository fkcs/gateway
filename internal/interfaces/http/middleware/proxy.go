// 资源初始化，前期数据过滤，以及推送
package middleware

import (
	"github.com/fkcs/gateway/internal/application/proxy/transport"
	context2 "github.com/fkcs/gateway/internal/context"
	"github.com/fkcs/gateway/internal/domain/proxy/endpoint"
	"github.com/fkcs/gateway/internal/infrastructure/logger"
	"github.com/fkcs/gateway/internal/interfaces/dto"
	"github.com/fkcs/gateway/internal/utils/common"
	errord "github.com/fkcs/gateway/internal/utils/error"
	"github.com/fkcs/gateway/internal/utils/wrapper/monitor"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type Proxy struct {
	*common.Trace
	ctx        *context2.Ctx
	dispatcher *transport.Transport
	endpoint   endpoint.EndPointer
}

func NewProxy(ctx *context2.Ctx) *Proxy {
	proxy := &Proxy{
		ctx:        ctx,
		dispatcher: transport.NewTransport(ctx),
		Trace:      common.NewWorker(ctx.Host()),
		endpoint:   endpoint.NewEndPoint(),
	}
	return proxy
}

func (x *Proxy) PreFilter(c *gin.Context) {
	var errorCode dto.ErrorCode
	startTime := time.Now()
	logger.Logger().Infof("PreFilter| START")
	errcode := x.filter(c)
	if errcode.Code > http.StatusBadRequest {
		errord.CustomErrorRequest(errorCode, c)
	}
	c.Next()
	endTime := time.Now()
	timeCost := endTime.Sub(startTime)
	monitor.ObserveMetric(timeCost, c.Request.URL.Path, errorCode.Code)
	logger.Logger().Infof("PreFilter| FINISH")
}

func (x *Proxy) BackFilter(c *gin.Context) {
	startTime := time.Now()
	c.Next()
	endTime := time.Now()
	timeCost := endTime.Sub(startTime)
	monitor.ObserveMetric(timeCost, c.Request.URL.Path, c.Writer.Status())
	logger.Logger().Infof("BackFilter| %v,%v", c.Request.URL.Path, c.Writer.Status())
}

func (x *Proxy) filter(r *gin.Context) dto.ErrorCode {
	return x.dispatcher.DoPreFilters(r, x.ctx)
}

func (x *Proxy) ForwardProxy(c *gin.Context) {
	x.endpoint.ForwardProxy(c.Writer, c.Request)
}

func (x *Proxy) ReverseProxy(c *gin.Context) {
	destNode := x.dispatcher.DestNode(c, x.ctx)
	x.endpoint.ReverseProxy(c.Writer, c.Request, destNode)
}

func (x *Proxy) Stop() {
	logger.Logger().Infof("stop proxy")
}
