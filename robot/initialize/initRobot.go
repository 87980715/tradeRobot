package initialize

import (
	"github.com/astaxie/beego/logs"
	"fmt"
	"tradeRobot/robot/utils"
	)

func InitRobot() (err error){

	filename := "../robot/conf/app.conf"

	err = InitConf("ini", filename)
	if err != nil {
		fmt.Printf("init conf failed, err:%v\n", err)
		panic("init conf failed")
		return
	}

	err = InitLogger()
	if err != nil {
		fmt.Printf("init logger failed, err:%v\n", err)
		panic("init logger failed")
		return
	}

	err = utils.InitRobotDB()
	if err != nil {
		fmt.Printf("init mysqldb failed, err:%v\n", err)
		panic("init db failed")
		return
	}

	err = InitRedis()
	if err != nil {
		logs.Error("init redis failed")
		return
	}

	initAccounts()
	logs.Info("init Accounts success")

	return
}
