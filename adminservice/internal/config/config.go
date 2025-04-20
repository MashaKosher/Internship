package config

import (
	"log"

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
		Host             string `mapstructure:"host"`
		Port             string `mapstructure:"port"`
		Partition        int32  `mapstructure:"partition"`
		AuthTopicSend    string `mapstructure:"auth_topic_send"`
		AuthTopicRecieve string `mapstructure:"auth_topic_recieve"`
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
	log.Println(AppConfig.DB.Host)
	log.Println(AppConfig.DB.Host)

	// viper.SetConfigName("config")
	// viper.SetConfigType("yaml")
	// viper.AddConfigPath("./internal/config")

	// err := viper.ReadInConfig() // Find and read the config file
	// if err != nil {             // Handle errors reading the config file
	// 	panic(fmt.Errorf("fatal error config file: %w", err))
	// }

	// Cfg.Server.Port = viper.GetString("server.port")
	// Cfg.Server.Host = viper.GetString("server.host")

	// Cfg.DB.User = viper.GetString("db.user")
	// Cfg.DB.Password = viper.GetString("db.password")
	// Cfg.DB.Name = viper.GetString("db.db")
	// Cfg.DB.Host = viper.GetString("db.host")
	// Cfg.DB.Port = viper.GetString("db.port")
	// Cfg.DB.SSLMode = viper.GetString("ssl_mode")

}
