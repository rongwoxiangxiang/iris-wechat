package config

import (
	"fmt"
	"github.com/go-xorm/xorm"
	"log"
	"strings"
)

type MysqlConfig struct {
	Host string `yaml:"host"`
	Port int64  `yaml:"port"`
	User string `yaml:"user"`
	Pass string `yaml:"pass"`
	Name string `yaml:"name"`
	Charset string `yaml:"charset"`
	ShowSql bool `yaml:"showSql"`
}

var (
	orm *xorm.Engine
	dataSourceName string
)

func initDb() {
	var err error
	configs := GetConfigs()
	if strings.ToLower(configs.DriverName) == "mysql" {
		dataSourceName = getMysqlConfig(configs)
	}
	orm, err = xorm.NewEngine(configs.DriverName, dataSourceName)
	orm.ShowSQL(configs.MysqlConfig.ShowSql)
	if err != nil {
		log.Fatalf("orm failed to initialized: %v", err)
	}
}

func getMysqlConfig(configs *configs) string {
	if configs.MysqlConfig.Charset == "" {
		configs.MysqlConfig.Charset = "utf8"
	}
	return fmt.Sprintf("%s:%s@(%s:%d)/%s?charset=%s",
		configs.MysqlConfig.User,
		configs.MysqlConfig.Pass,
		configs.MysqlConfig.Host,
		configs.MysqlConfig.Port,
		configs.MysqlConfig.Name,
		configs.MysqlConfig.Charset,
	)
}

func GetDb() *xorm.Engine {
	return orm
}