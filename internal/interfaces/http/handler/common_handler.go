package handler

import (
	"github.com/fkcs/gateway/internal/application/adaptor"
	"github.com/fkcs/gateway/internal/application/service"
	"github.com/fkcs/gateway/internal/context"
	"github.com/fkcs/gateway/internal/infrastructure/logger"
	"github.com/fkcs/gateway/internal/interfaces/command"
	errord "github.com/fkcs/gateway/internal/utils/error"
	"github.com/gin-gonic/gin"
)

type ApplicationService struct {
	ctx *context.Ctx
	app adaptor.CommonInterface
}

func NewApplicationService() *ApplicationService {
	return &ApplicationService{}
}

func (x *ApplicationService) Init(ctx *context.Ctx) error {
	x.ctx = ctx
	x.app = service.NewCommonAppImpl(ctx)
	return nil
}

func (x *ApplicationService) Route(router *gin.RouterGroup) error {
	root := router.Group("/v1")
	{
		root.GET("/log/level", x.getLogLevel)
		root.POST("/log/:level", x.setLogLevel)

		root.GET("/clusters/server", x.getServer)
		root.GET("/clusters/servers", x.getServers)
		root.GET("/clusters", x.getCluster)
		root.POST("/:cluster/:server", x.addServer)
		root.DELETE("/:cluster/:server", x.deleteServer)
	}
	return nil
}

func (x *ApplicationService) setLogLevel(c *gin.Context) {
	level := c.Param("level")
	logger.Logger().Infof("set log level is %v", level)
	rsp := x.app.SetLogLevel(level)
	errord.Response(rsp, c)
}

func (x *ApplicationService) getLogLevel(c *gin.Context) {
	rsp := x.app.GetLogLevel()
	errord.Response(rsp, c)
}

func (x *ApplicationService) getServer(c *gin.Context) {
	addr := c.Query("addr")
	rsp := x.app.GetServer(addr)
	errord.Response(rsp, c)
}

func (x *ApplicationService) getServers(c *gin.Context) {
	rsp := x.app.GetServers()
	errord.Response(rsp, c)
}

func (x *ApplicationService) addServer(c *gin.Context) {
	info := command.ServerInfo{
		Cluster: c.Param("cluster"),
		Addr:    c.Param("server"),
	}
	rsp := x.app.AddServer(&info)
	errord.Response(rsp, c)
}

func (x *ApplicationService) deleteServer(c *gin.Context) {
	info := command.ServerInfo{
		Cluster: c.Param("cluster"),
		Addr:    c.Param("server"),
	}
	rsp := x.app.DeleteServer(&info)
	errord.Response(rsp, c)
}

func (x *ApplicationService) getCluster(c *gin.Context) {
	cluster := c.Query("name")
	rsp := x.app.GetCluster(cluster)
	errord.Response(rsp, c)
}
