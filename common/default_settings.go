package common

//tb_domain_entry init
//entry_type = group
var SENSITIVE_GROUP = []string{
	"Administrators",
	"Power Users",
	"Account Operators",
	"Server Operators",
	"Print Operators",
	"Backup Operators",
	"Replicators",
	"Network Configuration Operators",
	"Incoming Forest Trust Builders",
	"Domain Admins",
	"Domain Controllers",
	"Group Policy Creator Owners",
	"read-only Domain Controllers",
	"Enterprise Read-only Domain Controllers",
	"Schema Admins",
	"Enterprise Admins",
	"Microsoft Exchange Servers",
	"Remote Desktop Users",
	"DnsAdmins",
}

// KerberosList kerberos的种类
var KerberosList = []string{"ST_max_cycle", "learn_day", "TGT_max_cycle", "high_risk_spn", "high_risk_delegation"}

// KerberosConfMap kerberos默认值
var KerberosConfMap = map[string][]string{
	"TGT_max_cycle":        {"10"},
	"ST_max_cycle":         {"600"},
	"learn_day":            {"10"},
	"high_risk_spn":        {"MSSQLSvc", "MSSQL", "FIMService", "AGPMServer", "exchangeMDB", "TERMSERV", "WSMAN", "Microsoft Virtual Console Service", "STS"},
	"high_risk_delegation": {"ldap/", "http/", "HOST/", "cifs/", "krbtgt/", "mssqlsvc/"},
}

// 前端类型展示map
var KerberosShouTypeMap = map[string]string{
	"TGT_max_cycle":        "input",
	"ST_max_cycle":         "input",
	"learn_day":            "input",
	"high_risk_spn":        "tag",
	"high_risk_delegation": "tag",
}

// 前端字段后缀展示map
var KerberosFileSuffixMap = map[string]string{
	"TGT_max_cycle":        "小时",
	"ST_max_cycle":         "分钟",
	"learn_day":            "天",
	"high_risk_spn":        "",
	"high_risk_delegation": "",
}

var KerberosRedisMap = map[string]string{
	"TGT_max_cycle":        "max_tgt_active_minute",
	"ST_max_cycle":         "max_st_active_minute",
	"learn_day":            "engine_learn_day",
	"high_risk_spn":        "",
	"high_risk_delegation": "",
}

// KerberosNameMap kerberos中文释义
var KerberosNameMap = map[string]string{
	"TGT_max_cycle":        "TGT最大化生命周期",
	"ST_max_cycle":         "ST最大化生命周期",
	"learn_day":            "学习周期",
	"high_risk_spn":        "高风险SPN前缀",
	"high_risk_delegation": "高风险委托前缀",
}

// KerberosDescMap kerberos中文释义
var KerberosDescMap = map[string]string{
	"TGT_max_cycle":        "AD域内的Kerberos认证TGT票据的最大周期，此配置影响安全检测规则准确度。",
	"ST_max_cycle":         "AD域内的Kerberos认证ST票据的最大周期，此配置影响安全检测规则准确度。",
	"learn_day":            "在该周期内域安全管家将收集当前域内的特征数据",
	"high_risk_spn":        "AD域内用户的SPN被请求时可能伴随恶意攻击，此配置为常见的高危SPN，此配置影响安全检测规则准确度。",
	"high_risk_delegation": "AD域内可能会造成安全风险的委派前缀，此配置影响安全检测规则准确度。",
}

// KerberosRuleTypeMap kerberos中文释义
var KerberosRuleTypeMap = map[string]string{
	"TGT_max_cycle":        "kerberos",
	"ST_max_cycle":         "kerberos",
	"learn_day":            "UEBA",
	"high_risk_spn":        "kerberos",
	"high_risk_delegation": "kerberos",
}
