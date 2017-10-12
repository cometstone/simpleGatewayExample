package gateway

import (
	"github.com/coreos/etcd/clientv3"
	"time"
	"go.uber.org/zap"
)

var etcdCli *clientv3.Client

func initEtcd() {
	cfg := clientv3.Config{
		Endpoints:Conf.Etcd.Addrs,
		DialTimeout: 10*time.Second,
	}
	var err error
	etcdCli,err = clientv3.New(cfg)

	if err != nil {
		Logger.Fatal("init etcd err:",zap.Error(err))
	}



}