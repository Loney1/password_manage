package test

import (
	"context"
	"os"
	"testing"
	"time"

	v1 "adp_backend/api/v1"
	"adp_backend/common"
	"adp_backend/util"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

const (
	username = "admin"
	grpcAddr = "192.168.30.245:5000"
)

var ADMCli *ADMGrpcClient

type ADMGrpcClient struct {
	cli v1.ADMClient
	ctx context.Context
}

func TestMain(m *testing.M) {
	conn, err := grpc.Dial(grpcAddr, grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	defer func() {
		_ = conn.Close()
	}()

	client := v1.NewADMClient(conn)
	exp := time.Now().AddDate(0, 0, 90).Unix()
	token, err := util.GenerateToken(username, common.RoleMgr, common.PrivSuper, exp)
	if err != nil {
		panic(err)
	}

	md := metadata.Pairs("authorization", "Bearer "+token)
	ctx := metadata.NewOutgoingContext(context.Background(), md)

	ADMCli = &ADMGrpcClient{
		cli: client,
		ctx: ctx,
	}
	os.Exit(m.Run())
}
