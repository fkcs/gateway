package config

import (
	path2 "github.com/fkcs/gateway/internal/utils/path"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"strings"
)

type Config struct {
	Addr       string
	Port       int
	MetricPort int
	*CfgInfo
}

type CfgInfo struct {
	GateWay GateWay `yaml:"gate_way"`
}

type GateWay struct {
	Routes      []Route       `yaml:"routes"`
	Tasks       []Task        `yaml:"tasks"`
	RateLimits  RateLimit     `yaml:"rate_limit"`
	HystrixInfo Hystrix       `yaml:"hystrix"`
	HealthCheck HealthCheck   `yaml:"health_check"`
	Queue       QueueCfg      `yaml:"queue"`
	Notify      NotifyCfg     `yaml:"notify"`
	BlackList   []string      `yaml:"black_list"`
	Prometheus  Prometheus    `yaml:"prometheus"` // Prometheus配置信息
	Aggregation []Aggregation `yaml:"aggregation_api"`
}

type NotifyCfg struct {
	User   string `yaml:"user"`
	Passwd string `yaml:"passwd"`
	Host   string `yaml:"host"`
	Port   int    `yaml:"port"`
}

type QueueCfg struct {
	Try int `yaml:"try"` // 重试任务个数
	TTL int `yaml:"ttl"` // 多久监听任务是否过期
}

type HealthCheck struct {
	Timeout         int `yaml:"timeout"`
	Internal        int `yaml:"internal"`
	ServiceAfterDel int `yaml:"service_after_del"`
}

type Route struct {
	Id          string   `yaml:"id"`
	Uri         string   `yaml:"cluster"`
	LoadBalance string   `yaml:"lb"`
	Path        string   `yaml:"path"`
	Filters     []Filter `yaml:"filters"`
	Limits      Limit    `yaml:"limit"`
}

type Filter struct {
	Name string                 `yaml:"name"`
	Args map[string]interface{} `yaml:"args"`
}

type Task struct {
	Name    string `yaml:"event_type"`
	Path    string `yaml:"path"`
	Uri     string `yaml:"cluster"`
	TimeOut int    `yaml:"timeout"`
	DoLen   int    `yaml:"do_len"`
	PendLen int    `yaml:"pend_len"`
	Limits  Limit  `yaml:"limit"`
}

type Limit struct {
	Cpu    float64 `yaml:"cpu"`
	Memory uint32  `yaml:"memory"`
	Disk   uint32  `yaml:"disk"`
}

type RateLimit struct {
	Capacity int64 `yaml:"capacity"`
}

type Hystrix struct {
	RequestThreshold    int64   `yaml:"request_threshold"`
	SleepWinTime        int64   `yaml:"sleep_win_time"`
	ErrThresholdPercent float64 `yaml:"err_threshold_percent"`
}

type Prometheus struct {
	Switch int    `yaml:"switch"`
	Ip     string `yaml:"ip"`
	Port   int    `yaml:"port"`
	Mount  string `yaml:"mount"`
}

type Aggregation struct {
	Type string      `yaml:"type"`
	Url  string      `yaml:"url" `
	Apis interface{} `yaml:"apis"`
}

// TODO err待优化优雅退出
func NewCfgInfo(path string) *CfgInfo {
	if !strings.HasPrefix(path, "/") {
		absPath := path2.AbsPath(path)
		path = absPath
	}
	file, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}
	cfg := &CfgInfo{}
	err = yaml.Unmarshal(file, cfg)
	if err != nil {
		panic(err)
	}
	return cfg
}
