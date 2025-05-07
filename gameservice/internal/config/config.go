package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	Server struct {
		Port string `mapstructure:"port"`
		Host string `mapstructure:"host"`
	} `mapstructure:"server"`

	Logger struct {
		FileName string `mapstructure:"filename"`
	} `mapstructure:"logger"`

	Clickhouse struct {
		Host     string `mapstructure:"host"`
		Port     string `mapstructure:"port"`
		Name     string `mapstructure:"name"`
		User     string `mapstructure:"user"`
		Password string `mapstructure:"password"`
	} `mapstructure:"clickhouse"`

	Kafka struct {
		Host                     string `mapstructure:"host"`
		Port                     string `mapstructure:"port"`
		Partition                int32  `mapstructure:"partition"`
		AuthTopicSend            string `mapstructure:"auth_topic_send"`
		AuthTopicRecieve         string `mapstructure:"auth_topic_recieve"`
		GameSettingsTopicRecieve string `mapstructure:"gameSettings_topic_recieve"`
		MatchTopicSend           string `mapstructure:"match_topic_send"`
	} `mapstructure:"kafka"`

	Redis struct {
		Host     string `mapstructure:"host"`
		Port     string `mapstructure:"port"`
		Password string `mapstructure:"password"`
		DB       int    `mapstructure:"db"`
	} `mapstructure:"redis"`

	Random struct {
		ApiKey string `mapstructure:"api-key"`
		ApiUrl string `mapstructure:"api-url"`
	} `mapstructure:"random"`

	Jaeger struct {
		Host        string `mapstructure:"host"`
		Port        string `mapstructure:"port"`
		ServiceName string `mapstructure:"service_name"`
	} `mapstructure:"jaeger"`
}

func MustParseConfig() Config {

	var Cfg Config

	viper.AddConfigPath("./internal/config")
	viper.SetConfigType("yaml")
	viper.SetConfigName("config")

	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}

	viper.Unmarshal(&Cfg)

	return Cfg

}
