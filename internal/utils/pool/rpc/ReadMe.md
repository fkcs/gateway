## gRPC连接池

#### 配置结构体
``` go
type Options struct {
	factory     func(address string) (*grpc.ClientConn, error)
	MaxIdle     int           // 最大空闲连接数
	MaxCap      int           // 最大连接数
	IdleTimeout time.Duration // 连接最大空闲时间，超过该时间将会关闭，避免空闲时连接EOF
	WaitTimeout time.Duration // 当连接满的时候，等待超时时间，创建新的连接或者返回错误
	mode        bool          // false: 不新建，true：进行新建
}
```
创建gRPC连接工厂函数如下：
``` go
func grpcDial(address string) (*grpc.ClientConn, error) {
	ctx, cancel := context.WithTimeout(context.Background(), bean.DialTimeout)
	defer cancel()
	return grpc.DialContext(ctx, address, grpc.WithInsecure(),
		grpc.WithBackoffMaxDelay(bean.BackoffMaxDelay),
		grpc.WithInitialWindowSize(bean.InitialWindowSize),
		grpc.WithInitialConnWindowSize(bean.InitialConnWindowSize),
		grpc.WithDefaultCallOptions(grpc.MaxCallSendMsgSize(bean.MaxSendMsgSize)),
		grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(bean.MaxRecvMsgSize)),
		grpc.WithKeepaliveParams(keepalive.ClientParameters{
			Time:                bean.KeepAliveTime,
			Timeout:             bean.KeepAliveTimeout,
			PermitWithoutStream: true,
		}))
}
```
#### gRPC连接参数调优

+ MaxSendMsgSizeGRPC：最大允许发送的字节数，默认4MiB，如果超过了GRPC会报错。Client和Server我们都调到4GiB。

+ MaxRecvMsgSizeGRPC：最大允许接收的字节数，默认4MiB，如果超过了GRPC会报错。Client和Server我们都调到4GiB。

+ InitialWindowSize：基于Stream的滑动窗口，类似于TCP的滑动窗口，用来做流控，默认64KiB，吞吐量上不去，Client和Server我们调到1GiB。

+ InitialConnWindowSize：基于Connection的滑动窗口，默认16 * 64KiB，吞吐量上不去，Client和Server我们也都调到1GiB。

+ KeepAliveTime：每隔KeepAliveTime时间，发送PING帧测量最小往返时间，确定空闲连接是否仍然有效，我们设置为10S。

+ KeepAliveTimeout：超过KeepAliveTimeout，关闭连接，我们设置为3S。

+ PermitWithoutStream：如果为true，当连接空闲时仍然发送PING帧监测，如果为false，则不发送忽略。我们设置为true。

#### 示例
```
// 创建连接池
pool, err := grpc.NewPool(rpcAddr, &grpc.DefaultOptions)
if err != nil {
   return
}

// 获取连接
conn,err := pool.Get()
if err != nil {
   return
}
defer conn.Close()
// 建立连接
...

```
