package gdata

/* 从Mysql中加载API信息到内存中*/


// 通过Api查找服务器地址时，返回的Server信息

type ServerInfo struct {
	APIName		string
	IP 			string
	Path 		string
	Load 		float64
	Type 		int			//1.更新 2.删除
}