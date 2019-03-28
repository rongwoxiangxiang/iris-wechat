package config

import (
	"github.com/kataras/iris"
)


func init()  {

}

func Runner() iris.Runner {
	return iris.Addr("localhost:8080")
}
func Configuration() iris.Configurator {
	return iris.WithConfiguration(iris.YAML("config/iris.yml"))
}

func Server() (iris.Runner, iris.Configurator) {
	return Runner(),Configuration()
}


