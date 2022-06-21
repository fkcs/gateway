package types

import (
	"time"
)

const (
	MetricQueryRoute = "/api/v1/query"

	CapacityUnitMB = 1024 * 1024
	CapacityUnitGB = 1024 * 1024 * 1024

	PrometheusSwitchOn  = 1
	PrometheusSwitchOff = 0

	MemoryMaxLimit = 100  // 100G
	CpuMaxLimit    = 1    // 100%
	DiskMaxLimit   = 1000 // 1T
)

type VectorOpts struct {
	Namespace string
	Subsystem string
	Name      string
	Help      string
	Labels    []string
}

type RspInfo struct {
	Rsp  []byte
	Code int
	Cost time.Duration
}

// Prometheus采集磁盘占用大小
type Disk struct {
	Device     string `json:"device"`
	Fstype     string `json:"fstype"`
	Instance   string `json:"instance"`
	Job        string `json:"job"`
	Mountpoint string `json:"mountpoint"`
}

type DiskResult struct {
	Metric Disk          `json:"metric"`
	Value  []interface{} `json:"value"`
}

type MetricDisk struct {
	ResultType string       `json:"resultType"`
	Result     []DiskResult `json:"result"`
}

type NodeFilesystemAvail struct {
	Status string     `json:"status"`
	Data   MetricDisk `json:"data"`
}

// Prometheus采集CPU剩余大小
type CpuResult struct {
	Value []interface{} `json:"value"`
}

type MetricCpu struct {
	ResultType string      `json:"resultType"`
	Result     []CpuResult `json:"result"`
}

type NodeCpuSecondsTotal struct {
	Status string    `json:"status"`
	Data   MetricCpu `json:"data"`
}

// Prometheus采集内存剩余大小
type Memory struct {
	Instance string `json:"instance"`
	Job      string `json:"job"`
}

type MemoryResult struct {
	Metric Memory        `json:"metric"`
	Value  []interface{} `json:"value"`
}

type MetricMemory struct {
	ResultType string         `json:"resultType"`
	Result     []MemoryResult `json:"result"`
}

type NodeMemoryAvail struct {
	Status string       `json:"status"`
	Data   MetricMemory `json:"data"`
}

type CurNodeResource struct {
	CpuIdle     float64
	DiskLeft    uint32
	MemoryAvail uint32
}
