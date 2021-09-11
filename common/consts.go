// author: s0nnet
// time: 2020-09-01
// desc:

package common

// 常量定义
const (
	CONF_PATH   = "./apiserver.yaml" // 线上环境通过API_SRV_CONF_PATH指定
	JWT_SECRET  = "DFCaAXUdKm3scpQW"
)

//return status
const (
	RESP_SUCCESS = "success"
	RESP_FAILED  = "failed"
)

// 用户角色定义
const (
	RoleDev = "dev"
	RoleOps = "ops"
	RoleSec = "sec"
	RoleMgr = "mgr"
)

// 用户权限等级定义
const (
	PrivSuper = 1
)

const RDX_CRYPT_SECRET = "1d34c6b89a3bc27f" // redis数据加密密钥

//login config
const (
	LoginErrorCount      = 5   //登录错误次数限制
	LoginErrorExpiration = 300 //登录错误过期时间
)
