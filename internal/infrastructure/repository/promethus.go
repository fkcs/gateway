package repository

import (
	"github.com/fkcs/gateway/internal/infrastructure/logger"
	"github.com/fkcs/gateway/internal/utils/common"
	"github.com/fkcs/gateway/internal/utils/types"
	"encoding/json"
	"fmt"
	"github.com/valyala/fasthttp"
	"strconv"
	"strings"
)

type NodeExporter struct {
	Addr string
	*common.Method
}

func NewNodeExporter(addr string) *NodeExporter {
	return &NodeExporter{
		Addr:   addr,
		Method: common.NewMethod(),
	}
}

func (x *NodeExporter) query(query string) (int, []byte) {
	args := map[string]string{
		"query": query,
	}
	return x.Method.SetHost(x.Addr).Do("GET", types.MetricQueryRoute, nil, args, nil)
}

func (x *NodeExporter) MetricNodeCpuIdle(node string) (float64, error) {
	node = strings.Split(node, ":")[0]
	query := fmt.Sprintf(`sum(increase(node_cpu_seconds_total{instance="%s",mode="idle"}[1m])) `+
		`/ sum(increase(node_cpu_seconds_total{instance="%s"}[1m]))`, node, node)
	statusCode, rsp := x.query(query)
	if statusCode != fasthttp.StatusOK {
		return 0, fmt.Errorf("%v, failed to get cpu info", node)
	}
	logger.Logger().Debugf("%v,%v", node, rsp)
	var nodeCpuInfo types.NodeCpuSecondsTotal
	if err := json.Unmarshal([]byte(rsp), &nodeCpuInfo); err != nil {
		logger.Logger().Errorf("%v", err)
		return 0, err
	}
	results := nodeCpuInfo.Data.Result
	if len(results) != 0 {
		cpuIdleStr := fmt.Sprintf("%v", results[0].Value[1])
		cpuIdle, err := strconv.ParseFloat(cpuIdleStr, 64)
		if err != nil {
			logger.Logger().Errorf("failed to parse float!%v", err)
			return 0, err
		}
		logger.Logger().Debugf("[Metric] %v, CPU:%v", node, cpuIdle)
		return cpuIdle, nil
	}
	logger.Logger().Warnf("{%v}, no cpu idle info", node)
	return types.CpuMaxLimit, nil
}

func (x *NodeExporter) MetricNodeMemAvailable(node string) (uint32, error) {
	node = strings.Split(node, ":")[0]
	query := fmt.Sprintf(`node_memory_MemAvailable_bytes{instance="%s"}`, node)
	statusCode, rsp := x.query(query)
	if statusCode != fasthttp.StatusOK {
		return 0, fmt.Errorf("%v, failed to get mem info", node)
	}
	logger.Logger().Debugf("%v,%v", node, rsp)
	var nodeMemInfo types.NodeMemoryAvail
	if err := json.Unmarshal([]byte(rsp), &nodeMemInfo); err != nil {
		logger.Logger().Errorf("%v", err)
		return 0, err
	}
	results := nodeMemInfo.Data.Result
	if len(results) != 0 {
		memoryAvailStr := fmt.Sprintf("%v", results[0].Value[1])
		memoryAvail, err := strconv.ParseUint(memoryAvailStr, 10, 64)
		if err != nil {
			logger.Logger().Errorf("%v", err)
			return 0, err
		}
		memoryAvail = memoryAvail / types.CapacityUnitGB
		logger.Logger().Debugf("[Metric] %v, Memory:%v GB", node, memoryAvail)
		return uint32(memoryAvail), nil
	}
	logger.Logger().Warnf("{%v}, no memory info", node)
	return types.MemoryMaxLimit, nil
}

func (x *NodeExporter) MetricNodeDiskAvail(node string, device string) (uint32, error) {
	node = strings.Split(node, ":")[0]
	query := fmt.Sprintf(`node_filesystem_avail_bytes{instance="%s",mountpoint="%s"}`,
		node, device)
	statusCode, rsp := x.query(query)
	if statusCode != fasthttp.StatusOK {
		return 0, fmt.Errorf("%v, failed to get disk info", node)
	}
	logger.Logger().Debugf("%v,%v", node, rsp)
	var nodeDiskInfo types.NodeFilesystemAvail
	if err := json.Unmarshal([]byte(rsp), &nodeDiskInfo); err != nil {
		logger.Logger().Errorf("%v", err)
		return 0, err
	}
	results := nodeDiskInfo.Data.Result
	if len(results) != 0 {
		diskStr := fmt.Sprintf("%v", results[0].Value[1])
		disk, err := strconv.ParseUint(diskStr, 10, 64)
		if err != nil {
			logger.Logger().Errorf("%v", err)
			return 0, err
		}
		disk = disk / types.CapacityUnitGB
		logger.Logger().Debugf("[Metric] %v, Disk:%v GB, Mount:%v", node, disk, device)
		return uint32(disk), nil
	}
	logger.Logger().Warnf("{%v}, no disk info", node)
	return types.DiskMaxLimit, nil
}

func (x *NodeExporter) CurNodeLeftResource(node string, device string) (types.CurNodeResource, error) {
	cpuIdle, err := x.MetricNodeCpuIdle(node)
	if err != nil {
		logger.Logger().Errorf("%v", err)
		return types.CurNodeResource{}, err
	}
	memAvail, err := x.MetricNodeMemAvailable(node)
	if err != nil {
		logger.Logger().Errorf("%v", err)
		return types.CurNodeResource{}, err
	}
	diskAvail, err := x.MetricNodeDiskAvail(node, "/data")
	if err != nil {
		logger.Logger().Errorf("%v", err)
		return types.CurNodeResource{}, err
	}
	return types.CurNodeResource{
		CpuIdle:     cpuIdle,
		MemoryAvail: memAvail,
		DiskLeft:    diskAvail,
	}, nil
}
