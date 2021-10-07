package cryptor

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"io"
)

// BHash - hash (bcrypt).
func BHash(data string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(data), 14)
	return string(bytes), err
}

// BHashCheck - check hash (bcrypt).
func BHashCheck(data string, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(data))
	return err == nil
}

// AESEncrypt - get text and return encrypted text (AES).
func AESEncrypt(text string, secret string) (encrypted string, error AESError) {
	key := []byte(secret)
	plaintext := []byte(text)
	block, err := aes.NewCipher(key)
	if err != nil {
		var err = AESError{HasErrors: true, AdditionalErr: ErrMakeBlock, OriginalErr: err}
		return "", err
	}
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		var err = AESError{HasErrors: true, AdditionalErr: ErrMakeGCM, OriginalErr: err}
		return "", err
	}
	nonce := make([]byte, aesGCM.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		var err = AESError{HasErrors: true, AdditionalErr: ErrMakeNonceFromGCM, OriginalErr: err}
		return "", err
	}
	ciphertext := aesGCM.Seal(nonce, nonce, plaintext, nil)
	encrypted = fmt.Sprintf("%x", ciphertext)
	return encrypted, AESError{HasErrors: false}
}

// AESDecrypt - get encrypted and return decrypted text (AES).
func AESDecrypt(encrypted string, secret string) (text string, error AESError) {
	key := []byte(secret)
	// get hex from string (this is encrypted data)
	hexed, err := hex.DecodeString(encrypted)
	if err != nil {
		// if presented string not a hex
		var err = AESError{HasErrors: true, AdditionalErr: ErrDecodeHEX, OriginalErr: err}
		return "", err
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		var err = AESError{HasErrors: true, AdditionalErr: ErrMakeBlock, OriginalErr: err}
		return "", err
	}
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		var err = AESError{HasErrors: true, AdditionalErr: ErrMakeGCM, OriginalErr: err}
		return "", err
	}
	nonceSize := aesGCM.NonceSize()
	nonce, ciphertext := hexed[:nonceSize], hexed[nonceSize:]
	plaintext, err := aesGCM.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		var err = AESError{HasErrors: true, AdditionalErr: ErrDecryption, OriginalErr: err}
		return "", err
	}
	var decrypted = fmt.Sprintf("%s", plaintext)
	return decrypted, AESError{HasErrors: false}
}
