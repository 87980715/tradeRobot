package controllers

import (
	"github.com/astaxie/beego"
	"tradeRobot/robot/models"
	"github.com/astaxie/beego/logs"
	"strings"
	"encoding/json"
	"github.com/jinzhu/gorm"
	"strconv"
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
	symbol := tempSymbol[0] + "_" + "CNZ"

	db, err := LoadRobotDB()

	if err != nil {
		logs.Error("load robotDB failed err:", err)
		result["code"] = 1001
		result["message"] = "操作失败"
	}

	defer db.Close()
	var tradeResult []models.ZGTradeResults

	db.Model(&models.ZGTradeResults{}).Where(&models.ZGTradeResults{Symbol: symbol}).Find(&tradeResult)

	bytes, _ := json.Marshal(tradeResult)

	result["results"] = string(bytes)
}

func LoadRobotDB() (*gorm.DB, error) {
	conf := new(DBConfig)
	conf.Path = "127.0.0.1"
	conf.Port = 3306
	conf.DbName = "tradeRobot"
	conf.User = "root"
	conf.Password = "root"

	str := conf.User + ":" +
		conf.Password + "@tcp(" +
		conf.Path + ":" +
		strconv.FormatUint(uint64(conf.Port), 10) + ")/" +
		conf.DbName + "?" +
		conf.Charset + "&parseTime=True&loc=Local"
	db, err := gorm.Open("mysql", str)
	if err != nil {
		return nil, err
	}
	return db, nil
}
