package manager

import (
	"github.com/labstack/echo"
	"simpleGatewayExample/common"
	"simpleGatewayExample/sdk"
	"simpleGatewayExample/global/gdata"
	"net/http"
	"strings"
	"fmt"
	"go.uber.org/zap"
	"time"
)


type ExtApiRes struct {
	Suc  bool                   `json:"suc"`
	Data map[string]interface{} `json:"data"`
}


func apiCreate(c echo.Context) error {
	//设置跨域
	c.Response().Header().Set("Access-Control-Allow-Origin","*")
	//链路id号
	rid := common.RequestID()
	//时间 跟设置
	ts,debugOn := sdk.LogExtra(c)
	//拼接参数
	api := &gdata.API{}
	api.FullName = c.FormValue("api_name")
	if api.FullName == "" {
		return c.JSON(http.StatusOK,&ExtApiRes{
			Suc:false,
			Data:map[string]interface{}{
				"msg":"必填项不能为空",
			},
		})
	}

	api.Method = c.FormValue("method")
	api.UpstreamMode = c.FormValue("upstream_mode")
	api.UpstreamValue = c.FormValue("upstream_value")
	api.ProxyMode = "1"

	//检查 FullName 是否符合规范
	names := strings.Split(api.FullName,".")
	if len(names) != 5 {
		return c.JSON(http.StatusOK,&ExtApiRes{
			Suc:false,
			Data:map[string]interface{}{
				"msg":"Api名必须为(公司.产品.系统.接口.版本号)格式,例如cometstone.simpGatewayExample.manager.apiCreate.v1",
			},
		})
	}

	api.Company, api.Product, api.System, api.Interface, api.Version = names[0], names[1], names[2], names[3], names[4]
	query := fmt.Sprintf("INSERT INTO api (`full_name`,`company`,`product`,`system`,`interface`,`version`,`method`,`proxy_mode`,`upstream_mode`,`upstream_value`) VALUES ('%s', '%s','%s','%s','%s','%s','%s','%s','%s','%s')", api.FullName, api.Company, api.Product, api.System, api.Interface, api.Version, api.Method, api.ProxyMode,
		api.UpstreamMode, api.UpstreamValue)

	//用于调试查看sql
	sdk.DebugLog(rid,debugOn,"创建api sql",zap.String("sql",query))

	_, err := db.Exec(query)

	if err != nil {
		Logger.Info("api create, insert error", zap.String("rid", rid), zap.Error(err), zap.String("query", query))
		return c.String(http.StatusOK, "create api error")
	}

	Logger.Info("api创建成功", zap.String("rid", rid), zap.Int64("TimeDifference", time.Now().Sub(ts).Nanoseconds()/1000))

	updateApi(api.FullName, 1)

	return c.JSON(http.StatusOK, &ExtApiRes{
		Suc: true,
		Data: map[string]interface{}{
			"api": api,
		},
	})

}

func apiUpdate(c echo.Context) error {
	c.Response().Header().Set("Access-Control-allow-Origin","*")
	rid := common.RequestID()
	ts,debugOn := sdk.LogExtra(c)

	api := &gdata.API{}
	api.FullName = c.FormValue("api_name")
	api.Method = c.FormValue("method")
	api.ProxyMode = c.FormValue("proxy_mode")
	api.UpstreamMode = c.FormValue("upstream_mode")
	api.UpstreamValue = c.FormValue("upstream_value")

	if api.FullName == "" {
		Logger.Info("api_name不能为空", zap.String("rid", rid))
		return c.JSON(http.StatusOK, &ExtApiRes{
			Suc: false,
			Data: map[string]interface{}{
				"msg": "api名不能为空",
			},
		})
	}

	Logger.Info("api 更新",zap.String("rid",rid),zap.Any("api",*api))

	query := fmt.Sprintf("UPDATE api SET `method`='%s',`proxy_mode`='%s',`upstream_mode`='%s',`upstream_value`='%s' WHERE `full_name`='%s'",
		api.Method, api.ProxyMode, api.UpstreamMode, api.UpstreamValue, api.FullName)
	sdk.DebugLog(rid, debugOn, "更新api sql", zap.String("sql", query))

	_,err := db.Exec(query)
	if err != nil {
		Logger.Info("api update  error", zap.String("rid", rid), zap.Error(err), zap.String("query", query))
		return c.String(http.StatusOK, "update api error")
	}

	Logger.Info("api更新成功", zap.String("rid", rid), zap.Int64("TimeDifference", time.Now().Sub(ts).Nanoseconds()/1000))

	updateApi(api.FullName,2)

	return c.JSON(http.StatusOK,&ExtApiRes{
		Suc:true,
	})

}

func apiQuery(c echo.Context) error {

	return nil
}

func apiDelete(c echo.Context) error {

	return nil
}

func apiList(c echo.Context) error {

	return nil
}
