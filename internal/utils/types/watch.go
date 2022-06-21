package types

type EvtType int
type EvtSrc string

const (
	EventSrcCluster = EvtSrc("cluster")
	EventSrcServer  = EvtSrc("server")
	EventSrcBind    = EvtSrc("bind")
)

const (
	EventTypeNew    = EvtType(0)
	EventTypeUpdate = EvtType(1)
	EventTypeDelete = EvtType(2)
)

const (
	ServerDownStatus = "DOWN"
	ServerUpStatus   = "UP"
)

const (
	WatchEventPrefix = "/api/gateway"
	ClusterPrefix    = "/api/cluster"
)

// Evt event
type Evt struct {
	Src   EvtSrc
	Type  EvtType
	Key   string
	Value string
}
