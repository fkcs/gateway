package common

import (
	"sync/atomic"
)

type OnceChan struct {
	Channel chan interface{}
	wrote   int32
}

func NewOnceChan() *OnceChan {
	return &OnceChan{
		Channel: make(chan interface{}),
		wrote:   0,
	}
}

func (oc *OnceChan) IsNil(err interface{}) bool {
	if err == nil {
		return true
	}
	if atomic.AddInt32(&oc.wrote, 1) > 1 {
		return false
	}
	oc.Channel <- err
	return false
}
