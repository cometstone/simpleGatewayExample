package manager

import (
	"github.com/coreos/etcd/clientv3"
	"sync"
	"simpleGatewayExample/sdk"
	"simpleGatewayExample/global/apilist"
	"go.uber.org/zap"
)

var etcdCli *clientv3.Client
var Servers = &sync.Map{}


func initGatewayUpdate() {
	go func() {
		etcdCli = sdk.InitEtcd(Conf.Etcd.Addrs,Logger)
		resCh, errCh  := sdk.QueryServerByAPI(etcdCli,[]string{apilist.SimpleGatewayExampleGatewayUpdateApi, apilist.StaffCheckLogin},0)

		for {
			select {
			case servers := <- resCh:
				Servers.Store(servers.ApiName,servers.Servers)
				Logger.Debug("更新服务器列表",zap.Any("gateway_addr",servers))
			case err := <-errCh:
				Logger.Warn("请求etcd异常", zap.Error(err), zap.Any("etcd_addr", Conf.Etcd.Addrs))
			}
		}

	}()
}
