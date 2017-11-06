package gateway

import (
	"github.com/labstack/echo"
	"simpleGatewayExample/global/gdata"
	"simpleGatewayExample/global/apilist"
	"simpleGatewayExample/sdk"
	"go.uber.org/zap"
	"strings"
	"fmt"
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
	apiNames := strings.Split(c.FormValue("api_name"),",")
	tp :=c.FormValue("type")
	Logger.Info("api update",zap.Any("apiName",apiNames),zap.String("type",tp))
	//查找数据库

	for apiName := range apiNames{
		switch tp {
		//创建api,更新api
		case "1", "2":
			query := fmt.Sprintf("select * from api where full_name = %s",apiName)
			rows,err := db.Query(query);
			if err != nil {
				Logger.Warn("query simpleGatewayExample.api error:",zap.Error(err))
			}

			for rows.Next() {
				load(rows)
			}
		//删除
		case "3":
			apis.Delete(apiName)
		}
	}

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