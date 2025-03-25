package config

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

// func main() {
// 	ReadKeys()
// 	fmt.Println(RSAkeys.PublicKey)
// 	fmt.Println(RSAkeys.PrivateKey)
// }

type Keys struct {
	PublicKey  *rsa.PublicKey
	PrivateKey *rsa.PrivateKey
}

var RSAkeys Keys

func ReadKeys() {
	publicPath, privatePath := composePathes()

	var wg sync.WaitGroup
	wg.Add(2)
	go readPublic(publicPath, &wg)
	go readPrivate(privatePath, &wg)
	wg.Wait()
}

func composePathes() (string, string) {
	currentDir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	publicPath := filepath.Join(currentDir, "config", "internal", Envs.PublicKeyFile)
	privatePath := filepath.Join(currentDir, "config", "internal", Envs.PrivateKeyFile)

	return publicPath, privatePath
}

func readPublic(filePath string, wg *sync.WaitGroup) {
	defer wg.Done()
	keyData, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Println("Ошибка чтения публичного ключа:", err)
		return
	}

	block, _ := pem.Decode(keyData)
	if block == nil || block.Type != "PUBLIC KEY" {
		fmt.Println("неверный формат публичного ключа")
		return
	}

	publicKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		fmt.Println("Ошибка парсинга публичного ключа:", err)
		return
	}

	RSAkeys.PublicKey = publicKey.(*rsa.PublicKey)
	fmt.Println("Публичный ключ загружен:", RSAkeys.PublicKey)
}

func readPrivate(filePath string, wg *sync.WaitGroup) {
	defer wg.Done()
	keyData, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Println("Ошибка чтения приватного ключа:", err)
		return
	}

	block, _ := pem.Decode(keyData)
	if block == nil || block.Type != "PRIVATE KEY" {
		fmt.Println("неверный формат приватного ключа")
		return
	}

	privateKey, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		fmt.Println("Ошибка парсинга приватного ключа:", err)
		return
	}

	RSAkeys.PrivateKey = privateKey.(*rsa.PrivateKey)
	fmt.Println("Приватный ключ загружен:", RSAkeys.PrivateKey)
}
