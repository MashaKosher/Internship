package config

import "github.com/spf13/viper"

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

	RSAKeys struct {
		PublicKeyFile  string `mapstructure:"public"`
		PrivateKeyFile string `mapstructure:"private"`
	} `mapstructure:"rsa_keys"`

	Kafka struct {
		Host         string `mapstructure:"host"`
		Port         string `mapstructure:"port"`
		TopicSend    string `mapstructure:"auth_topic_send"`
		TopicRecieve string `mapstructure:"auth_topic_recieve"`
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
