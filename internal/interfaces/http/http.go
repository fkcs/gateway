package http

import (
	"github.com/fkcs/gateway/internal/context"
	"github.com/fkcs/gateway/internal/infrastructure/logger"
	"github.com/fkcs/gateway/internal/interfaces/http/handler"
	"github.com/fkcs/gateway/internal/interfaces/http/middleware"
	"github.com/fkcs/gateway/internal/utils/types"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
	"net/http"
)

type HttpServer struct {
	ctx *context.Ctx
	r   *gin.Engine
	*middleware.Proxy
}

func NewHttpServer(ctx *context.Ctx) *HttpServer {
	gin.SetMode(gin.ReleaseMode)
	gin.DisableConsoleColor()
	httpServer := &HttpServer{
		r:     gin.Default(),
		ctx:   ctx,
		Proxy: middleware.NewProxy(ctx),
	}
	logger.Logger().Infof("success to init http")
	return httpServer
}

// swagger http://localhost:10101/swagger/index.html
func (x *HttpServer) Init() bool {
	url := ginSwagger.URL(fmt.Sprintf("http://%s/swagger/doc.json", x.ctx.Host()))
	x.r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))
	for _, v := range x.ctx.Config.GateWay.Routes {
		root := x.r.Group(v.Path)
		root.Use(middleware.CORSMiddleware())
		root.Use(x.Proxy.PreFilter)
		root.Use(x.Proxy.BackFilter)
		root.Any("/*action", x.Proxy.ReverseProxy)
	}
	root := x.r.Group(types.GatewayRelativePath)
	scheduler := NewScheduleJob(x.ctx, root)
	scheduler.RegisterRoute(handler.NewApplicationService())
	if err := scheduler.Observe(); !x.ctx.Panic.IsNil(err) {
		logger.Logger().Errorf("%v", err)
		return false
	}
	return true
}

func (x *HttpServer) startMetricHTTP() {
	http.Handle("/metrics", promhttp.Handler())
	addr := fmt.Sprintf(":%d", x.ctx.MetricPort)
	logger.Logger().Infof("Starting prometheus agent at %s", addr)
	if err := http.ListenAndServe(addr, nil); !x.ctx.Panic.IsNil(err) {
		logger.Logger().Errorf("%v", err)
	}
}

func (x *HttpServer) Run(port int) {
	logger.Logger().Infof("start listen tcp %v", x.ctx.Host())
	go x.startMetricHTTP()
	if err := x.r.Run(fmt.Sprintf(":%d", port)); !x.ctx.Panic.IsNil(err) {
		logger.Logger().Errorf("%v", err)
		return
	}
}

func (x *HttpServer) Stop() {
	x.ctx.Close()
	x.Proxy.Stop()
}
