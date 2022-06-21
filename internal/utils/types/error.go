package types

const (
	PackageLanguageCN = "zh_CN"
	PackageLanguageEN = "zh_EN"
	PackageLanguageTW = "zh_YW"
)

const (
	OK                 = "OK"
	ErrCodeNotExist    = "GATEWAY.CODE_NOT_FOUND"
	NoHealthService    = "GATEWAY.NO_HEALTH"
	ServerInternal     = "GATEWAY.SYSTEM_ERROR"
	InvalidUrl         = "GATEWAY.INVALID_URL"
	AggretationInvalid = "GATEWAY.AGGRERATION_INVALID"
	NoSessionIdErrCode = "GATEWAY.SESSION_NULL"
	RateLimitErr       = "GATEWAY.RATE_LIMIT"
	InterfaceInvalid   = "GATEWAY.INTERFACE_INVALIDED"
	InvalidLogLevel    = "INVALID_LOG_LEVEL"
	RedisInternal      = "GATEWAY.REDIS_INTERNAL"
	MySqlInternal      = "GATEWAY.MYSQL_INTERNAL"
	UnmarshalErr       = "GATEWAY.UNMARSHALL_EXCEPTION"
	InvalidParam       = "GATEWAY.INVALID_PARAM"
)

type PackageLanguage struct {
	MsgCN string `json:"msg_cn"`
	MsgTW string `json:"msg_tw"`
	MsgEN string `json:"msg_en"`
}
