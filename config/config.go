// author: s0nnet
// time: 2020-09-01
// desc:

package config

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"runtime/pprof"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type LogCfg struct {
	LogPath   string `yaml:"LogPath"`
	LogLevel  string `yaml:"LogLevel"`
	IsStdOut  bool   `yaml:"IsStdOut"`
	IsPProf   bool   `yaml:"IsPProf"`
	PathPProf string `yaml:"PathPProf"`
}

type GrpcSrvCfg struct {
	Address string `yaml:"Address"`
}

type MySQLCfg struct {
	User        string        `yaml:"User"`
	Passwd      string        `yaml:"Passwd"`
	Host        string        `yaml:"Host"`
	DbName      string        `yaml:"DbName"`
	DSN         string        `yaml:"DSN"`
	Active      int           `yaml:"Active"`
	Idle        int           `yaml:"Idle"`
	IdleTimeout time.Duration `yaml:"IdleTimeout"`
}

type LicenceCfg struct {
	Ip   string `yaml:"Ip"`
	IsVm bool   `yaml:"IsVm"`
}

type Config struct {
	ProjectName string     `yaml:"ProjectName"`
	Licence     LicenceCfg `yaml:"Licence"`
	Log         LogCfg     `yaml:"Log"`
	GrpcSrv     GrpcSrvCfg `yaml:"GrpcSrv"`
	MySQL       MySQLCfg   `yaml:"MySQL"`
}

type Env struct {
	Cfg      *Config
	MysqlCli *gorm.DB
}

var (
	cpuProfilingFile,
	memProfilingFile,
	blockProfilingFile,
	goroutineProfilingFile,
	threadCreateProfilingFile *os.File
)

func InitLog(setting *Config) error {
	logrus.SetFormatter(&logrus.JSONFormatter{})
	if setting.Log.LogLevel == "" {
		return nil
	}

	lvl, err := logrus.ParseLevel(setting.Log.LogLevel)
	if err != nil {
		return err
	}
	logrus.SetLevel(lvl)

	if setting.Log.IsStdOut {
		logrus.SetOutput(os.Stdout)
		logrus.SetFormatter(&logrus.TextFormatter{})
	}
	if setting.Log.LogPath != "" {
		logFile := filepath.Join(setting.Log.LogPath, setting.ProjectName+"_stdout.log")
		file, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err == nil {
			logrus.SetOutput(file)
		}
	}

	logrus.SetReportCaller(true)

	return nil
}

func InitDebugPProf(setting *Config) error {
	if setting.Log.IsPProf {
		_, err := os.Stat(setting.Log.PathPProf)
		if err != nil && os.IsNotExist(err) {
			_ = os.MkdirAll(setting.Log.PathPProf, os.ModePerm)
		}

		pathPrefix := path.Join(setting.Log.PathPProf, fmt.Sprintf("%d", os.Getpid()))
		logrus.Infof("start pprof, and will save to %s", pathPrefix)
		cpuProfilingFile, _ = os.Create(pathPrefix + "-cpu.prof")
		memProfilingFile, _ = os.Create(pathPrefix + "-mem.prof")
		blockProfilingFile, _ = os.Create(pathPrefix + "-block.prof")
		goroutineProfilingFile, _ = os.Create(pathPrefix + "-goroutine.prof")
		threadCreateProfilingFile, _ = os.Create(pathPrefix + "-threadcreat.prof")
		_ = pprof.StartCPUProfile(cpuProfilingFile)
	}
	return nil
}

// SaveProfile try to save pprof into local file
func (e *Env) SaveProfile() {
	if e.Cfg.Log.IsPProf {
		goroutine := pprof.Lookup("goroutine")
		goroutine.WriteTo(goroutineProfilingFile, 1)
		heap := pprof.Lookup("heap")
		heap.WriteTo(memProfilingFile, 1)
		block := pprof.Lookup("block")
		block.WriteTo(blockProfilingFile, 1)
		threadcreate := pprof.Lookup("threadcreate")
		threadcreate.WriteTo(threadCreateProfilingFile, 1)
		pprof.StopCPUProfile()
	}
}

func InitMySQLClient(setting *Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local",
		setting.MySQL.User, setting.MySQL.Passwd, setting.MySQL.Host, setting.MySQL.DbName)
	db, err := gorm.Open(mysql.New(mysql.Config{
		DriverName:                "",
		DSN:                       dsn,
		Conn:                      nil,
		SkipInitializeWithVersion: false,
		DefaultStringSize:         256,
		DefaultDatetimePrecision:  nil,
		DisableDatetimePrecision:  true,
		DontSupportRenameIndex:    true,
		DontSupportRenameColumn:   true,
		DontSupportForShareClause: false,
	}), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		panic(err)
	}
	sqlDB.SetMaxIdleConns(setting.MySQL.Idle)
	sqlDB.SetMaxOpenConns(setting.MySQL.Active)
	sqlDB.SetConnMaxLifetime(time.Duration(setting.MySQL.IdleTimeout) / time.Second)

	return db, nil
}

func Init(confPath string) (*Env, error) {
	content, err := ioutil.ReadFile(confPath)
	if err != nil {
		panic(err)
	}

	var setting Config
	err = yaml.Unmarshal(content, &setting)
	if err != nil {
		panic(err)
	}

	err = InitDebugPProf(&setting)
	if err != nil {
		return nil, err
	}

	err = InitLog(&setting)
	if err != nil {
		return nil, err
	}

	mysqlCli, err := InitMySQLClient(&setting)
	if err != nil {
		return nil, err
	}

	return &Env{&setting, mysqlCli}, nil
}
