package keys

import (
	"authservice/internal/di"
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

func ReadRSAKeys(cfg di.ConfigType, logger di.LoggerType, RSAKeys *di.RSAKeys) {
	publicPath, privatePath := composePathes(cfg)

	logger.Info("publicPath: " + fmt.Sprint(publicPath))
	logger.Info("privatePath: " + fmt.Sprint(privatePath))

	var wg sync.WaitGroup
	wg.Add(2)
	go readPath(PUBLIC_KEY, publicPath, &wg, logger, RSAKeys)
	go readPath(PRIVATE_KEY, privatePath, &wg, logger, RSAKeys)
	wg.Wait()
}

func composePathes(cfg di.ConfigType) (string, string) {
	currentDir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	publicPath := filepath.Join(currentDir, cfg.RSAKeys.PublicKeyFile)
	privatePath := filepath.Join(currentDir, cfg.RSAKeys.PrivateKeyFile)
	return publicPath, privatePath
}

func readPath(keyType, filePath string, wg *sync.WaitGroup, logger di.LoggerType, RSAKeys *di.RSAKeys) {
	defer wg.Done()
	keyData, err := os.ReadFile(filePath)
	if err != nil {
		logger.Fatal("Error while reading " + keyType + " file")
	}

	block, _ := pem.Decode(keyData)
	if block == nil || (block.Type != PUBLIC_KEY && block.Type != PRIVATE_KEY) {
		logger.Fatal("Invalid Key format")
	}

	if keyType == PUBLIC_KEY {
		publicKey, err := x509.ParsePKIXPublicKey(block.Bytes)
		if err != nil {
			logger.Fatal("Error while parsing " + keyType + " file")
		}
		RSAKeys.PublicKey = publicKey.(*rsa.PublicKey)
	} else {
		privateKey, err := x509.ParsePKCS8PrivateKey(block.Bytes)
		if err != nil {
			logger.Fatal("Error while parsing " + keyType + " file")
		}
		RSAKeys.PrivateKey = privateKey.(*rsa.PrivateKey)
	}

	logger.Info(keyType + " readed succesfully")
}
