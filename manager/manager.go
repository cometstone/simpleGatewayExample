package manager

import (
	"go.uber.org/zap"
	"simpleGatewayExample/common"
	"database/sql"
	"simpleGatewayExample/sdk"
	"simpleGatewayExample/global/servicelist"
)

var Logger *zap.Logger
var Conf *common.Config

var db *sql.DB

func Start() {
	common.InitConfig()
	Conf = common.Conf

	sdk.InitLogger(Conf.Common.LogLevel,Conf.Common.LogPath,Conf.Common.IsDebug,servicelist.SimpleGatewayExampleManager)
	db = sdk.InitMysql(&sdk.MysqlConf{
		Acc:		Conf.Mysql.Acc,
		Pw:			Conf.Mysql.Pw,
		Addr:		Conf.Mysql.Addr,
		Port:		Conf.Mysql.Port,
		Database:	Conf.Mysql.Database,
	})

	//初始化 获取在活的服务器列表
	initGatewayUpdate()

}
