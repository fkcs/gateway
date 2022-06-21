package adapter

import (
	"github.com/fkcs/gateway/internal/domain/vo"
)

// 负载均衡实现
type LoadBalance interface {
	Init()
	Select(servers []*vo.LbVO) string
}
