package error

import (
	"github.com/fkcs/gateway/internal/utils/types"
)

var errorCodeMap = map[string]*types.PackageLanguage{
	types.ErrCodeNotExist: {
		MsgCN: "参数错误",
		MsgEN: "Param error",
		MsgTW: "参数错误",
	},
	types.NoHealthService: {
		MsgCN: "系统异常",
		MsgEN: "System error",
		MsgTW: "系统异常",
	},
	types.ServerInternal: {
		MsgCN: "系统异常",
		MsgEN: "System error",
		MsgTW: "系统异常",
	},
	types.InvalidUrl: {
		MsgCN: "无效路由",
		MsgEN: "Invalid Url",
		MsgTW: "无效路由",
	},
	types.AggretationInvalid: {
		MsgCN: "无效参数",
		MsgEN: "Invalid Type",
		MsgTW: "无效参数",
	},
	types.NoSessionIdErrCode: {
		MsgCN: "认证失效",
		MsgEN: "Invalid Token",
		MsgTW: "认证失效",
	},
	types.RateLimitErr: {
		MsgCN: "达到流控上限",
		MsgEN: "Rate Limit",
		MsgTW: "达到流控上限",
	},
	types.InterfaceInvalid: {
		MsgCN: "接口已失效",
		MsgEN: "Interface invalid",
		MsgTW: "接口已失效",
	},
	types.InvalidLogLevel: {
		MsgCN: "无效日志等级",
		MsgEN: "Invalid Log Level",
		MsgTW: "无效日志等级",
	},
	types.RedisInternal: {
		MsgCN: "Redis服务异常",
		MsgEN: "Redis service internal",
		MsgTW: "Redis服务异常",
	},
	types.MySqlInternal: {
		MsgCN: "Mysql服务异常",
		MsgEN: "Mysql service internal",
		MsgTW: "Mysql服务异常",
	},
	types.UnmarshalErr: {
		MsgCN: "反序列数据异常",
		MsgEN: "Unmarshal exception",
		MsgTW: "反序列数据异常",
	},
	types.InvalidParam: {
		MsgCN: "无效参数",
		MsgEN: "Invalid Param",
		MsgTW: "无效参数",
	},
}
