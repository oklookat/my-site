package cryptor

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"github.com/pkg/errors"
	"io"
)

// AESEncrypt - get text and return encrypted text (AES).
func AESEncrypt(text string, secret string) (encrypted string, err error) {
	key := []byte(secret)
	plaintext := []byte(text)
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", errors.Wrap(err, "AESEncrypt: make cipher failed. Error:")
	}
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", errors.Wrap(err, "AESEncrypt: make GCM failed. Error:")
	}
	nonce := make([]byte, aesGCM.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return "", errors.Wrap(err, "AESEncrypt: read nonce failed. Error:")
	}
	ciphertext := aesGCM.Seal(nonce, nonce, plaintext, nil)
	encrypted = fmt.Sprintf("%x", ciphertext)
	return encrypted, err
}

// AESDecrypt - get encrypted and return decrypted text (AES).
func AESDecrypt(encrypted string, secret string) (text string, err error) {
	if len(encrypted) < 4 {
		return "", errors.New("AESDecrypt: encrypted must be min 4 bit length")
	}
	if len(secret) < 32 {
		return "", errors.New("AESDecrypt: secret must be 32 bit length")
	}
	key := []byte(secret)
	// get hex from string (this is encrypted data)
	hexed, err := hex.DecodeString(encrypted)
	if err != nil {
		// if presented string not a hex
		return "", errors.Wrap(err, "AESDecrypt: failed to decode hex. Error:")
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", errors.Wrap(err, "AESDecrypt: failed to make cipher. Error:")
	}
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", errors.Wrap(err, "AESDecrypt: failed to make GCM. Error:")
	}
	nonceSize := aesGCM.NonceSize()
	nonce, ciphertext := hexed[:nonceSize], hexed[nonceSize:]
	plaintext, err := aesGCM.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", errors.Wrap(err, "AESDecrypt: failed to open GCM. Error:")
	}
	var decrypted = fmt.Sprintf("%s", plaintext)
	return decrypted, err
}
