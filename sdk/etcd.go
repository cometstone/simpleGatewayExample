package sdk

import (
	"go.uber.org/zap"
	"github.com/coreos/etcd/clientv3"
	"time"
)

func InitEtcd(addrs []string, l *zap.Logger) *clientv3.Client {
	cfg := clientv3.Config{
		Endpoints:   addrs,
		DialTimeout: 10 * time.Second,
	}

	cli, err := clientv3.New(cfg)
	if err != nil {
		l.Fatal("Etcd init error", zap.Error(err), zap.Any("etcd_addr", addrs))
	}

	return cli
}
