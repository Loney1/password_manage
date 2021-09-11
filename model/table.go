// author: s0nnet
// time: 2020-09-01
// desc:

package model

import (
	"time"
)

//用户表
type User struct {
	ID           int32     `gorm:"primary_key"`          // ID
	UserName     string    `gorm:"column:username"`      // username
	Password     string    `gorm:"column:password"`      // passwd
	Pri          int32     `gorm:"column:pri"`           // 权限类别： 1:super 2:admin
	Role         string    `gorm:"column:role"`          // 用户角色
	PassStrength string    `gorm:"column:pass_strength"` // 密码强度: high/middle/low
	Mobile       string    `gorm:"column:mobile"`        // 手机号
	Email        string    `gorm:"column:email"`         // 邮箱
	Remark       string    `gorm:"column:remark"`        // 备注
	Token        string    `gorm:"column:token"`         // token
	RealName     string    `gorm:"column:real_name"`     // 真实姓名
	Department   string    `gorm:"column:department"`    // 部门
	Post         string    `gorm:"column:post"`          // 岗位
	Address      string    `gorm:"column:address"`       // 所在地址
	Secret       string    `gorm:"column:secret"`        // 密钥
	MfaStatus    string    `gorm:"column:mfa_status"`    // 标志状态 enable、disable stop
	CreatedAt    time.Time `gorm:"column:created_at"`    // 添加时间
	UpdatedAt    time.Time `gorm:"column:updated_at"`    // 密码更新时间 add by xc
}

func (a *User) TableName() string {
	return "user"
}

//机器表
type MachineUser struct {
	ID          int32     `gorm:"primary_key"`         //ID
	DN          string    `gorm:"column:dn"`           //DN
	Domain      string    `gorm:"column:domain"`       // 所在域
	MachineName string    `gorm:"column:machine_name"` //机器名
	MachinePwd  string    `gorm:"column:machine_pwd"`  //历史密码
	Remark      string    `gorm:"column:remark"`       //描述
	ExpiredAt   time.Time `gorm:"column:expired_at"`   //过期时间
	CreatedAt   time.Time `gorm:"column:created_tm"`   //生成时间
	UpdatedAt   time.Time `gorm:"column:update_tm"`    //更新时间
}

func (a *MachineUser) TableName() string {
	return "machine_user"
}

//域用户表
type Domain struct {
	ID         int32     `gorm:"primary_key"`        // ID
	Name       string    `gorm:"column:name"`        // 域名
	DCHostName string    `gorm:"column:dc_hostname"` // 域控DC主机名
	DNS        string    `gorm:"column:dns"`         // 域dns
	DN         string    `gorm:"column:dn"`          // 域的dn
	UserName   string    `gorm:"column:user_name"`   // 域用户信息
	Password   string    `gorm:"column:password"`    // 域用户密码
	UserDN     string    `gorm:"column:user_dn"`     // 域用户dn
	Status     int       `gorm:"column:status"`      // run|stop|init|error
	ErrMsg     string    `gorm:"column:err_msg"`     // 错误信息
	CreatedAt  time.Time `gorm:"column:created_tm"`  // 添加时间
}

func (a *Domain) TableName() string {
	return "domain"
}

//日志审计表
type AuditLog struct {
	ID          int32     `gorm:"primary_key"`         // ID
	LoginUser   string    `gorm:"column:login_user"`   //登录用户
	SourceIp    string    `gorm:"column:source_ip"`    //源ip
	Event       string    `gorm:"column:event"`        //事件
	EventArgs   string    `gorm:"column:event_args"`   //事件参数
	EventResult string    `gorm:"column:event_result"` //事件结果 //成功 失败
	Status      int32     `gorm:"column:status"`       //数据状态 1删除 0正常
	UpdatedAt   time.Time `gorm:"column:update_tm"`    //更新时间
	CreatedAt   time.Time `gorm:"column:create_tm"`    //添加时间
}

func (a *AuditLog) TableName() string {
	return "audit_log"
}
