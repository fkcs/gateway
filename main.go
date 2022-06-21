package main

import (
	"github.com/fkcs/gateway/internal/context"
	_ "github.com/fkcs/gateway/internal/domain/proxy/filter/adapter"
	"github.com/fkcs/gateway/internal/infrastructure/logger"
	"github.com/fkcs/gateway/internal/infrastructure/memory"
	"github.com/fkcs/gateway/internal/infrastructure/persistence"
	"github.com/fkcs/gateway/internal/infrastructure/watch"
	"github.com/fkcs/gateway/internal/interfaces"
	transport "github.com/fkcs/gateway/internal/interfaces/event"
	"github.com/fkcs/gateway/internal/interfaces/http"
	"github.com/fkcs/gateway/internal/utils"
	"github.com/fkcs/gateway/internal/utils/common"
	"github.com/fkcs/gateway/internal/utils/config"
	"github.com/fkcs/gateway/internal/utils/path"
	_ "github.com/fkcs/gateway/internal/utils/wrapper/lb"
	"flag"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	path.SetProjectRootEnv("API_GATEWAY_ROOT")
	hostIP := flag.String("host_ip", "127.0.0.1", "host IP address")
	port := flag.Int("port", 10101, "port for http service")
	metricPort := flag.Int("metric_port", 10111, "port for collect metric service")
	cfgFilename := flag.String("cfg", "/opt/cfg.yaml", "yaml格式的配置文件")
	etcdEP := flag.String("etcd", "172.16.30.34:2379", "etcd end point")
	dbIP := flag.String("db_ip", "172.16.30.34", "mysql DB IP address")
	dbPort := flag.Int("db_port", 3306, "mysql DB port")
	dbUsername := flag.String("db_username", "root", "mysql DB username")
	dbPassword := flag.String("db_password", "uWXf87plmQGz8zMM", "mysql DB password")
	dbName := flag.String("db_name", "nlp", "mysql DB name")
	redisAddr := flag.String("redis_addr", "172.16.30.34:47001", "redis address")
	redisPassword := flag.String("redis_password", "AhspHJ2l0ychcves", "redis password")
	logFilename := flag.String("log", "/logs/api-gateway.log", "log filename")
	logLevel := flag.String("log_level", "DEBUG", "log level")
	flag.Parse()

	utils.FlagMustBePresent("api-gateway", "host_ip", hostIP)
	utils.FlagMustBePresent("api-gateway", "etcd", etcdEP)
	utils.FlagMustBePresent("api-gateway", "db_ip", dbIP)
	utils.FlagMustBePresent("api", "redis_addr", redisAddr)
	utils.FlagMustBePresent("api", "cfg", cfgFilename)
	logger.LogInit(*logLevel, *logFilename, false, true)
	if *cfgFilename == "" {
		logger.Logger().Errorf("missing -cfg argument, run '%s --help' for more info", "api-gateway")
		os.Exit(1)
	}

	etcdClient := watch.NewEtcdClient(*etcdEP)
	etcdClient.MustConnect()
	redisClient := memory.NewRedisClient(*redisAddr, *redisPassword)
	redisClient.MustConnect()
	sqlDB, db := persistence.NewDB(*dbUsername, *dbPassword, *dbIP, *dbName, *dbPort)
	cfgInfo := &config.Config{
		Addr:       *hostIP,
		Port:       *port,
		MetricPort: *metricPort,
		CfgInfo:    config.NewCfgInfo(*cfgFilename),
	}
	beanCtx := context.NewCtx(db, sqlDB, etcdClient, redisClient, cfgInfo)
	transport.NewTransportEvent(beanCtx).Init()
	proxy := http.NewHttpServer(beanCtx)
	common.BatchGoSafe(func() {
		interfaces.NewINotify(beanCtx).Init(*cfgFilename)
	}, func() {
		if proxy.Init() {
			proxy.Run(*port)
		}
	})
	waitStop(proxy, beanCtx.Panic)
}

func waitStop(p *http.HttpServer, onceChan *common.OnceChan) {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGHUP, syscall.SIGKILL, os.Interrupt, os.Kill)
	done := make(chan bool, 1)
	go func() {
		select {
		case sig := <-ch:
			logger.Logger().Infof("%v", sig)
		case err := <-onceChan.Channel:
			logger.Logger().Errorf("panic:%v", err)
		}
		p.Stop()
		done <- true
	}()
	<-done
	logger.Logger().Infof("done!")
	time.Sleep(time.Second)
}
