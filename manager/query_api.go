package manager

import (
	"github.com/coreos/etcd/clientv3"
	"sync"
	"simpleGatewayExample/sdk"
	"simpleGatewayExample/global/apilist"
	"go.uber.org/zap"
	"github.com/valyala/fasthttp"
)

var etcdCli *clientv3.Client
var Servers = &sync.Map{}
var httpCli = &fasthttp.Client{}


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
//更新gateway 中服务器列表
func updateApi(apiName string, tp int) {
	args := &fasthttp.Args{}
	args.Set("api_name", apiName)
	args.SetUint("type", tp)

	serversS, ok := Servers.Load(apilist.SimpleGatewayExampleGatewayUpdateApi)
	if !ok {
		Logger.Warn("gateway没有节点存活")
		return
	}
	servers := serversS.(*sdk.QueryServerRes)

	for _, server := range servers.Servers {
		url := "http://" + server.IP + server.Path
		code, _, err := httpCli.Post(nil, url, args)
		if err != nil {
			Logger.Warn("manager update api error", zap.Error(err), zap.String("url", url), zap.String("args", args.String()))
			continue
		}

		if code != 200 {
			Logger.Warn("manager update api code invalid", zap.Int("code", code), zap.String("url", url), zap.String("args", args.String()))
			continue
		}
	}
}
