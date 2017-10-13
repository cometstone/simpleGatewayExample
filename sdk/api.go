package sdk

import (
	"github.com/coreos/etcd/clientv3"
	"simpleGatewayExample/global/gdata"
	"time"
	"simpleGatewayExample/global/gconst"
	"strconv"
	"context"
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
			//每分钟上报一次服务信息
			time.Sleep(time.Second * gconst.APIStoreInterval)
		}

	}()

	return errCh
}