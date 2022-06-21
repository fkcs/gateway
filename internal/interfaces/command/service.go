package command

type ServerInfo struct {
	Cluster string `json:"cluster"`
	Addr    string `json:"addr"`
}

type BindInfo struct {
	ClusterName string `json:"cluster_name"`
	Addr        string `json:"addr"`
}

type MenuListTable struct {
}
