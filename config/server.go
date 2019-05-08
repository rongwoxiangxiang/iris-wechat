package config
//
//import (
//	"fmt"
//	_ "github.com/go-sql-driver/mysql"
//	"github.com/kataras/iris"
//	"gopkg.in/yaml.v2"
//	"io/ioutil"
//	"log"
//	"os"
//)
//
//
//
//type configs struct {
//	ServiceConfig *ServiceConfig `yaml:"service"`
//	DriverName string `yaml:"driverName"`
//	MysqlConfig *MysqlConfig `yaml:"database"`
//	RedisConfig *RedisConfig `yaml:"redis"`
//	OtherConfigs *OtherConfigs `yaml:"others"`
//}
//
//var (
//    conf *configs
//)
//
//func init()  {
//	initYamlConfig()
//	initDb()
//}
//
//
//func Run(app *iris.Application)  {
//	app.Run(runner(), configuration())
//}
//
//func runner() iris.Runner {
//	return iris.Addr(fmt.Sprintf("%s:%d",conf.ServiceConfig.Host , conf.ServiceConfig.Port))
//}
//func configuration() iris.Configurator {
//	configurate := iris.YAML("server.yml")
//	return iris.WithConfiguration(configurate)
//}
//
//func initYamlConfig()  {
//	data, err := ioutil.ReadFile(os.Getenv("GOPATH")+"/src/iris/application.yml")
//	if err != nil {
//		log.Fatalf("error: %v", err)
//	}
//	err = yaml.Unmarshal([]byte(data), &conf)
//	if err != nil {
//		log.Fatalf("error: %v", err)
//	}
//}
//
//func GetConfigs() *configs {
//	return conf
//}