package gateway

import (
	"go.uber.org/zap"
	"database/sql"
	"simpleGatewayExample/common"
)


var Logger *zap.Logger
var Conf *common.Config
var db *sql.DB

func Start() {
	//加载文件
	common.InitConfig()
	Conf = common.Conf
	//根据配置初始化相关服务

}