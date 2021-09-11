// author: s0nnet
// time: 2020-09-12
// desc: 用户访问权限控制列表 ACL

package v1

import (
	"strings"

	"adp_backend/common"

	logger "github.com/sirupsen/logrus"
)

var serviceName = _ADM_serviceDesc.ServiceName

// URL事件映射定义
var URLEventMap = map[string]string{
	// 登录页
	"/" + serviceName + "/" + "Login":  "登录",
	"/" + serviceName + "/" + "Logout": "退出登录",

	// 检测中心
	// 告警事件
	"/" + serviceName + "/" + "ExportThreatEvent": "导出告警事件",
	"/" + serviceName + "/" + "UpdateThreatEvent": "更新告警事件",
	// 监测配置
	"/" + serviceName + "/" + "AddDomainEntry":     "新增域条目配置",
	"/" + serviceName + "/" + "DeleteDomainEntry":  "修改域条目配置",
	"/" + serviceName + "/" + "UpdateKerberosConf": "修改Kerberos配置",
	// 白名单管理
	"/" + serviceName + "/" + "AddRuleWhitelist":    "添加白名单",
	"/" + serviceName + "/" + "DeleteRuleWhitelist": "删除白名单",
	"/" + serviceName + "/" + "UpdateRuleWhitelist": "编辑白名单",
	// 联动配置
	"/" + serviceName + "/" + "SetAlertConf":  "配置告警邮件",
	"/" + serviceName + "/" + "TestEmailSend": "测试邮件",

	// 主动监测
	"/" + serviceName + "/" + "ScanInspection":           "一键巡检",
	"/" + serviceName + "/" + "SetCronTask":              "设置定期扫描任务",
	"/" + serviceName + "/" + "ExportScanEvent":          "导出主动检测事件列表",
	"/" + serviceName + "/" + "ScanLeakEvent":            "漏洞检测",
	"/" + serviceName + "/" + "UpdateScanPluginEnable":   "更新漏洞扫描器配置",
	"/" + serviceName + "/" + "UpdateScanPluginMetaData": "设置plugin的MetaData",

	// 报表报告
	"/" + serviceName + "/" + "GenerateEventReport": "生成报表",
	"/" + serviceName + "/" + "DownloadEventReport": "导出报表",
	"/" + serviceName + "/" + "DeleteEventReport":   "删除报表",

	// 系统管理
	//域服务器配置
	"/" + serviceName + "/" + "AddDomain":    "添加域控",
	"/" + serviceName + "/" + "UpdateDomain": "编辑域控信息",
	"/" + serviceName + "/" + "TestDomain":   "测试域连接",
	"/" + serviceName + "/" + "DeleteDomain": "删除域控",
	// 传感器管理
	"/" + serviceName + "/" + "DownloadAgent":      "下载域控传感器",
	"/" + serviceName + "/" + "DownHttps":          "下载证书",
	"/" + serviceName + "/" + "UpdateAgent":        "编辑域控传感器",
	"/" + serviceName + "/" + "DeleteAgent":        "删除域控传感器",
	"/" + serviceName + "/" + "UpdateAgentVersion": "更新域控传感器",
	"/" + serviceName + "/" + "DeleteWecBeat":      "删除日志",
	"/" + serviceName + "/" + "SetMsRCP":           "设置MSRPC日志采集",
	// 个人中心
	"/" + serviceName + "/" + "UpdateUserPassword": "修改密码",
	"/" + serviceName + "/" + "UpdateUser":         "修改用户信息",
	"/" + serviceName + "/" + "UpdateAvatar":       "上传头像",
	"/" + serviceName + "/" + "EnableMfa":          "开启登录二次校验",
	"/" + serviceName + "/" + "DisableMfa":         "关闭登录二次校验",
	// 子账户管理
	"/" + serviceName + "/" + "AddUser":       "添加子用户",
	"/" + serviceName + "/" + "ResetPassword": "重置密码",
	"/" + serviceName + "/" + "DeleteUser":    "删除子账户",
	// 系统信息
	"/" + serviceName + "/" + "UpdateReboot":      "更新重启",
	"/" + serviceName + "/" + "UpdateSystemIcon":  "上传Logo",
	"/" + serviceName + "/" + "UpdateLicence":     "更新授权",
	"/" + serviceName + "/" + "DownloadSystemLog": "下载系统日志",
	"/" + serviceName + "/" + "SetSystemTime":     "更新系统时间",

	// 日志审计
	"/" + serviceName + "/" + "ExportAuditLog": "导出审计日志",
	"/" + serviceName + "/" + "DeleteAuditLog": "清空审计日志",
	// 网络检测
	"/" + serviceName + "/" + "NetworkDiag": "网络诊断",
	// 通知模块
	"/" + serviceName + "/" + "UpdateNotifyConf":       "更新通知信息",
	"/" + serviceName + "/" + "UpdateNotifyConfEnable": "修改通知启动状态",
}

