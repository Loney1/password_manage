package service

import (
	"context"
	"fmt"
	"strconv"
	"time"

	v1 "adp_backend/api/v1"
	"adp_backend/common"
	"adp_backend/config"
	utime "adp_backend/infra/time"
	"adp_backend/model"
	"adp_backend/server"
	"adp_backend/util"

	ldap3 "github.com/go-ldap/ldap/v3"
	logger "github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func SyncComputerList(env *config.Env) {
	for {
		time.Sleep(5 * time.Minute)
		//time.Sleep(5 * time.Second)

		// 查询所有的计算机列表
		domainList, err := server.FindDomainList(env.MysqlCli)
		if err != nil {
			continue
		}

		// 处理计算机列表数据
		var list []model.MachineUser
		for _, domain := range domainList {
			passwordDecode, err := util.PasswordDecode(domain.Password)
			if err != nil {
				logger.Warnf("password decode err:%v", err)
			}
			list, err = findComputerList(domain.DNS, domain.UserDN, domain.DN, passwordDecode, domain.Name)
			if err != nil {
				logger.Errorf("find computer by domain list err:%v", err)
			}
		}

		// 判断是应该查询还是更新
		for _, cUser := range list {
			err = server.UpsetMachineUserByName(env.MysqlCli, cUser)
			if err != nil {
				logger.Warnf("upset machine user by name err:%v", err)
				continue
			}
		}
	}
}

func (s *ADMServiceV1) ListMachineName(ctx context.Context, in *v1.ListMachineNameReq) (*v1.ListMachineNameReply, error) {
	list, count, err := server.FindMachineNameList(s.env.MysqlCli, in.SearchName, in.StartTime, in.EndTime, in.PageIdx, in.PageSize, in.SortTime)
	if err != nil {
		return nil, status.Error(codes.Internal, "获取机器列表失败")
	}

	var ret v1.ListMachineNameReply
	for _, user := range list {
		exTime := user.ExpiredAt.Format("2006-01-02 15:04:05")
		if time.Now().AddDate(-1, 0, 0).Sub(user.ExpiredAt) > 0 {
			exTime = "目标机器未安装LAPS服务"
		}
		ret.List = append(ret.List, &v1.ListMachineNameReply_Details{
			Id:       user.ID,
			HostName: user.MachineName,
			Dn:       user.DN,
			ExTime:   exTime,
			Remark:   user.Remark,
		})

	}

	ret.Page = &v1.ModelPage{
		PageIdx:  in.PageIdx,
		PageSize: in.PageSize,
		Total:    int32(count),
	}
	return &ret, nil
}

func (s *ADMServiceV1) GetMachinePwd(ctx context.Context, in *v1.GetMachinePwdReq) (*v1.GetMachinePwdReply, error) {
	mfaCode, err := strconv.Atoi(in.HasMfa)
	if err != nil {
		logger.Errorf("strconv atoi err:%v", err)
		return nil, status.Error(codes.Unauthenticated, "输入错误的验证码")
	}

	user, err := server.FirstUser(s.env)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "获取用户信息失败")
	}

	check := util.TotpCheck(user.Secret, mfaCode)
	if !check {
		return nil, status.Error(codes.Unauthenticated, "验证码错误")
	}

	machineUser, err := server.GetMachineUser(s.env.MysqlCli, in.Id)
	if err != nil {
		return nil, status.Error(codes.Internal, "获取机器用户信息失败")
	}

	domain, err := server.GetDomainByName(s.env.MysqlCli, machineUser.Domain)
	if err != nil {
		return nil, status.Error(codes.Internal, "获取域信息失败")
	}

	decodePWD, err := util.PasswordDecode(domain.Password)
	if err != nil {
		return nil, status.Error(codes.Internal, "域账户信息异常")
	}
	pwd, ex := getPWD(domain.DNS, domain.UserDN, domain.DN, decodePWD, machineUser.MachineName)
	fileTimeStamp := util.Atoll(ex)
	lastLogonTime := utime.FileTime2Time(fileTimeStamp)
	LastPwd := lastLogonTime.Format("2006-01-02 15:04:05")
	if LastPwd == "1601-01-01 00:00:00" {
		LastPwd = "目标机器未安装LAPS服务"
	}
	ret := v1.GetMachinePwdReply{
		MachinePwd: pwd,
		LastPwd:    LastPwd,
	}

	return &ret, nil
}

