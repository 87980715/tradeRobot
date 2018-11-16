package controllers

import (
	"github.com/astaxie/beego"
	"tradeRobot/robot/models"
	"strings"
	"github.com/jinzhu/gorm"
	"strconv"
	"github.com/astaxie/beego/logs"
)

type ZGFinishedController struct {
	beego.Controller
}

type DBConfig struct {
	User     string `default:"root"`
	Password string `default:""`
	Path     string `required:"true"`
	Port     uint   `default:"3306"`
	DbName   string `required:"true"`
	Charset  string `default:"utf8"`
}

func (c *ZGFinishedController) ZGFinished() {

	result := make(map[string]interface{})

	result["code"] = 0
	result["message"] = "操作成功"

	defer func() {
		c.Data["json"] = result
		c.ServeJSON()
	}()

	s := strings.ToUpper(c.GetString("symbol"))
	tempSymbol := strings.Split(s, "-")
	symbol := tempSymbol[0] + "_" + "CNT"

	db,err := LoadRobotDB()
	if err != nil  {
		result["code"] = 1001
		result["message"] = "操作失败"
		return
	}

	defer db.Close()

	var tradeResult []models.ZGTradeResults
	db.Model(&models.ZGTradeResults{}).Where(&models.ZGTradeResults{Symbol: symbol}).Find(&tradeResult)

	result["results"] = tradeResult
}

func LoadRobotDB() (*gorm.DB, error) {
	conf := new(DBConfig)
	conf.Path = "47.244.14.215"
	conf.Port = 3306
	conf.DbName = "tradeRobot"
	conf.User = "robot"
	conf.Password = "robot"

	str := conf.User + ":" +
		conf.Password + "@tcp(" +
		conf.Path + ":" +
		strconv.FormatUint(uint64(conf.Port), 10) + ")/" +
		conf.DbName + "?" +
		conf.Charset + "&parseTime=True&loc=Local"
	db, err := gorm.Open("mysql", str)
	if err != nil {
		logs.Error("gorm open db failed")
		return nil, err
	}
	return db, nil
}
