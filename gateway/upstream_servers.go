package gateway

import (
	"context"
	"github.com/coreos/etcd/clientv3"
	"fmt"
	"bytes"
	"simpleGatewayExample/global/apilist"
	"go.uber.org/zap"
	"sort"
)

//定时watch 服务状态
func watchUpstramServers() {
	for true {
		// 监听 /api/servers/ 下的所有服务
		rch := etcdCli.Watch(context.Background(),Conf.Etcd.ServerKey,clientv3.WithPrefix())
		for resp := range rch {
			for _, ev := range resp.Events {
				//mvccpb
				//const (
				//PUT    Event_EventType = 0
				//DELETE Event_EventType = 1
				//)

				//put
				if ev.Type == 0 {
					//解析load 跟ip
					ip,load,ipIndex := ipAndLoad(ev.Kv.Key,ev.Kv.Value)
					fmt.Printf("--watchUpstramServers-- ip:%s	load:%d		ipIndex:%d",ip,load,ipIndex)

					rest := ev.Kv.Key[:ipIndex]
					fmt.Printf("--watchUpstramServers-- rest:%s",rest)
					// 解析出apiName
					apiIndex :=bytes.LastIndex(rest,[]byte{'/'})
					apiName := string(rest[apiIndex+1:])

					if apiName == apilist.SimpleGatewayExampleGatewayUpdateApi {
						continue
					}
					//获取服务集合
					apiI, ok := apis.Load(apiName)
					if !ok {
						// api不存在 返回错误
						Logger.Info("api不存在，但是取到了服务器的打点信息", zap.String("api", apiName))
						continue
					}
					api := apiI.(*Api)
					// 更新对应的ip
					ipExist := false

					for _, server := range api.UpstreamServers {
						if server.IP == ip {
							//更新load 根据延迟负载
							server.Load = load
							ipExist = true
						}
					}
					//如果没找到新增
					if !ipExist {
						newServer := &UpstreamServer{
							IP:ip,
							Load:load,
						}

						api.UpstreamServers = append(api.UpstreamServers,newServer)
					}

					sort.Slice(api.UpstreamServers, func(i, j int) bool {
						return api.UpstreamServers[i].Load < api.UpstreamServers[j].Load
					})

					fmt.Printf("watch, key: %s, ip : %s, load: %v\n", apiName, ip, load)
					for _, s := range api.UpstreamServers {
						fmt.Println("etcd watch插入,该api的最新服务器列表: ", *s)
					}

				}else if ev.Type == 1 {
					//delete
					// 解析出ip
					ip, _, ipIndex := ipAndLoad(ev.Kv.Key, ev.Kv.Value)

					rest := ev.Kv.Key[:ipIndex]

					// 解析出apiName
					apiIndex := bytes.LastIndex(rest, []byte{'/'})
					apiName := string(rest[apiIndex+1:])

					apiI, ok := apis.Load(apiName)

					if !ok {
						// api不存在，返回错误
						Logger.Info("api不存在，但是取到了服务器的打点信息", zap.String("api", apiName))
					}

					api := apiI.(*Api)
					for i,server := range api.UpstreamServers {

						if i == len(api.UpstreamServers) {
							break
						}

						if server.IP == ip {
							//删除服务 ip
							api.UpstreamServers = append(api.UpstreamServers[:i],api.UpstreamServers[i+1:]...)
						}
					}

					for _, s := range api.UpstreamServers {
						fmt.Println("etcd watch 删除,该api的最新服务器列表: ", *s)
					}
				}

			}
		}


	}
}
