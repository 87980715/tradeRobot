package initialize

import (
	"github.com/astaxie/beego/logs"
	"fmt"
	"tradeRobot/models"
	)

func InitRobot() (err error){

	filename := "./conf/app.conf"

	err = InitConf("ini", filename)
	if err != nil {
		fmt.Printf("init conf failed, err:%v\n", err)
		panic("init conf failed")
		return
	}

	err = InitLogger()
	if err != nil {
		fmt.Printf("init logger failed---, err:%v\n", err)
		panic("init logger failed")
		return
	}

	err = InitRedis()
	if err != nil {
		fmt.Printf("init redis failed, err:%v\n", err)
		logs.Error("init redis failed")
		return
	}

	initSymbol()
	logs.Info("OKexSymbols:%s,HuobiSymbols:%s",models.OKexSymbols,models.HuobiSymbols)

	initAccounts()
	logs.Info("init Accounts success")
	return
}
