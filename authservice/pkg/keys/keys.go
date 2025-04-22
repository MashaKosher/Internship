package keys

import (
	"authservice/internal/config"
	"authservice/pkg/logger"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

const PUBLIC_KEY = "PUBLIC KEY"
const PRIVATE_KEY = "PRIVATE KEY"

type Keys struct {
	PublicKey  *rsa.PublicKey
	PrivateKey *rsa.PrivateKey
}

var RSAkeys Keys

func ReadRSAKeys() {
	publicPath, privatePath := composePathes()

	logger.Logger.Info("publicPath: " + fmt.Sprint(publicPath))
	logger.Logger.Info("privatePath: " + fmt.Sprint(privatePath))

	var wg sync.WaitGroup
	wg.Add(2)
	go readPath(PUBLIC_KEY, publicPath, &wg)
	go readPath(PRIVATE_KEY, privatePath, &wg)
	wg.Wait()
}

func composePathes() (string, string) {
	currentDir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	publicPath := filepath.Join(currentDir, config.AppConfig.RSAKeys.PublicKeyFile)
	privatePath := filepath.Join(currentDir, config.AppConfig.RSAKeys.PrivateKeyFile)
	return publicPath, privatePath
}

func readPath(keyType, filePath string, wg *sync.WaitGroup) {
	defer wg.Done()
	keyData, err := os.ReadFile(filePath)
	if err != nil {
		logger.Logger.Fatal("Error while reading " + keyType + " file")
	}

	block, _ := pem.Decode(keyData)
	if block == nil || (block.Type != PUBLIC_KEY && block.Type != PRIVATE_KEY) {
		logger.Logger.Fatal("Invalid Key format")
	}

	if keyType == PUBLIC_KEY {
		publicKey, err := x509.ParsePKIXPublicKey(block.Bytes)
		if err != nil {
			logger.Logger.Fatal("Error while parsing " + keyType + " file")
		}
		RSAkeys.PublicKey = publicKey.(*rsa.PublicKey)
	} else {
		privateKey, err := x509.ParsePKCS8PrivateKey(block.Bytes)
		if err != nil {
			logger.Logger.Fatal("Error while parsing " + keyType + " file")
		}
		RSAkeys.PrivateKey = privateKey.(*rsa.PrivateKey)
	}

	logger.Logger.Info(keyType + " readed succesfully")
}
