package trade

import (
	"tradeRobot/robot/initialize"
	"github.com/astaxie/beego/logs"
	"github.com/garyburd/redigo/redis"
	"fmt"
	"time"
)


func ReadTradesFromRedis() {
Loop:
	for {
		time.Sleep(500*time.Millisecond)
		conn := initialize.RedisPool.Get()
		reply, err := conn.Do("RPOP", "T")
		trade, err := redis.String(reply, err)
		if err != redis.ErrNil {
			logs.Error("RPOP return nil err:",err)
			conn.Close()
			continue Loop
		}
		fmt.Println("redis:",trade)
		//TradesChan <- trade
		conn.Close()
	}
}

