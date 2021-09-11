// author: s0nnet
// time: 2020-09-01
// desc:

package service

import (
	"context"
	"fmt"
	"regexp"
	"strconv"
	"time"

	v1 "adp_backend/api/v1"
	"adp_backend/common"
	"adp_backend/model"
	"adp_backend/server"
	"adp_backend/util"

	logger "github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	tokenExpired    = 60 * 4
	fileMaxSize     = 512 * 1024
	needChangePwdTm = 90 * 24 * time.Hour //needChangePwdTm 提醒用户修改密码周期
)

var UserLoginCountInfo model.UserBucket

func (s *ADMServiceV1) Login(ctx context.Context, in *v1.LoginReq) (*v1.LoginReply, error) {
	user, err := server.GetUserByUserName(s.env, in.Username)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "无效的用户名或密码")
	}

	clt := UserLoginCountInfo.Get(in.Username) //用户登录次数
	if clt.LoginErrCount >= common.LoginErrorCount {
		return nil, status.Error(codes.PermissionDenied, "登录无效锁定五分钟")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(in.Password)) //判断密码是否正确
	if err != nil {
		_ = UserLoginCountInfo.SetLoginErrCount(in.Username, 1)
		return nil, status.Error(codes.Unauthenticated, "无效的用户名或密码")
	}

	if user.MfaStatus == "enable" && in.TotpCode == "" {
		return nil, status.Error(codes.PermissionDenied, "未输入二次认证")
	}

	// 新增二次验证判断
	if user.MfaStatus == "enable" {
		totpCode, err := strconv.Atoi(in.TotpCode)
		if err != nil {
			logger.Errorf("string to int failed：%v", err)
			return nil, status.Error(codes.Unauthenticated, "输入的验证码有误")
		}

		check := util.TotpCheck(user.Secret, totpCode)
		if !check {
			_ = UserLoginCountInfo.SetLoginErrCount(in.Username, 1)

			return nil, status.Error(codes.Unauthenticated, "验证码错误")
		}
	}

	//如果登录成功, 错误变为0
	_ = UserLoginCountInfo.SetLoginErrCount(in.Username, 0)

	// 生产jwt-token
	exp := time.Now().Add(time.Minute * tokenExpired).Unix()
	token, err := util.GenerateToken(in.Username, user.Role, user.Pri, exp)

	//写最后登录过期时间
	if err != nil {
		return nil, status.Error(codes.Internal, "验证生成异常")
	}

	var needChangePwd bool
	t := user.UpdatedAt
	if t.IsZero() {
		t = user.CreatedAt
	}

	if time.Since(t) >= needChangePwdTm {
		logger.Debugf("用户[%s]已有90天没有更新密码", user.UserName)
		needChangePwd = true
	}

	userInfo := v1.LoginReply{
		Id:            user.ID,
		Username:      user.UserName,
		Role:          user.Role,
		Pri:           user.Pri,
		Mobile:        user.Mobile,
		Email:         user.Email,
		Remark:        user.Remark,
		Token:         token,
		NeedChangePwd: needChangePwd,
	}
	return &userInfo, nil
}

func (s *ADMServiceV1) Logout(ctx context.Context, in *v1.LogoutReq) (*v1.LogoutReply, error) {
	// passed
	return &v1.LogoutReply{
		Result: common.RESP_SUCCESS,
	}, nil
}

