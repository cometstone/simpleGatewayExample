package common

import (
	"io/ioutil"
	"log"
	"gopkg.in/yaml.v2"
)

// 配置结构体
type Config struct {
	Common struct{
		Version  string
		IsDebug  bool `yaml:"debug"`
		LogPath  string
		LogLevel string
		Service  string
		RealIp   string
	}

	Api struct {
		GatewayPort   string
		ServerID      int
		ApiUpdatePort string
	}

	Admin struct {
		ManagerPort string
	}

	Mysql struct {
		Addr     string
		Port     string
		Database string
		Acc      string
		Pw       string
	}

	Etcd struct {
		Addrs     []string
		ServerKey string
	}
}

var Conf = &Config{}

func InitConfig() {
	data,err := ioutil.ReadFile("simpleGatewayExample.yaml")
	if err != nil {
		log.Fatal("read config err :",err)
	}
	err = yaml.Unmarshal(data,&Conf)
	if err != nil {
		log.Fatal("yaml decode err:",err)
	}

	log.Println(Conf)
}