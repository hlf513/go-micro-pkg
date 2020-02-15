package aes

import (
	"crypto/aes"
	"encoding/base64"
)

// ECBPKCS5DecryptFromBase64
func ECBPKCS5DecryptFromBase64(crypt, key string) (string, error) {
	encrypt, err := base64.StdEncoding.DecodeString(crypt)
	if err != nil {
		return "", err
	}
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", err
	}
	blockMode := NewECBDecrypter(block)
	origData := make([]byte, len(encrypt))
	blockMode.CryptBlocks(origData, encrypt)
	origData = pkcs5UnPadding(origData)
	return string(origData), nil
}
