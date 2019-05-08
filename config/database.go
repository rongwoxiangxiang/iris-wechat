package config

import (
	"fmt"
	"github.com/go-xorm/xorm"
	_ "github.com/go-sql-driver/mysql"
	"log"
)

//
//import (
//	"fmt"
//	"github.com/go-xorm/xorm"
//	"log"
//	"strings"
//)

type DataBaseConfig struct {
	Host string `yaml:"host"`
	Port string `yaml:"port"`
	User string `yaml:"user"`
	Pass string `yaml:"pass"`
	Name string `yaml:"name"`
	Charset string `yaml:"charset"`
	ShowSql bool `yaml:"showSql"`
}

func setDataBaseConfigDefault()  {
	local_conf.DataBaseConfig = &DataBaseConfig{}
	local_conf.DataBaseConfig.Host = "127.0.0.1"
	local_conf.DataBaseConfig.Port = "3306"
	local_conf.DataBaseConfig.User = "root"
	local_conf.DataBaseConfig.ShowSql = true
}

var (
	orm *xorm.Engine
	dataSourceName string
)

func initMysqlDb() {
	var err error
	configs := GetConfigs()
	dataSourceName = getMysqlConfig(configs)
	orm, err = xorm.NewEngine("mysql", dataSourceName)
	orm.ShowSQL(configs.DataBaseConfig.ShowSql)
	if err != nil {
		log.Fatalf("orm failed to initialized: %v", err)
	}
}

func getMysqlConfig(configs *configs) string {
	if configs.DataBaseConfig.Charset == "" {
		configs.DataBaseConfig.Charset = "utf8"
	}
	return fmt.Sprintf("%s:%s@(%s:%s)/%s?charset=%s",
		configs.DataBaseConfig.User,
		configs.DataBaseConfig.Pass,
		configs.DataBaseConfig.Host,
		configs.DataBaseConfig.Port,
		configs.DataBaseConfig.Name,
		configs.DataBaseConfig.Charset,
	)
}

func GetDb() *xorm.Engine {
	if orm == nil {
		initMysqlDb()
	}
	return orm
}