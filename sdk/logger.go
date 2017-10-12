package sdk

import (
	"go.uber.org/zap"
	"fmt"
	"encoding/json"
	"go.uber.org/zap/zapcore"
	"log"
)

var Logger *zap.Logger
//https://github.com/uber-go/zap/blob/master/example_test.go
func InitLogger(lv string,lp string, isDebug bool, service string) {
	js := fmt.Sprintf(`{
		"level": "%s",
		"encoding": "json",
		"outputPaths": ["stdout","%s"],
		"errorOutputPaths": ["stderr","%s"]
	}`, lv, lp, lp)

	var cfg zap.Config
	if err := json.Unmarshal([]byte(js), &cfg); err != nil {
		panic(err)
	}

	cfg.EncoderConfig = zap.NewProductionEncoderConfig()
	cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	cfg.InitialFields = map[string]interface{}{
		"service": service,
	}

	var err error
	Logger, err = cfg.Build()
	if err != nil {
		log.Fatal("init logger error: ", err)
	}

	Logger.Info("logger初始化成功")

	// Output:
	// {"level":"info","message":"logger初始化成功","service":"xxx"}
}