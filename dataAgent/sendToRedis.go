package dataAgent

import (
	"tradeRobot/initialize"
	"github.com/astaxie/beego/logs"
)

func SendTickersToRedis(data, cionName string, toCionNames map[int]string) {

	conn := initialize.RedisPool.Get()
	defer conn.Close()
	for _, toCionName := range toCionNames {
		args := cionName + "_" + toCionName + "_" + "tickers"
		_, err := conn.Do("LPUSH", args, data)
		if err != nil {
			logs.Error("LPUSH tickers failed, err:", err)
		}
	}
}

func SendDepthsToRedis(data, cionName string, toCionNames map[int]string) {

	conn := initialize.RedisPool.Get()
	defer conn.Close()
	for _, toCionName := range toCionNames {
		args := cionName + "_" + toCionName + "_" + "depths"
		_, err := conn.Do("LPUSH", args, data)
		if err != nil {
			logs.Error("LPUSH depth failed, err:", err)
		}
	}
}
