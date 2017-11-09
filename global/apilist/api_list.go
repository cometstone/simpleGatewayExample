package apilist

import "simpleGatewayExample/global/servicelist"
/* Api列表在此处定义 */

// gateway监听此端口，等待api更新的数据


const SimpleGatewayExampleGatewayUpdateApi = servicelist.SimpleGatewayExampleGateway + "." + "updatApi"



// Account.staff
const StaffCheckLogin = "cometstone.account.staff.checkLogin"