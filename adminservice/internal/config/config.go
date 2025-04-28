package config

import (
	"github.com/spf13/viper"
)

type config struct {
	Server struct {
		Port string `mapstructure:"port"`
		Host string `mapstructure:"host"`
	} `mapstructure:"server"`

	DB struct {
		Host     string `mapstructure:"host"`
		Port     string `mapstructure:"port"`
		Name     string `mapstructure:"name"`
		User     string `mapstructure:"user"`
		Password string `mapstructure:"password"`
		SSLMode  string `mapstructure:"ssl_mode"`
	} `mapstructure:"db"`

	Logger struct {
		FileName string `mapstructure:"filename"`
	} `mapstructure:"logger"`

	Kafka struct {
		Host                  string `mapstructure:"host"`
		Port                  string `mapstructure:"port"`
		Partition             int32  `mapstructure:"partition"`
		AuthTopicSend         string `mapstructure:"auth_topic_send"`
		AuthTopicRecieve      string `mapstructure:"auth_topic_recieve"`
		SeasonTopicSend       string `mapstructure:"season_topic_send"`
		DailyTaskTopicSend    string `mapstructure:"dailyTasks_topic_send"`
		GameSettingsTopicSend string `mapstructure:"gameSettings_topic_send"`
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
