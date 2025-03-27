package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Server
	DB
	Logger
	RSAKeys
}

type Server struct {
	Port int
}

type DB struct {
	Host     string
	Port     string
	Name     string
	User     string
	Password string
	SSLMode  string
}

type Logger struct {
	LogFileName string
}

type RSAKeys struct {
	PublicKeyFile  string
	PrivateKeyFile string
}

var Cfg Config

func LoadEnvs() {
	godotenv.Load()

	// DB set up
	Cfg.DB.Host = os.Getenv("POSTGRES_HOST")
	Cfg.DB.Port = os.Getenv("POSTGRES_PORT")
	Cfg.DB.Name = os.Getenv("POSTGRES_DB")
	Cfg.DB.User = os.Getenv("POSTGRES_USER")
	Cfg.DB.Password = os.Getenv("POSTGRES_PASSWORD")
	Cfg.DB.SSLMode = os.Getenv("SSL_MODE")

	// Logger set up
	Cfg.Logger.LogFileName = os.Getenv("LOG_FILE_NAME")

	// RSA keys
	Cfg.RSAKeys.PublicKeyFile = os.Getenv("PUBLIC_KEY_FILE_NAME")
	Cfg.RSAKeys.PrivateKeyFile = os.Getenv("PRIVATE_KEY_FILE_NAME")

}
