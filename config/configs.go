package config

//其他配置
type OtherConfigs struct {
	*QQAiServer  `yaml:"ai"`
}

type QQAiServer struct {
	AppId string `yaml:"qq_ai_app_id"`
	AppKey string  `yaml:"qq_ai_app_key"`
}