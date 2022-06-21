package grpc

import (
	"google.golang.org/grpc"
	"time"
)

type Conn struct {
	Pool     *Pool
	conn     *grpc.ClientConn
	unusable bool // 凡是超过容量新建的，状态为true，则Put时候直接close
	time     time.Time
}

func (c *Conn) Close() error {
	if c.unusable {
		c.Pool.cur--
		if c.conn != nil {
			return c.conn.Close()
		}
		return nil
	}
	return c.Pool.Put(c)
}

func (c *Conn) Value() *grpc.ClientConn {
	return c.conn
}

func (x *Pool) WrapConn(conn *grpc.ClientConn, unusable bool, time time.Time) *Conn {
	return &Conn{
		Pool:     x,
		conn:     conn,
		unusable: unusable,
		time:     time,
	}
}
