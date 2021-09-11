// author: s0nnet
// time: 2020-09-01
// desc:

package main

import (
	"os"

	"adp_backend/common"
	"adp_backend/config"
	"adp_backend/service"

	logger "github.com/sirupsen/logrus"
)

func main() {
	logger.Info("starting apiserver for ADM")

	confPath := os.Getenv("API_SRV_CONF_PATH")
	if confPath == "" {
		// dev mode, configure path load from file
		logger.Info("load configure from file")
		confPath = common.CONF_PATH
	}

	logger.Infof("load configure from %s", confPath)
	env, err := config.Init(confPath)
	if err != nil {
		logger.Panic(err)
	}

	s := service.New(env)

	go s.SignalHandler()

	if err := s.Start(env.Cfg.GrpcSrv.Address); err != nil {
		logger.Panic(err)
	}
}
