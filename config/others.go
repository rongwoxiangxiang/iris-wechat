package config

type OthersConfig struct {
	*QQAiServerConfig `yaml:"qqai"`
	*ElasticSearchServerConfig `yaml:"elsticsearch"`
}

type QQAiServerConfig struct {
	AppId string `yaml:"qq_ai_app_id"`
	AppKey string  `yaml:"qq_ai_app_key"`
}

type ElasticSearchServerConfig struct {
	Host string `yaml:"host"`
	Port string `yaml:"port"`
}
