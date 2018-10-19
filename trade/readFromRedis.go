package trade

import (
	"tradeRobot/initialize"
	"github.com/astaxie/beego/logs"
	"github.com/garyburd/redigo/redis"
	"fmt"
	"time"
	"tradeRobot/models"
	"tradeRobot/utils"
	"encoding/json"
)

func ReadTickersFromRedis(coinNames map[int]string, toCionNames map[int]string) {

	// tempMap := make(map[string]string) 在这申明map，放入通道，和从通道读取，会存在数据竞争，
	// 因为写入的map和从通道读取的map最终都会指向同一个内存地址。
Loop:
	for {
		time.Sleep(1 * time.Second)
		conn := initialize.RedisPool.Get()
		for _, toCionName := range toCionNames {
			for _, cionName := range coinNames {
				args := cionName + "_" + toCionName + "_" + "tickers"
				reply, err := conn.Do("RPOP", args)
				if err != nil {
					logs.Error("RPOP tickers failed, err:", err)
					conn.Close()
					continue Loop
				}
				data, err := redis.String(reply, err)
				if err != nil {
					conn.Close()
					continue Loop
				}

				tempMap := make(map[string]string)
				key := cionName + "_" + toCionName
				tempMap[key] = data
				fmt.Println("tickers:", tempMap)
				//TickersChan <- tempMap
			}
		}
		conn.Close()
	}
}

func ReadDepthsFromRedis() {
	Loop:
	for {
		time.Sleep(time.Millisecond*130)
		conn := initialize.RedisPool.Get()
		// 读取okex买卖深度
		for _, symbol := range models.AllSymbols {
			reply, err := conn.Do("RPOP", symbol)
			data, err := redis.String(reply, err)
			if err == redis.ErrNil {
				//logs.Error("return nil")
				conn.Close()
				continue Loop
			}
			tempMap := make(map[string]string)
			tempMap[symbol] = data
			DepthsChan <- tempMap
		}
		conn.Close()
	}
}

func ReadOrdersForOKex() {
Loop:
	for {
		conn := initialize.RedisPool.Get()
		reply, err := conn.Do("RPOP", "OKex")
		str, err := redis.String(reply, err)
		if err != redis.ErrNil {
			logs.Error("RPOP orders failed")
			conn.Close()
			continue Loop
		}
		okexPostDataLimit := &utils.OKexPostDataLimit{}
		err = json.Unmarshal([]byte(str),okexPostDataLimit)
		if err != nil {
			logs.Error("json.Unmarshal okexPostDataLimit failed err:",err)
		}

		fmt.Println("OKexOrders:", okexPostDataLimit)
		conn.Close()
	}
}

