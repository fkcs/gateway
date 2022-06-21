// 监听ETCD的信息
package repository

// 有没有必要interface
type EtcdStore struct {
}

func NewEtcdSore() *EtcdStore {
	return &EtcdStore{}
}

func (x *EtcdStore) WatchService() {
	return
}

func (x *EtcdStore) Close() {

}
