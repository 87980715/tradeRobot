package initialize

import (
	"fmt"
	"github.com/astaxie/beego/config"
	"github.com/labstack/gommon/log"
	"sync"
)

var (
	AppConfig *Config
)

type Config struct {
	logLevel           string
	logPath            string

	RedisConf *RedisConf

	Wg     sync.WaitGroup
	RWlock sync.RWMutex
}

type RedisConf struct {
	RedisAddr        string
	RedisMaxIdle     int
	RedisMaxActive   int
	RedisIdleTimeout int
}


func InitConf(confType, filename string) (err error) {

	conf, err := config.NewConfig(confType, filename)
	if err != nil {
		fmt.Println("new config failed, err:", err)
		return
	}
	redisConf := &RedisConf{}

	AppConfig = &Config{
		RedisConf: redisConf,
	}
	AppConfig.logLevel = conf.String("logs::log_level")
	if len(AppConfig.logLevel) == 0 {
		AppConfig.logLevel = "debug"
	}

	AppConfig.logPath = conf.String("logs::log_path")
	if len(AppConfig.logPath) == 0 {
		AppConfig.logPath = "../robot/logs"
	}

	AppConfig.RedisConf.RedisAddr = conf.String("redis::server_addr")
	if len(AppConfig.RedisConf.RedisAddr) == 0 {
		err = fmt.Errorf("invalid redis addr")
		return
	}

	AppConfig.RedisConf.RedisMaxIdle, err = conf.Int("redis::redis_idle")
	if err != nil {
		log.Warn("load redis config failed")
	}
	AppConfig.RedisConf.RedisMaxActive, err = conf.Int("redis::redis_active")
	if err != nil {
		log.Warn("load redis config failed")
	}

	AppConfig.RedisConf.RedisIdleTimeout, err = conf.Int("redis::redis_idle_timeout")
	if err != nil {
		log.Warn("load redis config failed")
	}

	return
}
