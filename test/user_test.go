package test

import (
	v1 "adp_backend/api/v1"
	"testing"
)

func TestGetUserName(t *testing.T) {
}

func TestLogin(t *testing.T) {
	req := v1.LoginReq{
		Username: "admin",
		Password: "admin12345",
		TotpCode: "678043",
	}

	login, err := ADMCli.cli.Login(ADMCli.ctx, &req)
	if err != nil {
		panic(err)
	}

	t.Log(login)
}
