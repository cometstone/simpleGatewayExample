package gateway

import (
	"sync"
	"fmt"
	"go.uber.org/zap"
	"database/sql"
	"simpleGatewayExample/global/gdata"
	"context"
	"simpleGatewayExample/global/gconst"
	"github.com/coreos/etcd/clientv3"
	"strings"
	"strconv"
	"bytes"
	"sort"
)

type Api struct {
	// domain.group.service.version
	// 只能由字母、数字、点组成
	FullName string
	//GET POST
	Method string
	// 1.Raw: 将请求的Path直接append在upstream_value(url)后
	// 2.Indirect: 直接访问upstream_value(url)
	ProxyMode string

	// 1.直接寻址： url = upstream_value
	// 2.间接寻址: 在etcd中取出key为Api.Name的值，返回的数据结构存储在UpstreamValue
	UpstreamMode    string
	UpstreamServers []*UpstreamServer
}
//IP跟权重
type UpstreamServer struct {
	Load 	float64
	IP		string
}

type Apis struct {
	*sync.Map
}

var apis *Apis

func (a *Apis) LoadAll() {
	query := fmt.Sprintf("select * from api")
	rows,err := db.Query(query)
	if err != nil {
		Logger.Fatal("query gateway.api error ",zap.Error(err))
	}

	for rows.Next() {
		load(rows)
	}
}

func load(rows *sql.Rows) {
	rowApi := &gdata.API{}
	err := rows.Scan(&rowApi.ID,&rowApi.FullName,&rowApi.Company,&rowApi.Product,&rowApi.System,&rowApi.Interface, &rowApi.Version,
		&rowApi.Method, &rowApi.ProxyMode, &rowApi.UpstreamMode, &rowApi.UpstreamValue)

	if err != nil {
		Logger.Fatal("scan simpleGatewayExample.api error ",zap.Error(err))
	}

	//添加到服务组里

	api := &Api{}
	api.FullName = rowApi.FullName
	api.Method = rowApi.Method
	api.ProxyMode = rowApi.ProxyMode
	api.UpstreamMode = rowApi.UpstreamMode

	if api.UpstreamMode == "1" {
		//直接寻址 放入servers
		api.UpstreamServers = []*UpstreamServer{
			&UpstreamServer{
				Load:1,
				IP:rowApi.UpstreamValue,
			},
		}
	}else {
		//在etcd中取出key 前缀为gconst.APIsRootPath+api.FullName 的服务
		resp,err:=etcdCli.Get(context.Background(),gconst.APIsRootPath+api.FullName,clientv3.WithPrefix())
		if err != nil {
			Logger.Fatal("etcd get err",zap.Error(err))
		}
		//从etcd 读出服务报春到apis里
		servers := make([]*UpstreamServer,len(resp.Kvs))

		for _, v := range resp.Kvs {
			//从etcd 的k v 中解析出
			ip,load,_ := ipAndLoad(v.Key,v.Value)
			servers = append(servers,&UpstreamServer{
				IP:ip,
				Load:load,
			})
		}
		//负载排序
		sort.Slice(servers, func(i, j int) bool {
			return servers[i].Load < servers[j].Load
		})
		api.UpstreamServers = servers
		for _,s := range api.UpstreamServers {
			fmt.Printf("api load: %s 的最新服务器列表: %v\n", api.FullName, *s)
		}
	}

	fmt.Println(*api)
	//将服务组放入 apis里
	apis.Store(api.FullName,api)

}

func ipAndLoad(key []byte, val []byte) (string, float64, int) {
	// 解析load
	loadPath := strings.Split(string(val), "--")

	load, _ := strconv.ParseFloat(loadPath[0], 64)
	path := loadPath[1]

	// 解析出ip
	ipIndex := bytes.LastIndex(key, []byte{'/'})
	ip := "http://" + string(key[ipIndex+1:]) + path

	return ip, load, ipIndex
}
