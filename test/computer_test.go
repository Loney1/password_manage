package test

import (
	v1 "adp_backend/api/v1"
	"fmt"
	"testing"
)

func TestGetMachinePwd(t *testing.T) {
	req := v1.GetMachinePwdReq{
		Id:     47,
		HasMfa: "238743",
	}

	pwd, err := ADMCli.cli.GetMachinePwd(ADMCli.ctx, &req)
	if err != nil {
		panic(err)
	}

	fmt.Println(pwd)
}
func TestListMachinePwd(t *testing.T) {
	req := v1.ListMachineNameReq{
		PageIdx:    2,
		PageSize:   10,
		SearchName: "",
		StartTime:  "",
		EndTime:    "",
		SortTime:   0,
	}

	pwd, err := ADMCli.cli.ListMachineName(ADMCli.ctx, &req)
	if err != nil {
		panic(err)
	}

	fmt.Println(pwd)
}
