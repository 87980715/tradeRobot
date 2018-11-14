package utils

import (
	"github.com/jinzhu/gorm"
	"github.com/astaxie/beego/logs"
	"tradeRobot/robot/models"
	"strconv"
	_ "github.com/go-sql-driver/mysql"
	"database/sql"
)

var (
	RobotDB *gorm.DB
	ExchangeDB *sql.DB
)

type DBConfig struct {
	User     string `default:"root"`
	Password string `default:""`
	Path     string `required:"true"`
	Port     uint   `default:"3306"`
	DbName   string `required:"true"`
	Charset  string `default:"utf8"`
}

func LoadRobotDB() (*gorm.DB, error) {
	var err error
	dbConfig := new(DBConfig)
	dbConfig.Path = "127.0.0.1"
	//dbConfig.Path = "192.168.0.33"
	dbConfig.Port = 3306
	dbConfig.DbName = "tradeRobot"
	dbConfig.User = "root"
	dbConfig.Password = "root"

	RobotDB, err = GetDBConnection(dbConfig)
	if err != nil {
		return nil, err
	}
	return RobotDB, err
}

func InitRobotDB() (error) {
	var err error
	RobotDB, err = LoadRobotDB()

	if err != nil {
		logs.Error("loadDB failed err:", err)
	}
	if !RobotDB.HasTable(&models.HuobiTradeResults{}) {
		if err = RobotDB.Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8").CreateTable(&models.HuobiTradeResults{}).Error; err != nil {
			return err
		}
	}
	if !RobotDB.HasTable(&models.ZGTradeResults{}) {
		if err = RobotDB.Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8").CreateTable(&models.ZGTradeResults{}).Error; err != nil {
			return err
		}
	}
	logs.Error("create trade result tables success")
	return nil
}

func GetDBConnection(conf *DBConfig) (db *gorm.DB, err error) {

	str := conf.User + ":" +
		conf.Password + "@tcp(" +
		conf.Path + ":" +
		strconv.FormatUint(uint64(conf.Port), 10) + ")/" +
		conf.DbName + "?" +
		conf.Charset + "&parseTime=True&loc=Local"
	db, err = gorm.Open("mysql", str)
	if err != nil {
		logs.Error("gorm open mamria failed err:", err)
	}
	return
}

func LoadExchangeDB() (*sql.DB, error) {
	conf := new(DBConfig)
	//conf.Path = "47.99.74.117"
	conf.Path = "exchange-readonly.crdydmbkv0de.ap-northeast-1.rds.amazonaws.com"
	conf.User = "exchange"
	//conf.Password = "exchange"
	conf.Password = "7wdUVYriIyeMI2zaicpYnNy2nT6YFUUm"
	conf.Port = 3306
	conf.DbName = "trade_history"

	str := conf.User + ":" +
		conf.Password + "@tcp(" +
		conf.Path + ":" +
		strconv.FormatUint(uint64(conf.Port), 10) + ")/" +
		conf.DbName + "?" +
		conf.Charset

	db, err := sql.Open("mysql", str)
	if err != nil {
		logs.Error("sql open mysql failed err:", err)
		return nil, err
	}
	return db, nil
}