func (s *ADMServiceV1) ListUser(ctx context.Context, in *v1.ListUserReq) (*v1.ListUserReply, error) {
	//username := s.GetUser(ctx)
	var username = ""
	//如果是管理员查询子用户列表，则清空查询条件
	//if s.IsSuper(ctx) && !in.IsSelf {
	//	username = ""
	//}

	var limit, offset = in.PageSize, in.PageSize * (in.PageIndex - 1) // 每页的数据量, 总偏移量
	res, total, err := server.FindAllUser(s.env, limit, offset, in.Search, username)
	if err != nil {
		logger.Error("Query ListUser:%v", err)
		return nil, status.Error(codes.Internal, "查询用户列表发生错误")
	}

	ret := v1.ListUserReply{}

	for _, r := range res {
		ret.List = append(ret.List,
			&v1.ListUserReply_Details{
				UserId:       r.ID,
				UserName:     r.UserName,
				PassStrength: r.PassStrength,
				Role:         r.Role,
				Pri:          r.Pri,
				Mobile:       r.Mobile,
				Email:        r.Email,
				Remark:       r.Remark,
				CreateTM:     r.CreatedAt.String(),
				HasMfa:       r.MfaStatus == "enable", //是否开启二次认证
				UpdateTm: func(u *model.User) string {
					if u.UpdatedAt.IsZero() {
						return u.CreatedAt.String()
					}
					return u.UpdatedAt.String()
				}(&r),
				RealName:   r.RealName,
				Address:    r.Address,
				Department: r.Department,
				Post:       r.Post,
			})
	}
	ret.Page = &v1.ModelPage{PageSize: in.PageSize, PageIdx: in.PageIndex, Total: int32(total)}
	if (limit + offset) < int32(total) {
		ret.Exhausted = false
	} else {
		ret.Exhausted = true
	}
	return &ret, nil
}

func (s *ADMServiceV1) UpdateUser(ctx context.Context, in *v1.UpdateUserReq) (*v1.UpdateUserReply, error) {
	_, err := server.GetUserByUserName(s.env, in.UserName)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "无效的用户名或密码")
	}

	user, err := server.GetUserByUserName(s.env, in.UserName)
	if err != nil {
		logger.Error("Find user by name wrong:%v", err)
		return nil, status.Error(codes.Unauthenticated, "获取用户信息失败")
	}

	user, err = server.UpdateUser(s.env, in.UserId, in.UserName, in.Passwd, in.Mobile, in.Email, in.Remark)
	if err != nil {
		logger.Error("Update user wrong:%v", err)
		return nil, status.Error(codes.Unauthenticated, "更新用户信息失败")
	}
	fmt.Println(user)

	return &v1.UpdateUserReply{Result: common.RESP_SUCCESS}, nil
}

func (s *ADMServiceV1) DeleteUser(ctx context.Context, in *v1.DeleteUserReq) (*v1.DeleteUserReply, error) {
	user, err := server.DeleteUser(s.env, in.UserId)

	if err != nil {
		logger.Error("DeleteUser Error:%v", err)
		return nil, status.Error(codes.Internal, "删除用户失败")
	}
	fmt.Println(user)

	return &v1.DeleteUserReply{Result: common.RESP_SUCCESS}, nil
}

func (s *ADMServiceV1) AddUser(ctx context.Context, in *v1.AddUserReq) (*v1.AddUserReply, error) {
	if in.UserName == "" || in.Passwd == "" {
		return nil, status.Error(codes.Unauthenticated, "密码和用户名不能为空")
	}

	_, err := server.GetUserByUserName(s.env, in.UserName)
	if err == nil {
		return nil, status.Error(codes.Unauthenticated, "用户名已存在")
	}

	if (len(in.Mobile)) != 11 {
		return nil, status.Error(codes.InvalidArgument, "手机号长度应为11位")
	}

	if len(in.Email) > 100 {
		return nil, status.Error(codes.InvalidArgument, "邮箱长度过长, 请输入100位以下长度")
	}

	if len(in.Remark) > 1000 {
		return nil, status.Error(codes.InvalidArgument, "备注长度应小于1000位")
	}

	if CheckPassword(in.Passwd) == false {
		return nil, status.Error(codes.InvalidArgument, "密码长度应大于12位")
	}

	user, err := server.AddUser(s.env, in.UserName, in.Passwd, in.Mobile, in.Email, in.Remark)
	if err != nil {
		logger.Error("AddUser Error:%v", err)
		return nil, status.Error(codes.Internal, "新增用户失败")
	}
	fmt.Println(user)

	return &v1.AddUserReply{Result: common.RESP_SUCCESS}, nil
}

func CheckPassword(password string) (b bool) {
	var str = "^[a-zA-Z0-9!@#$%^&*()_+]{12,}$"
	if ok, _ := regexp.MatchString(str, password); !ok {
		return false
	}

	return true
}
