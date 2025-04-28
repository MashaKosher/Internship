package config

import (
	"github.com/spf13/viper"
)

type config struct {
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
	} `mapstructure:"kafka"`
}

var AppConfig config

func Load() {
	viper.AddConfigPath("./internal/config")
	viper.SetConfigType("yaml")
	viper.SetConfigName("config")

	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}

	viper.Unmarshal(&AppConfig)

}
