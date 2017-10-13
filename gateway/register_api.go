package gateway

import (
	"github.com/labstack/echo"
	"simpleGatewayExample/global/gdata"
	"simpleGatewayExample/global/apilist"
	"simpleGatewayExample/sdk"
	"go.uber.org/zap"
)

func initUpdateApi() {

	go func() {
		e := echo.New()
		e.POST("/api/update",apiUpdate)
		e.Logger.Fatal(e.Start(":"+Conf.Api.ApiUpdatePort))
	}()

	go registerApiUpdate()
}

func apiUpdate(c echo.Context) error {
	return nil
}
//注册更新服务
func registerApiUpdate() {

	services := []*gdata.ServerInfo{
		&gdata.ServerInfo{
			APIName:	apilist.SimpleGatewayExampleGatewayUpdateApi,		//cometstone.simpleGatewayExample.gateway.updatApi
			IP:			Conf.Common.RealIp+":"+Conf.Api.ApiUpdatePort,		//10.50.36.184:1321
			Path:		"/api/update",
			Load:		1,
		},
	}

	//上报到etcd
	errCh := sdk.StoreServersByApi(etcdCli,services)

	//监听errch通道
	for {
		select {
		case err := <-errCh:
			Logger.Warn("请求ETCD异常",zap.Error(err),zap.Any("etcd_addr",Conf.Etcd.Addrs))
		}
	}



}