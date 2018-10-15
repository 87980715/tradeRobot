package initialize

import (
	"fmt"
	"github.com/astaxie/beego/config"
	"github.com/labstack/gommon/log"
	"sync"
	"time"
)

var (
	AppConfig *Config
)

type Config struct {
	logLevel           string
	logPath            string


	TradeConf *TradeConf
	RedisConf *RedisConf
	EtcdConf  *EtcdConf

	Wg     sync.WaitGroup
	RWlock sync.RWMutex
}

type RedisConf struct {
	RedisAddr        string
	RedisMaxIdle     int
	RedisMaxActive   int
	RedisIdleTimeout int
}

type EtcdConf struct {
	EtcdAddr string
	//Timeout           int
	EtcdSecKeyPrefix  string
	EtcdSecProductKey string
}

type TradeConf struct {
	TradeAmountMutiple float64
	TradeInsprctiontime time.Duration
}

func InitConf(confType, filename string) (err error) {

	conf, err := config.NewConfig(confType, filename)
	if err != nil {
		fmt.Println("new config failed, err:", err)
		return
	}
	redisconf := &RedisConf{}
	AppConfig = &Config{
		RedisConf: redisconf,
	}
	AppConfig.logLevel = conf.String("logs::log_level")
	if len(AppConfig.logLevel) == 0 {
		AppConfig.logLevel = "debug"
	}

	AppConfig.logPath = conf.String("logs::log_path")
	if len(AppConfig.logPath) == 0 {
		AppConfig.logPath = "./logs"
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

	AppConfig.TradeConf.TradeAmountMutiple, err = conf.Float("redis::trade_amount_multiple")
	if err != nil {
		log.Warn("load trade config failed")
	}

	//AppConfig.EtcdAddr = conf.String("etcd::addr")
	//if len(AppConfig.etcdAddr) == 0 {
	//	err = fmt.Errorf("invalid etcd addr")
	//	return
	//}

	//appConfig.etcdKey = conf.String("etcd::configKey")
	//if len(appConfig.etcdKey) == 0 {
	//	err = fmt.Errorf("invalid etcd key")
	//	return
	//}

	return
}
