package helpers

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
)

//create encoding
func NewRSAKey() (string, string, error) {
	key, err := rsa.GenerateKey(rand.Reader, 1024)
	if err != nil {
		return "", "", err
	}
	derStream := x509.MarshalPKCS1PrivateKey(key)
	block := &pem.Block{
		Bytes: derStream,
	}
	str := string(pem.EncodeToMemory(block))
	fmt.Println(str)

	publicKey := &key.PublicKey
	derPkix, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		return "", "", err
	}
	block = &pem.Block{
		Bytes: derPkix,
	}
	str2 := string(pem.EncodeToMemory(block))
	fmt.Println(str2)

	return str, str2, nil
}

func RSA_encrypt(origData []byte, publicKey string) ([]byte, error) {
	block, _ := pem.Decode([]byte(publicKey))
	if block == nil {
		return nil, errors.New("public key error")
	}
	pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	pub := pubInterface.(*rsa.PublicKey)
	return rsa.EncryptPKCS1v15(rand.Reader, pub, origData)
}

func RSA_decrypt(ciphertext []byte, privateKey string) ([]byte, error) {
	block, _ := pem.Decode([]byte(privateKey))
	if block == nil {
		return nil, errors.New("private key error!")
	}
	priv, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	return rsa.DecryptPKCS1v15(rand.Reader, priv, ciphertext)
}
