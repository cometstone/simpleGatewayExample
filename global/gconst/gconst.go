package gconst


// 存储api和节点信息的过期时间
const APIStoreInterval = 60
const APILeaseTime = APIStoreInterval + 10

// ApisRootPath 是通过Api查找服务器地址时，etcd中的根目录
const APIsRootPath = "/api/servers/"
const InApiRootPath = "/in/api/servers/"

