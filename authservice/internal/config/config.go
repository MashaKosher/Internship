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

	RSAKeys struct {
		PublicKeyFile  string `mapstructure:"public"`
		PrivateKeyFile string `mapstructure:"private"`
	} `mapstructure:"rsa_keys"`

	Kafka struct {
		Host           string `mapstructure:"host"`
		Port           string `mapstructure:"port"`
		TopicSend      string `mapstructure:"auth_topic_send"`
		TopicRecieve   string `mapstructure:"auth_topic_recieve"`
		UserSignupSend string `mapstructure:"user_signup_send"`
	} `mapstructure:"kafka"`

	Jaeger struct {
		Host        string `mapstructure:"host"`
		Port        string `mapstructure:"port"`
		ServiceName string `mapstructure:"service_name"`
	} `mapstructure:"jaeger"`

	Memcached struct {
		Host string `mapstructure:"host"`
		Port string `mapstructure:"port"`
	} `mapstructure:"memcached"`
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
