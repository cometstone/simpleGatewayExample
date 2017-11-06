package gateway

import (
	"github.com/labstack/echo"
)
//统一入口 跟route分发
func apiRoute(c echo.Context) error {
	// 生成request_id
	// 获取debug选项
	// 查询是否存在此Api
	// 记录请求IP、参数
	// api不存在，返回错误
	// 记录请求的API信息
	// 生成url
	// 使用fasthttp    或者(http) (grpc) 用于通信
	// 设置Method
	// 透传参数
	// 透传cookie
	// 设置X-FORWARD-FOR 		记录真实IP
	// 请求
	// 拼接url
	return  nil
}