func (s *ADMServiceV1) UpdateMachineRemark(ctx context.Context, in *v1.UpdateMachineRemarkReq) (*v1.UpdateMachineRemarkReply, error) {

	err := server.UpdateMachineUserRemark(s.env.MysqlCli, in.Id, in.Remark)
	if err != nil {
		logger.Errorf("UpdateMachineUserRemark err:%v", err)
		return nil, status.Error(codes.Internal, "更新机器备注失败")
	}

	var ret v1.UpdateMachineRemarkReply
	ret.Result = common.RESP_SUCCESS

	return &ret, nil
}

func findComputerList(ip, userDN, dn, pwd, domainName string) ([]model.MachineUser, error) {
	dial, err := ldap3.Dial("tcp", fmt.Sprintf("%s:%d", ip, 389))
	if err != nil {
		return nil, err
	}
	defer dial.Close()

	err = dial.Bind(userDN, pwd)
	if err != nil {
		panic(err)
	}

	searchRequest := ldap3.NewSearchRequest(
		dn,
		ldap3.ScopeWholeSubtree, ldap3.NeverDerefAliases, 0, 0, false,
		"(objectclass=computer)",
		[]string{"cn", "distinguishedName", "ms-Mcs-AdmPwd", "ms-Mcs-AdmPwdExpirationTime"},
		nil,
	)

	search, err := dial.SearchWithPaging(searchRequest, 5)
	if err != nil {
		logger.Errorf("search(%s) computer list err:%v", ip, err)
		return nil, err
	}

	userList := []model.MachineUser{}

	for _, entry := range search.Entries {
		cn := entry.GetAttributeValue("cn")
		dn := entry.GetAttributeValue("distinguishedName")
		ex := entry.GetAttributeValue("ms-Mcs-AdmPwdExpirationTime")
		fileTimeStamp := util.Atoll(ex)
		lastLogonTime := utime.FileTime2Time(fileTimeStamp)

		userList = append(userList, model.MachineUser{
			DN:          dn,
			Domain:      domainName,
			MachineName: cn,
			MachinePwd:  "",
			Remark:      "",
			ExpiredAt:   lastLogonTime,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		})
	}

	return userList, nil
}

func getPWD(ip, userDN, dn, pwd, cName string) (string, string) {
	dial, err := InitLdapConn(ip, userDN, pwd)
	defer dial.Close()

	//for _, cName := range computerList {
	searchPWDRequest := ldap3.NewSearchRequest(
		dn,
		ldap3.ScopeWholeSubtree, ldap3.NeverDerefAliases, 0, 0, false,
		fmt.Sprintf("(&(objectCategory=computer)(samaccountname=*%s*))", cName),
		[]string{"ms-Mcs-AdmPwd", "ms-Mcs-AdmPwdExpirationTime"},
		nil,
	)

	search, err := dial.SearchWithPaging(searchPWDRequest, 20)
	if err != nil {
		logger.Errorf("get pwd err:%v", err)
		return "", ""
	}

	if len(search.Entries) == 0 {
		return "", ""
	}

	entry := search.Entries[0]

	return entry.GetAttributeValue("ms-Mcs-AdmPwd"), entry.GetAttributeValue("ms-Mcs-AdmPwdExpirationTime")
}

func InitLdapConn(ip string, userDN string, pwd string) (*ldap3.Conn, error) {
	dial, err := ldap3.Dial("tcp", fmt.Sprintf("%s:%d", ip, 389))
	if err != nil {
		panic(err)
	}

	dial.SetTimeout(time.Second * 10)
	err = dial.Bind(userDN, pwd)
	if err != nil {
		defer dial.Close()
		panic(err)
	}
	return dial, err
}
