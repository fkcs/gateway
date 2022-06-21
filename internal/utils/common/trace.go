package common

import (
	"hash/fnv"
	"sync"
	"time"
)

const (
	workerBits  uint8 = 10
	numberBits  uint8 = 12
	workerMax   int64 = -1 ^ (-1 << workerBits)
	numberMax   int64 = -1 ^ (-1 << numberBits)
	timeShift   uint8 = workerBits + numberBits
	workerShift uint8 = numberBits
	epoch       int64 = 1525705533000
)

type Trace struct {
	mu        sync.Mutex // 添加互斥锁 确保并发安全
	timestamp int64      // 记录时间戳
	workerId  int64      // 该节点的ID
	number    int64      // 当前毫秒已经生成的id序列号(从0开始累加) 1毫秒内最多生成4096个ID
}

// 实例化一个工作节点
func NewWorker(addr string) *Trace {
	trace := &Trace{
		timestamp: 0,
		number:    0,
	}
	trace.workerId = int64(trace.hash(addr))
	return trace
}

func (w *Trace) GetId() int64 {
	w.mu.Lock()
	defer w.mu.Unlock()

	now := time.Now().UnixNano() / 1e6
	if w.timestamp == now {
		w.number++
		if w.number > numberMax {
			for now <= w.timestamp {
				now = time.Now().UnixNano() / 1e6
			}
		}
	} else {
		w.number = 0
		w.timestamp = now
	}
	ID := int64((now-epoch)<<timeShift | (w.workerId << workerShift) | (w.number))
	return ID
}

func (w *Trace) hash(s string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(s))
	return h.Sum32() % 1024
}