//URL事件脱敏参数定义
var URLEventMaskingMap = map[string][]string{
	"/" + serviceName + "/" + "Login":               []string{"password"},
	"/" + serviceName + "/" + "AddUser":             []string{"password"},
	"/" + serviceName + "/" + "UpdateUser":          []string{"password"},
	"/" + serviceName + "/" + "UpdateUserPassword":  []string{"oldPassword", "newPassword"},
	"/" + serviceName + "/" + "SetAlertConf":        []string{"config"},
	"/" + serviceName + "/" + "AddDomain":           []string{"password"},
	"/" + serviceName + "/" + "UpdateDomain":        []string{"password"},
	"/" + serviceName + "/" + "TestDomain":          []string{"password"},
	"/" + serviceName + "/" + "ResetPassword":       []string{"newPassword"},
	"/" + serviceName + "/" + "TestEmailSend":       []string{"Config"},
	"/" + serviceName + "/" + "UpdateNotifyConfReq": []string{"senderIdentity"},
}

var moduleMap = map[string][]string{
	// 风险大盘
	"RiskMarket": []string{"StatsAlertActivity", "StatsRiskAssets", "StatsAlertEvents", "StatsScanEvents", "StatsAssets", "AlarmAnalysis", "RiskTrend", "ListStatsAlertName", "ListStatsAlertType"},
	//告警列表
	"ThreatEventFind": []string{"ListThreatEvent", "ListThreatActivity", "ListThreatRawLog", "GetRuleInfo", "GetDCNameList", "GetTarget", "GetDomainFromAlert", "ListThreatEventSearch", "ListRuleTypes", "StateAlertEventByRule", "GetThreatEventByUniqueID"},
	// 告警列表操作
	"ThreatEventOperating": []string{"UpdateThreatEvent", "ExportThreatEvent"},
	// 主动检测
	"Scanner": []string{"GetScanRule", "ScanInspection", "GetScanTaskState", "GetScanScore", "SetCronTask", "ListCronTask", "EventList", "EventDetails",
		"LastScanInfo", "StopScan", "ListOnlineDomain", "ListDomainByScanEvent", "ExportScanEvent", "GetInstanceList", "ListTaskManagerGroup",
		"DetailTaskManagerGroup", "DeleteTaskManagerGroup", "ProtectInfo"},
	// 事件列表
	"ThreatList": []string{"ListRuleTypes", "ListDomainByThreat", "ListDCByThreat", "ListTypeByThreat", "ListLevelByThreat", "GetThreatList"},
	// 敏感组配置,蜜罐账户
	"SensitiveGroup": []string{"AddDomainEntry", "DeleteDomainEntry", "ListDomainEntry"},
	// 告警配置
	"Kerberos": []string{"GetKerberosConf", "UpdateKerberosConf", "ListKerberosConf"},
	// 白名单管理
	"RuleWhite": []string{"AddRuleWhitelist", "DeleteRuleWhitelist", "UpdateRuleWhitelist", "GetRuleWhitelist", "ListRuleWhitelist", "GetRuleWhitelistInfo", "ListWhiteField", "GetWhiteFieldValue"},
	// 联动配置
	"AlertConf": []string{"GetAlertConf", "SetAlertConf", "TestEmailSend"},
	// 通知模块
	"NotifyConf": []string{"ListNotifyConf", "UpdateNotifyConf", "UpdateNotifyConfEnable", "GetNotifyConfInfo", "ListNotifyTarget", "TestEmail", "SelectOptionNotify"},
	//域服务器配置
	"Domain": []string{"ListDomain", "AddDomain", "TestDomain", "UpdateDomain", "DeleteDomain", "GetDomainObjectInfo", "UpdateDomainData", "GetDomainInfo", "SetMsRCP"},
	// 运维管理员的域配置
	"OpsDomain": []string{"ListDomain", "AddDomain", "TestDomain", "UpdateDomain", "GetDomainObjectInfo", "UpdateDomainData", "GetDomainInfo", "SetMsRCP"},
	// 安全管理员的域配置
	"SecDomain": []string{"ListDomain", "GetDomainObjectInfo", "GetDomainInfo"},
	//传感器管理
	"Agent": []string{"UpdateAgent", "CmdAgent", "DownloadAgent", "DownCertificate",
		"DeleteAgent", "ListGateway", "ListWecBeat", "UpdateAgentVersion", "DeleteWecBeat", "GetDCList", "AddWecConf", "WecBeatInfo", "ListWecBeatEventInfo"},
	// 日志审计
	"AuditLog":       []string{"ListAuditLog", "ExportAuditLog"},
	"AuditLogDelete": []string{"DeleteAuditLog"},
	// 系统信息
	"System": []string{"GetSystemInfo", "DownloadSystemLog", "GetSystemLog", "UpdateReboot", "GetLicence", "UpdateLicence", "UpdateSystemIcon", "GetSystemIcon", "NetworkDiag", "SetSystemTime"},
	// 帮助中心
	"Help": []string{},
	// 个人中心
	"User": []string{"Login", "Logout", "ListUser", "AddUser", "UpdateUser", "UpdateUserPassword",
		"CheckMfa", "EnableMfa", "DisableMfa", "UpdateAvatar", "GetPwdUpdateTm"},
	"AccountManagement": []string{"DeleteUser", "ResetPassword"},
	//	消息模块
	"MessageNotify": []string{"ListNotify", "UpdateNotify", "AddNotifyEmailConf", "DeleteNotifyEmailConf", "UpdateNotifyEmailConf", "ListNotifyEmailConf", "StatsNotify"},
	// 事件报表
	"EventReport": []string{"GenerateEventReport", "ListEventReport", "StatusEventReport", "DownloadEventReport", "DeleteEventReport"},
	// 通用接口
	"All": {"ListAgent", "ListStatsAlertCount", "ListDomainNameForEventList", "ListDomainNameFromAgent", "ListDomain", "ListWhiteField", "ListDomainName", "GetDomainObject", "GetTaskState", "ListThreatEventSearch", "StateAlertEventByRule", "ListScanPluginType"},
	// 数据检索
	"Search": {"ListSearchLogEvent", "GetSearchLogField", "GetSearchChartData", "GetSearchFieldInfo"},
	// 资产相关接口
	"Assets": {"ListAssetsUser", "ListAssetsComputer", "ListAssetsGroup", "GetAssetsDetailsByAlert", "ListGroupByAssets",
		"GetAssetsActivities", "GetAssetsEntry", "GetAssetsLabel", "GetAssetsLabelInfo", "ListUsersSensitiveGroup", "StatsAssetsActivitiesLevel", "GetAssetsSensitiveGroupLabelInfo"},
	// 攻击路径
	"AttackPath": {"ListAttackPath", "ExportAttackPath"},
	// 漏洞检测
	"LeakEvent": {"ScanLeakEvent", "StatsLeakEvent", "ListLeakEvent", "ListScanPlugin", "UpdateScanPluginEnable", "UpdateScanPluginMetaData", "GetScanLeakEventStatus"},
}

