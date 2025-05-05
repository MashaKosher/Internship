package config

import "github.com/spf13/viper"

type Config struct {
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
		SeasonTopicRecieve    string `mapstructure:"season_topic_recieve"`
		DailyTaskTopicRecieve string `mapstructure:"daily_task_topic_recieve"`
		MatchTopicRecieve     string `mapstructure:"match_topic_recieve"`
		UserSignupRecieve     string `mapstructure:"user_signup_recieve"`
	} `mapstructure:"kafka"`

	Redis struct {
		Host     string `mapstructure:"host"`
		Port     string `mapstructure:"port"`
		Password string `mapstructure:"password"`
		DB       int    `mapstructure:"db"`
	} `mapstructure:"redis"`
}

var AppConfig Config

func MustParseConfig() Config {
	viper.AddConfigPath("./internal/config")
	viper.SetConfigType("yaml")
	viper.SetConfigName("config")

	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}

	viper.Unmarshal(&AppConfig)

	return AppConfig
}
