// Package encryption contains funtions to encrypt and decrypt text used by the application
package encryption

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
)

// CreateHash creates a MD5 sum from the encryption key
func CreateHash(key string) string {
	hasher := md5.New()
	hasher.Write([]byte(key))
	return hex.EncodeToString(hasher.Sum(nil))
}

// EncryptData encrypts a slice of bytes using a passphrase
func EncryptData(data string, passphrase string) string {
	block, _ := aes.NewCipher([]byte(CreateHash(passphrase)))
	ciphertext := make([]byte, aes.BlockSize+len(data))
	iv := ciphertext[:aes.BlockSize]
	cfb := cipher.NewCFBEncrypter(block, iv)
	cfb.XORKeyStream(ciphertext[aes.BlockSize:], []byte(data))

	// Encode the ciphertext to Base64
	encodedCiphertext := base64.StdEncoding.EncodeToString(ciphertext)
	return encodedCiphertext
}

// DecryptData decrypts a slice of bytes using a passphrase
func DecryptData(encodedCiphertext string, passphrase string) string {
	key := []byte(CreateHash(passphrase))
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err.Error())
	}

	// Decode the Base64 encoded ciphertext
	ciphertext, err := base64.StdEncoding.DecodeString(encodedCiphertext)
	if err != nil {
		panic(err.Error())
	}

	plaintext := make([]byte, len(ciphertext)-aes.BlockSize)
	iv := ciphertext[:aes.BlockSize]
	cfb := cipher.NewCFBDecrypter(block, iv)
	cfb.XORKeyStream(plaintext, ciphertext[aes.BlockSize:])
	return string(plaintext)
}