func moduleMapJoin(strList ...string) string {
	str := ""
	for _, v := range strList {
		str += strings.Join(moduleMap[v], ",") + ","
	}
	return str
}

var UserACL = map[string]string{
	common.RoleMgr: moduleMapJoin("All", "AccountManagement", "Agent", "AlertConf", "AttackPath", "AuditLog", "AuditLogDelete", "Domain", "EventReport", "LeakEvent", "Help", "Honeypot", "Kerberos", "MessageNotify", "RiskMarket", "RuleWhite", "Scanner", "SensitiveGroup", "System", "ThreatEventFind", "ThreatEventOperating", "ThreatList", "User", "Search", "Assets", "SetSystemTime"),
	common.RoleSec: moduleMapJoin("All", "Agent", "AlertConf", "AttackPath", "AuditLog", "AuditLogDelete", "EventReport", "Help", "Honeypot", "Kerberos", "LeakEvent", "MessageNotify", "RiskMarket", "RuleWhite", "Scanner", "SensitiveGroup", "SecDomain", "System", "ThreatEventFind", "ThreatEventOperating", "ThreatList", "User", "Search", "Assets"),
	common.RoleOps: moduleMapJoin("All", "Agent", "AlertConf", "AttackPath", "AuditLog", "AuditLogDelete", "EventReport", "Help", "LeakEvent", "MessageNotify", "OpsDomain", "RiskMarket", "Scanner", "System", "ThreatEventFind", "ThreatList", "User", "Search", "Assets"),
}

func CheckUserAccess(role, fullMethod string) bool {
	paths := strings.SplitN(fullMethod, "/", 3)
	if len(paths) != 3 {
		return false
	}
	if paths[1] != serviceName {
		return false
	}

	acl, ok := UserACL[role]
	if !ok {
		logger.Warnf("invalid user role:%s, ignored.", role)
		return false
	}

	for _, method := range strings.Split(acl, ",") {
		if method == paths[2] {
			return true
		}
	}

	return false
}
