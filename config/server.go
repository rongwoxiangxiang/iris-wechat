package config

type ServiceConfig struct {
	Host string `yaml:"host"`
	Port string `yaml:"port"`
}

func setServiceConfigDefault() {
	local_conf.ServiceConfig = &ServiceConfig{Host:"0.0.0.0", Port:"8888"}
}