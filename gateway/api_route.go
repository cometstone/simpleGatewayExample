package gateway

import (
	"github.com/labstack/echo"
	"time"
	"simpleGatewayExample/common"
	"go.uber.org/zap"
	"net/http"
	"simpleGatewayExample/sdk"
	"github.com/valyala/fasthttp"
	"fmt"
)

var cli = &fasthttp.Client{}

//统一入口 跟route分发
func apiRoute(c echo.Context) error {
	ts := time.Now()
	// 生成request_id
	rid :=common.RequestID()

	// 获取debug选项
	if c.FormValue("log_debug") == "on" {
		c.Set("debug_on",true)
	}else {
		c.Set("debug_on",false)
	}
	// 查询是否存在此Api
	apiName := c.FormValue("api_name")
	// 记录请求IP、参数
	apiI,ok := apis.Load(apiName)
	Logger.Info("收到新请求", zap.String("rid", rid), zap.String("ip", c.RealIP()), zap.String("api_name", apiName), zap.Bool("api_exist", ok))

	if !ok {
		// api不存在，返回错误
		return c.String(http.StatusOK, "Api不存在")
	}

	// 记录请求的API信息
	api := apiI.(*Api)
	sdk.DebugLog(rid,c.Get("debug_on").(bool),"请求的api", zap.Any("api", *api))

	// 生成url
	var upstreamUrl string
	if len(api.UpstreamServers) <= 0 {
		Logger.Warn("要请求的服务不存活",zap.String("rid",rid),zap.String("api_name",apiName))
		return c.String(http.StatusOK,"该Api没有服务器存活")
	}

	upstreamUrl = api.UpstreamServers[0].IP
	sdk.DebugLog(rid,c.Get("debug_on").(bool),"raw upstream url",zap.String("url",c.Request().RequestURI),zap.String("raw_upstream_url",upstreamUrl))

	// 使用fasthttp    或者(http) (grpc) 用于通信
	req := &fasthttp.Request{}
	resp := &fasthttp.Response{}

	// 设置Method
	req.Header.SetMethod(api.Method)

	// 透传参数
	args := &fasthttp.Args{}
	params,err := c.FormParams()
	if err != nil {
		Logger.Fatal("解析请求参数错误",zap.String("rid",rid),zap.Error(err),zap.String("parmas",params.Encode()))
		return c.String(http.StatusOK, "获取参数错误")
	}
	for parma,_:=range params {
		args.Set(parma,c.FormValue(parma))
	}
	args.Set("rid",rid)
	sdk.DebugLog(rid, c.Get("debug_on").(bool), "请求参数", zap.String("paramas", params.Encode()))
	// 透传cookie
	for _, cookie := range c.Cookies() {
		req.Header.SetCookie(cookie.Name,cookie.Value)
	}
	sdk.DebugLog(rid,c.Get("debug_on").(bool),"请求cookie",zap.Any("cookies",c.Cookies()))

	// 设置X-FORWARD-FOR 		记录真实IP
	req.Header.Set("X-Forwarded-For",c.RealIP())
	// 请求
	switch api.Method {
	case "POST":
		args.WriteTo(req.BodyWriter())
	case "GET":
		// 拼接url
		upstreamUrl = upstreamUrl + "?" + args.String()
	}

	sdk.DebugLog(rid,c.Get("debug_on").(bool),"最终 upstreamUrl",zap.String("upstreamUrl",upstreamUrl))

	req.SetRequestURI(upstreamUrl)
	err = cli.DoTimeout(req,resp,10*time.Second)
	if err != nil {
		Logger.Info("api请求错误",zap.String("rid",rid),zap.Error(err),zap.String("apiName",apiName))
		return c.String(resp.StatusCode(),err.Error())
	}

	if resp.StatusCode() != 200 {
		Logger.Info("api请求code不为200",zap.String("rid",rid),zap.Int("code",resp.StatusCode()),zap.String("apiName",apiName))
		return c.String(resp.StatusCode(),fmt.Sprintf("请求返回Code异常：%v",resp.StatusCode()))
	}

	Logger.Info("api请求成功",zap.String("rid",rid),zap.Int64("TimeDifference",time.Now().Sub(ts).Nanoseconds()/1000))
	sdk.DebugLog(rid, c.Get("debug_on").(bool), "api请求返回body", zap.String("body", string(resp.Body())))
	return c.String(http.StatusOK, string(resp.Body()))

}
