package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"github.com/spf13/viper"
	rand1 "math/rand"
	"strings"
)

func rsaEncrypt(data []byte, publicKey *rsa.PublicKey) ([]byte, error) {
	return rsa.EncryptOAEP(
		sha256.New(), // 使用SHA-256哈希算法
		rand.Reader,
		publicKey,
		data,
		nil,
	)
}

func aesEncrypt(data, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	// GCM模式需要一个12字节的nonce
	nonce := make([]byte, 12)
	if _, err := rand.Read(nonce); err != nil {
		return nil, err
	}
	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	return aesgcm.Seal(nonce, nonce, data, nil), nil
}

func EncryptPassword(password string) (string, error) {
	if password == "" {
		return "", nil
	}
	// 解码Base64编码的公钥
	rsaPublicKeyText := viper.GetString("jumpserver.jms_public_key")
	rsaPublicKeyText = strings.ReplaceAll(rsaPublicKeyText, "\"", "")
	publicKeyBytes, err := base64.StdEncoding.DecodeString(rsaPublicKeyText)
	if err != nil {
		return "", err
	}

	publicKey, err := parseRSAPublicKey(publicKeyBytes)
	if err != nil {
		return "", err
	}

	// 生成AES密钥
	aesKey := generateRandomKey(16)

	// RSA加密AES密钥
	keyCipher, err := rsaEncrypt([]byte(aesKey), publicKey)
	if err != nil {
		return "", err
	}

	// AES加密密码
	passwordCipher, err := aesEncrypt([]byte(password), []byte(aesKey))
	if err != nil {
		return "", err
	}

	// 将结果编码为Base64
	keyCipherBase64 := base64.StdEncoding.EncodeToString(keyCipher)
	passwordCipherBase64 := base64.StdEncoding.EncodeToString(passwordCipher)

	return fmt.Sprintf("%s:%s", keyCipherBase64, passwordCipherBase64), nil
}

func generateRandomKey(length int) string {
	const charset = "0123456789abcdefghijklmnopqrstuvwxyz"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand1.Intn(len(charset))]
	}
	return string(b)
}

func parseRSAPublicKey(pemBytes []byte) (*rsa.PublicKey, error) {
	// 解码PEM格式
	block, _ := pem.Decode(pemBytes)
	if block == nil || block.Type != "PUBLIC KEY" {
		return nil, fmt.Errorf("failed to decode PEM block containing public key")
	}

	// 解析DER格式的公钥
	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	// 断言公钥类型
	rsaPub, ok := pub.(*rsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("expected RSA public key, got %T", pub)
	}

	return rsaPub, nil
}
