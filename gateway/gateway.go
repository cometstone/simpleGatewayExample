package gateway

import (
	"go.uber.org/zap"
	"database/sql"
	"simpleGatewayExample/common"
	"simpleGatewayExample/sdk"
	"simpleGatewayExample/global/servicelist"
)


var Logger *zap.Logger
var Conf *common.Config
var db *sql.DB

func Start() {
	//加载文件
	common.InitConfig()
	Conf = common.Conf
	//根据配置初始化相关服务
	//从配置初始化LOG
	sdk.InitLogger(Conf.Common.LogLevel,Conf.Common.LogPath,Conf.Common.IsDebug,servicelist.SimpleGatewayExampleGateway)
	Logger = sdk.Logger

	//链接数据库
	db = sdk.InitMysql(&sdk.MysqlConf{
		Acc:		Conf.Mysql.Acc,
		Pw:			Conf.Mysql.Pw,
		Addr:		Conf.Mysql.Addr,
		Port:		Conf.Mysql.Port,
		Database:	Conf.Mysql.Database,
	})
	//初始化etcd
	initEtcd()
	//API更新服务初始化
	initUpdateApi()




}