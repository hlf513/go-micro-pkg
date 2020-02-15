package aes

import (
	"crypto/aes"
	"encoding/base64"
)

// ECBPKCS5EncryptToBase64 
func ECBPKCS5EncryptToBase64(src, key string) (string, error) {
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", err
	}
	ecb := NewECBEncrypter(block)
	content := []byte(src)
	content = pkcs5Padding(content, block.BlockSize())
	crypted := make([]byte, len(content))
	ecb.CryptBlocks(crypted, content)
	return base64.StdEncoding.EncodeToString(crypted), nil
}
