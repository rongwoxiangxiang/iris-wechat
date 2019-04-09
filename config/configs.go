package config

//其他配置
type OtherConfigs struct {
	*QQAiServer  `yaml:"ai"`
	*ElasticSearchServer  `yaml:"es"`
}

type QQAiServer struct {
	AppId string `yaml:"qq_ai_app_id"`
	AppKey string  `yaml:"qq_ai_app_key"`
}

type ElasticSearchServer struct {
	Host string `yaml:"host"`
}