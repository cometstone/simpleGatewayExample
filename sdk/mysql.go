package sdk

import (
	"database/sql"
	"fmt"
	"go.uber.org/zap"
	_ "github.com/go-sql-driver/mysql"
)

type MysqlConf struct {
	Acc      string
	Pw       string
	Addr     string
	Port     string
	Database string
}

func InitMysql(config *MysqlConf) *sql.DB {
	//格式化输出
	sqlcom := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s",config.Acc,config.Pw,config.Addr,config.Port,config.Database)
	//链接sql
	db,err := sql.Open("mysql",sqlcom)
	if err != nil {
		Logger.Fatal("init mysql error",zap.Error(err))
	}

	TestDb(db)

	return db
}

func TestDb(db *sql.DB) {
	err := db.Ping()
	if err != nil {
		Logger.Fatal("init mysql, ping error", zap.Error(err))
	}
}
