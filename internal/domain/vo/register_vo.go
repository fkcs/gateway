package vo

type ServerVo struct {
	Cluster string `json:"cluster"`
	Addr    string `json:"addr"`
}

type LbVO struct {
	Addr   string
	Weight int64
}
