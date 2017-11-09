package sdk

import (
	"github.com/coreos/etcd/clientv3"
	"simpleGatewayExample/global/gdata"
	"time"
	"simpleGatewayExample/global/gconst"
	"strconv"
	"context"
	"strings"
	"bytes"
	"fmt"
)
//  key='/APIsRootPath/ApiName/ip:port' val='load--path'    etcd put 格式
func StoreServersByApi(cli *clientv3.Client, servers []*gdata.ServerInfo) chan error{
	//上报错误通道
	errCh := make(chan error)

	go func() {
		for {
			for _,server:= range servers{
				//拼接参数 key val
				addr := server.IP
				key:= gconst.InApiRootPath+server.APIName+"/"+addr
				loadStr := strconv.FormatFloat(server.Load,'f',1,64)
				val := loadStr+"--"+server.Path

				//less组 节点的过期时间
				Grant,err := cli.Grant(context.TODO(),gconst.APILeaseTime)
				if err != nil {
					errCh <- err
					continue
				}
				//保存到etcd
				_,err = cli.Put(context.TODO(),key,val,clientv3.WithLease(Grant.ID))
				if err != nil {
					errCh <- err
					continue
				}
			}
			//每分钟更新一次务信息
			time.Sleep(time.Second * gconst.APIStoreInterval)
		}

	}()

	return errCh
}

// QueryServerByAPI 通过Api和Get请求查询在活的服务器列表
// n == 0 : 返回所有

type QueryServerRes struct {
	ApiName	string
	Servers []*gdata.ServerInfo
}
//每2分钟更新服务信息到 resch中
func QueryServerByAPI(cli *clientv3.Client, apiNames []string, n int) (chan *QueryServerRes, chan error) {

	resCh := make(chan *QueryServerRes)
	errCh := make(chan error)

	go func() {
		for {
			for _,apiName := range apiNames {
				res ,err := cli.Get(context.Background(),gconst.InApiRootPath+apiName,clientv3.WithPrefix())
				if err != nil {
					errCh <- err
					continue
				}

				var servers []*gdata.ServerInfo

				for i, kv := range res.Kvs {
					if n != 0 {
						if i+1 > n {
							break
						}
					}

					loadPath := strings.Split(string(kv.Value),"--")
					load,_ := strconv.ParseFloat(loadPath[0],10)

					// 解析出ip
					ipIndex := bytes.LastIndex(kv.Key,[]byte{'/'})
					ip := kv.Key[ipIndex+1:]
					servers = append(servers,&gdata.ServerInfo{
						APIName:apiName,
						IP:string(ip),
						Load:load,
						Path:loadPath[1],
					})
				}

				if len(servers) == 0 {
					errCh <- fmt.Errorf("%s没有服务器存活",apiName)
					continue
				}

				qsr := &QueryServerRes{
					ApiName:apiName,
					Servers:servers,
				}

				resCh <- qsr
			}

			time.Sleep(gconst.APIUpdateInterval*time.Second)

		}
	}()

	return resCh,errCh

}
