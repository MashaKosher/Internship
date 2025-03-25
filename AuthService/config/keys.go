package config

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
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

func ReadKeys() {
	publicPath, privatePath := composePathes()

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
	publicPath := filepath.Join(currentDir, Envs.PublicKeyFile)
	privatePath := filepath.Join(currentDir, Envs.PrivateKeyFile)
	return publicPath, privatePath
}

func readPath(keyType, filePath string, wg *sync.WaitGroup) {
	defer wg.Done()
	keyData, err := os.ReadFile(filePath)
	if err != nil {
		Logger.Error("Error while reading " + keyType + " file")
		panic("Error while reading " + keyType + " file")
	}

	block, _ := pem.Decode(keyData)
	if block == nil && (block.Type != PUBLIC_KEY || block.Type != PRIVATE_KEY) {
		Logger.Error("Invalid Key format")
		panic("Invalid Key format")
	}

	if keyType == PUBLIC_KEY {
		publicKey, err := x509.ParsePKIXPublicKey(block.Bytes)
		if err != nil {
			Logger.Error("Error while parsing " + keyType + " file")
			panic("Error while parsing " + keyType + " file")
		}
		RSAkeys.PublicKey = publicKey.(*rsa.PublicKey)
	} else {
		privateKey, err := x509.ParsePKCS8PrivateKey(block.Bytes)
		if err != nil {
			Logger.Error("Error while parsing " + keyType + " file")
			panic("Error while parsing " + keyType + " file")
		}
		RSAkeys.PrivateKey = privateKey.(*rsa.PrivateKey)
	}

	Logger.Info(keyType + " readed succesfully")
}

// func readPublic(filePath string, wg *sync.WaitGroup) {
// 	defer wg.Done()
// 	keyData, err := os.ReadFile(filePath)
// 	if err != nil {
// 		fmt.Println("Ошибка чтения публичного ключа:", err)
// 		return
// 	}

// 	block, _ := pem.Decode(keyData)
// 	if block == nil || block.Type != "PUBLIC KEY" {
// 		fmt.Println("неверный формат публичного ключа")
// 		return
// 	}

// 	publicKey, err := x509.ParsePKIXPublicKey(block.Bytes)
// 	if err != nil {
// 		fmt.Println("Ошибка парсинга публичного ключа:", err)
// 		return
// 	}

// 	RSAkeys.PublicKey = publicKey.(*rsa.PublicKey)
// 	fmt.Println("Публичный ключ загружен:", RSAkeys.PublicKey)
// }

// func readPrivate(filePath string, wg *sync.WaitGroup) {
// 	defer wg.Done()
// 	keyData, err := os.ReadFile(filePath)
// 	if err != nil {
// 		fmt.Println("Ошибка чтения приватного ключа:", err)
// 		return
// 	}

// 	block, _ := pem.Decode(keyData)
// 	if block == nil || block.Type != "PRIVATE KEY" {
// 		fmt.Println("неверный формат приватного ключа")
// 		return
// 	}

// 	privateKey, err := x509.ParsePKCS8PrivateKey(block.Bytes)
// 	if err != nil {
// 		fmt.Println("Ошибка парсинга приватного ключа:", err)
// 		return
// 	}

// 	RSAkeys.PrivateKey = privateKey.(*rsa.PrivateKey)
// 	fmt.Println("Приватный ключ загружен:", RSAkeys.PrivateKey)
// }
