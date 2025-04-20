// package config

// import (
// 	"os"

// 	"github.com/joho/godotenv"
// )

// type Config struct {
// 	Server
// 	DB
// 	Logger
// 	RSAKeys
// }

// type Server struct {
// 	Port string
// }

// type DB struct {
// 	Host     string
// 	Port     string
// 	Name     string
// 	User     string
// 	Password string
// 	SSLMode  string
// }

// type Logger struct {
// 	LogFileName string
// }

// type RSAKeys struct {
// 	PublicKeyFile  string
// 	PrivateKeyFile string
// }

// var Cfg Config

// func LoadEnvs() {
// 	godotenv.Load()

// 	// DB set up
// 	Cfg.DB.Host = os.Getenv("POSTGRES_HOST")
// 	Cfg.DB.Port = os.Getenv("POSTGRES_PORT")
// 	Cfg.DB.Name = os.Getenv("POSTGRES_DB")
// 	Cfg.DB.User = os.Getenv("POSTGRES_USER")
// 	Cfg.DB.Password = os.Getenv("POSTGRES_PASSWORD")
// 	Cfg.DB.SSLMode = os.Getenv("SSL_MODE")

// 	// Logger set up
// 	Cfg.Logger.LogFileName = os.Getenv("LOG_FILE_NAME")

// 	// RSA keys
// 	Cfg.RSAKeys.PublicKeyFile = os.Getenv("PUBLIC_KEY_FILE_NAME")
// 	Cfg.RSAKeys.PrivateKeyFile = os.Getenv("PRIVATE_KEY_FILE_NAME")

// }

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
