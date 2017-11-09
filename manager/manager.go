package manager

import (
	"go.uber.org/zap"
	"simpleGatewayExample/common"
	"database/sql"
	"simpleGatewayExample/sdk"
	"simpleGatewayExample/global/servicelist"
	"github.com/labstack/echo"
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

	e :=echo.New()
	e.POST("/api/create", apiCreate)
	e.POST("/api/update", apiUpdate)
	e.POST("/api/query", apiQuery)
	e.POST("/api/delete", apiDelete)
	e.GET("/api/list", apiList)

	e.Logger.Fatal(e.Start(":" + Conf.Admin.ManagerPort))

}
