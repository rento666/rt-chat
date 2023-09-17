package tencrypt

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// GenerateAndSaveKeyPair 生成RSA公钥和私钥，并保存到文件
func GenerateAndSaveKeyPair() error {
	privateKeyPath := getPath("private_key.pem")
	publicKeyPath := getPath("public_key.pem")
	// 生成RSA密钥对
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return fmt.Errorf("failed to generate RSA key pair: %w", err)
	}

	// 将私钥保存到文件
	privateKeyFile, err := os.Create(privateKeyPath)
	if err != nil {
		return fmt.Errorf("failed to create private key file: %w", err)
	}
	defer privateKeyFile.Close()

	privateKeyPEM := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	}
	err = pem.Encode(privateKeyFile, privateKeyPEM)
	if err != nil {
		return fmt.Errorf("failed to save private key to file: %w", err)
	}

	// 生成并保存公钥
	publicKey := privateKey.PublicKey
	publicKeyFile, err := os.Create(publicKeyPath)
	if err != nil {
		return fmt.Errorf("failed to create public key file: %w", err)
	}
	defer publicKeyFile.Close()

	publicKeyPEM := &pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: x509.MarshalPKCS1PublicKey(&publicKey),
	}
	err = pem.Encode(publicKeyFile, publicKeyPEM)
	if err != nil {
		return fmt.Errorf("failed to save public key to file: %w", err)
	}

	return nil
}

// RsaEncrypt 使用公钥加密数据
func RsaEncrypt(plaintext []byte) ([]byte, error) {

	parsedPublicKey, err := GetPublicKey()
	if err != nil {
		return nil, fmt.Errorf("failed to parse public key: %w", err)
	}

	ciphertext, err := rsa.EncryptPKCS1v15(rand.Reader, parsedPublicKey, plaintext)
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt data: %w", err)
	}

	return ciphertext, nil
}

// RsaDecrypt 使用私钥解密数据
func RsaDecrypt(ciphertext []byte) ([]byte, error) {

	ciphertext, er := base64.StdEncoding.DecodeString(string(ciphertext))
	if er != nil {
		return nil, fmt.Errorf("failed to decode base64: %w", er)
	}
	parsedPrivateKey, err := GetPrivateKey()
	if err != nil {
		return nil, fmt.Errorf("failed to parse private key: %w", err)
	}
	plaintext, err := rsa.DecryptPKCS1v15(rand.Reader, parsedPrivateKey, ciphertext)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt data: %w", err)
	}

	return plaintext, nil
}

// GetPublicKey 获取RSA公钥
func GetPublicKey() (*rsa.PublicKey, error) {
	publicKeyPath := getPath("public_key.pem")
	publicKeyFile, err := os.Open(publicKeyPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open public key file: %w", err)
	}
	defer publicKeyFile.Close()

	publicKeyBytes, err := io.ReadAll(publicKeyFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read public key file: %w", err)
	}

	block, _ := pem.Decode(publicKeyBytes)
	if block == nil {
		return nil, fmt.Errorf("failed to decode public key")
	}

	parsedPublicKey, err := x509.ParsePKCS1PublicKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse public key: %w", err)
	}
	return parsedPublicKey, nil
}

// GetPrivateKey 获取RSA私钥
func GetPrivateKey() (*rsa.PrivateKey, error) {

	privateKeyPath := getPath("private_key.pem")
	privateKeyFile, err := os.Open(privateKeyPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open private key file: %w", err)
	}
	defer privateKeyFile.Close()

	privateKeyBytes, err := io.ReadAll(privateKeyFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read private key file: %w", err)
	}

	block, _ := pem.Decode(privateKeyBytes)
	if block == nil {
		return nil, fmt.Errorf("failed to decode private key")
	}

	parsedPrivateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse private key: %w", err)
	}
	return parsedPrivateKey, nil
}

func getPath(name string) string {
	currentDir, err := os.Getwd()
	if err != nil {
		fmt.Println("Failed to get current working directory:", err)
		return ""
	}

	return filepath.Join(currentDir, "tkey", name)
}
