package grpc

import (
	"github.com/fkcs/gateway/internal/utils/types"
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"time"
)

// 配置参数
type Options struct {
	factory     func(address string) (*grpc.ClientConn, error)
	MaxIdle     int           // 最大空闲连接数
	MaxCap      int           // 最大连接数
	IdleTimeout time.Duration // 连接最大空闲时间，超过该时间将会关闭，避免空闲时连接EOF
	WaitTimeout time.Duration // 当连接满的时候，等待超时时间，创建新的连接或者返回错误
	mode        bool          // false: 不新建，true：进行新建
}

var DefaultOptions = Options{
	factory:     grpcDial,
	MaxIdle:     8,
	MaxCap:      64,
	IdleTimeout: 15 * time.Second,
	WaitTimeout: 5 * time.Second,
	mode:        true,
}

func grpcDial(address string) (*grpc.ClientConn, error) {
	ctx, cancel := context.WithTimeout(context.Background(), types.DialTimeout)
	defer cancel()
	return grpc.DialContext(ctx, address, grpc.WithInsecure(),
		grpc.WithBackoffMaxDelay(types.BackoffMaxDelay),
		grpc.WithInitialWindowSize(types.InitialWindowSize),
		grpc.WithInitialConnWindowSize(types.InitialConnWindowSize),
		grpc.WithDefaultCallOptions(grpc.MaxCallSendMsgSize(types.MaxSendMsgSize)),
		grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(types.MaxRecvMsgSize)),
		grpc.WithKeepaliveParams(keepalive.ClientParameters{
			Time:                types.KeepAliveTime,
			Timeout:             types.KeepAliveTimeout,
			PermitWithoutStream: true,
		}))
}
