package main

import (
	"adp_backend/util"
	"flag"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"os"
	"reflect"
	"strings"
	"time"

	"adp_backend/common"
	"adp_backend/config"
	"adp_backend/model"
	"adp_backend/server"
	"adp_backend/service"

	logger "github.com/sirupsen/logrus"
)

var env *config.Env

//./build -model domain -DomainName  -DCHostName  -DNS  -UserName  -UserDN  -UserPWD
func main() {
	modelType := flag.String("model", "", "change model type,domain or initUser")
	DomainName := flag.String("DomainName", "", "")
	DCHostName := flag.String("DCHostName", "", "")
	DNS := flag.String("DNS", "", "")
	UserName := flag.String("UserName", "", "")
	UserDN := flag.String("UserDN", "", "")
	Pwd := flag.String("UserPWD", "", "")
	//Secret := flag.String("Secret", "", "")

	flag.Parse()

	if *modelType == "" {
		fmt.Println("Please select operation model")
		return
	}

	switch *modelType {
	case "domain":
		if isNil(DomainName, DCHostName, DNS, UserName, UserDN, Pwd) {
			fmt.Println("Please enter the correct parameters, -DomainName  -DCHostName  -DNS  -UserName  -UserDN  -UserPWD")
			return
		}
		split := strings.SplitN(*DomainName, ".", 2)
		encrypt, err := util.PasswordEncrypt(*Pwd)
		if err != nil {
			logger.Panicf("pass word encrypt err:%v", err)
		}
		domain := model.Domain{
			Name:       *DomainName,
			DCHostName: *DCHostName,
			DNS:        *DNS,
			DN:         fmt.Sprintf("DC=%s,DC=%s", split[0], split[1]),
			UserName:   *UserName,
			Password:   encrypt,
			UserDN:     *UserDN,
			Status:     0,
			ErrMsg:     "",
			CreatedAt:  time.Time{},
		}

		_, err = service.InitLdapConn(domain.DNS, domain.UserDN, *Pwd)
		if err != nil {
			logger.Panicf("Domain status failed err:%v", err)
		}

		err = server.AddDomain(env.MysqlCli, domain)
		if err != nil {
			logger.Panicf("add domain err:%v", err)
		}
	case "initUser":
		if isNil(UserName, Pwd) {
			fmt.Println("Please enter the correct parameters")
			return
		}
		plainPwd, err := bcrypt.GenerateFromPassword([]byte(*Pwd), bcrypt.MinCost)
		if err != nil {
			fmt.Println("服务器内部错误")
			return
		}
		user := model.User{
			ID:           1,
			UserName:     *UserName,
			Password:     string(plainPwd),
			PassStrength: "high",
			Secret:       "JYYFQWCDGQYTGQ2WI5EUYNBQHFETSSCX",
			MfaStatus:    "enable",
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		}

		err = server.CreateUser(env, user)
		if err != nil {
			fmt.Printf("add user err:%v", err)
			return
		}

	}

}

func isNil(values ...interface{}) bool {
	for _, value := range values {
		valueOf := reflect.ValueOf(value)
		if valueOf.IsZero() {
			return true
		}
	}
	return false
}

func init() {
	confPath := os.Getenv("API_SRV_CONF_PATH")
	if confPath == "" {
		// dev mode, configure path load from file
		logger.Info("load configure from file")
		confPath = common.CONF_PATH
	}

	logger.Infof("load configure from %s", confPath)
	var err error
	env, err = config.Init(confPath)
	if err != nil {
		logger.Panic(err)
	}
}
