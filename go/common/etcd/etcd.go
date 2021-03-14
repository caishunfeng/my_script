package etcd

import (
	"github.com/coreos/etcd/clientv3"
	"sync"
)

type EtcdConfig struct {
	Endpoints []string
	Username  string
	Password  string
}

type EtcdClient struct {
	client *clientv3.Client
	lock   *sync.Mutex
}

func InitEtcdClient(cfg *EtcdConfig) (etcdClient *EtcdClient, err error) {
	etcdCfg := clientv3.Config{
		Endpoints: cfg.Endpoints,
		Username:  cfg.Username,
		Password:  cfg.Password,
	}
	client, err := clientv3.New(etcdCfg)
	if err != nil {
		return nil, err
	}
	etcdClient = &EtcdClient{
		client: client,
	}
	return
}
