package bootstrap

import (
	"fmt"
	"github.com/kataras/iris"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"iris/common"
	"iris/config"
	"iris/routers"
	"log"
	"time"
)

func init()  {
	initYamlConfig()
}

func Run()  {
	app := iris.New()
	routers.Routes(app)
	app.Run(runner(), configuration())
}

func runner() iris.Runner {
	return iris.Addr(fmt.Sprintf("%s:%s",
		config.GetConfigs().ServiceConfig.Host , config.GetConfigs().ServiceConfig.Port))
}

func configuration() iris.Configurator {
	serverConfigExist := common.CheckFileExist("server.yml")
	if serverConfigExist {
		configurate := iris.YAML("server.yml")
		return iris.WithConfiguration(configurate)
	}
	return iris.WithConfiguration(iris.Configuration{
		TimeFormat:                        time.RFC3339,
		Charset:                           "UTF-8",
	})
}

func initYamlConfig()  {
	serverConfigExist := common.CheckFileExist("application.yml")
	if serverConfigExist {
		data, err := ioutil.ReadFile("application.yml")
		if err != nil {
			log.Fatalf("error: %v", err)
		}
		err = yaml.Unmarshal([]byte(data), config.GetConfigs())
		if err != nil {
			log.Fatalf("error: %v", err)
		}
	}

}