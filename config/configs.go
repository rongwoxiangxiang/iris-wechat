package config

type configs struct {
	ServiceConfig  *ServiceConfig `yaml:"service"`
	DataBaseConfig *DataBaseConfig `yaml:"database"`
	RedisConfig  *RedisConfig `yaml:"redis"`
	OthersConfig *OthersConfig  `yaml:"others"`
}

var (
	local_conf *configs
)

func GetConfigs() *configs {
	if local_conf == nil {
		getDefaultConfigs()
	}
	return local_conf
}

func getDefaultConfigs(){
	local_conf = new(configs)
	setServiceConfigDefault()
	setDataBaseConfigDefault()
	setCacheConfigDefault()
}