package gdata

/* 从Mysql中加载API信息到内存中*/

type API struct {
	ID 			int 		`json:"id"`
	FullName 	string		`json:"api_name"`
	Company		string		`json:"company"`
	Product		string		`json:"product"`
	System    string `json:"system"`
	Interface string `json:"interface"`
	Version   string `json:"version"`

	Method string `json:"method"`

	ProxyMode string `json:"proxy_mode"`

	UpstreamMode  string `json:"upstream_mode"`
	UpstreamValue string `json:"upstream_value"`
}


// 通过Api查找服务器地址时，返回的Server信息
type ServerInfo struct {
	APIName		string
	IP 			string
	Path 		string
	Load 		float64
	Type 		int			//1.更新 2.删除
}