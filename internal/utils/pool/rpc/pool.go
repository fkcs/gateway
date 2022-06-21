package grpc

import (
	"github.com/fkcs/gateway/internal/infrastructure/logger"
	"errors"
	"fmt"
	"sync"
	"time"
)

type Pool struct {
	mux     *sync.RWMutex
	cur     int
	conns   chan *Conn
	address string
	*Options
}

func NewPool(address string, options *Options) (*Pool, error) {
	if options.MaxIdle <= 0 || options.MaxCap <= 0 || options.MaxIdle > options.MaxCap {
		logger.Logger().Errorf("invalid option! %+v", options)
		return nil, errors.New("invalid option")
	}
	pool := &Pool{
		mux:     new(sync.RWMutex),
		conns:   make(chan *Conn, options.MaxCap),
		Options: options,
		address: address,
		cur:     0,
	}
	for i := 0; i < options.MaxIdle; i++ {
		conn, err := options.factory(address)
		if err != nil {
			pool.Close()
			logger.Logger().Errorf("%v", err)
			return nil, err
		}
		pool.conns <- pool.WrapConn(conn, false, time.Now())
	}
	return pool, nil
}

func (x *Pool) Get() (*Conn, error) {
	x.mux.Lock()
	defer x.mux.Unlock()

	if x.conns == nil {
		return nil, fmt.Errorf("when get grpc connection, pool is closed")
	}
	for {
		select {
		case conn := <-x.conns:
			if conn == nil || conn.conn == nil {
				return nil, fmt.Errorf("connection is nil")
			}
			if timeout := x.IdleTimeout; timeout > 0 {
				if conn.time.Add(timeout).Before(time.Now()) {
					conn.conn.Close()
					continue
				}
			}
			x.cur++
			return conn, nil
		case <-time.After(x.WaitTimeout):
			if x.cur < x.MaxCap {
				rpcConn, err := x.factory(x.address)
				if err != nil {
					logger.Logger().Errorf("%v", err)
					return nil, err
				}
				x.cur++
				conn := x.WrapConn(rpcConn, false, time.Now())
				return conn, nil
			} else {
				if x.mode { // 如果是允许新建连接模式，新建连接并
					rpcConn, err := x.factory(x.address)
					if err != nil {
						logger.Logger().Errorf("when create grpc connection, %v", err)
						return nil, err
					}
					x.cur++
					conn := x.WrapConn(rpcConn, true, time.Now())
					return conn, nil
				} else { // 如果是不允许新建连接模式，直接返回失败
					return nil, fmt.Errorf("when get grpc connection, pool is full")
				}
			}
		}
	}

}

func (x *Pool) Put(conn *Conn) error {
	if conn == nil {
		return fmt.Errorf("connection is nil")
	}
	x.mux.Lock()
	defer x.mux.Unlock()

	// 表示pool已经关闭
	if x.conns == nil {
		logger.Logger().Warnf("rpc pool had bean closed")
		return conn.Close()
	}
	conn.time = time.Now()
	x.cur--
	select {
	case x.conns <- conn:
		return nil
	default:
		logger.Logger().Warnf("pool is full, close the connection")
		return conn.conn.Close()
	}
}

func (x *Pool) Len() int {
	return len(x.conns)
}

func (x *Pool) Close() {
	x.mux.Lock()
	defer x.mux.Unlock()

	conns := x.conns
	x.conns = nil
	x.cur = 0
	if conns == nil {
		return
	}
	close(conns)
	for conn := range conns {
		conn.Close()
	}
}
