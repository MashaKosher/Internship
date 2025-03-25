package config

import (
	"os"

	"github.com/joho/godotenv"
)

type envVariables struct {
	DBHost         string
	DBPort         string
	DBName         string
	DBUser         string
	DBPassword     string
	DBSSLMode      string
	LogFileName    string
	PublicKeyFile  string
	PrivateKeyFile string
}

var Envs envVariables

func LoadEnvs() {
	godotenv.Load()

	// DB set up
	Envs.DBHost = os.Getenv("POSTGRES_HOST")
	Envs.DBPort = os.Getenv("POSTGRES_PORT")
	Envs.DBName = os.Getenv("POSTGRES_DB")
	Envs.DBUser = os.Getenv("POSTGRES_USER")
	Envs.DBPassword = os.Getenv("POSTGRES_PASSWORD")
	Envs.DBSSLMode = os.Getenv("SSL_MODE")

	// Logger set up
	Envs.LogFileName = os.Getenv("LOG_FILE_NAME")

	// RSA keys
	Envs.PublicKeyFile = os.Getenv("PUBLIC_KEY_FILE_NAME")
	Envs.PrivateKeyFile = os.Getenv("PRIVATE_KEY_FILE_NAME")

}
