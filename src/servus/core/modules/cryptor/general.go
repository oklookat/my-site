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

// BHash hash (bcrypt)
func BHash(data string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(data), 14)
	return string(bytes), err
}

// BHashCheck check hash (bcrypt)
func BHashCheck(data string, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(data))
	return err == nil
}

// AESEncrypt get text and return encrypted text (AES)
func AESEncrypt(text string, secret string) (encrypted string, error AESError) {
	key := []byte(secret)
	plaintext := []byte(text)
	block, err := aes.NewCipher(key)
	if err != nil {
		var err = AESError{HaveErrors: true, ErrorCode: "PIPE_MAKE_BLOCK_FAILED", Error: err}
		return "", err
	}
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		var err = AESError{HaveErrors: true, ErrorCode: "PIPE_MAKE_GCM_FAILED", Error: err}
		return "", err
	}
	nonce := make([]byte, aesGCM.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		var err = AESError{HaveErrors: true, ErrorCode: "PIPE_MAKE_NONCE_FAILED", Error: err}
		return "", err
	}
	ciphertext := aesGCM.Seal(nonce, nonce, plaintext, nil)
	encrypted = fmt.Sprintf("%x", ciphertext)
	return encrypted, AESError{HaveErrors: false}
}

// AESDecrypt get encrypted and return decrypted text (AES)
func AESDecrypt(encrypted string, secret string) (text string, error AESError) {
	key := []byte(secret)
	enc, err := hex.DecodeString(encrypted)
	if err != nil {
		// if not hex in string
		var err = AESError{HaveErrors: true, ErrorCode: "PIPE_DECODE_HEX_STRING_FAILED", Error: err}
		return "", err
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		var err = AESError{HaveErrors: true, ErrorCode: "PIPE_MAKE_BLOCK_FAILED", Error: err}
		return "", err
	}
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		var err = AESError{HaveErrors: true, ErrorCode: "PIPE_MAKE_GCM_FAILED", Error: err}
		return "", err
	}
	nonceSize := aesGCM.NonceSize()
	nonce, ciphertext := enc[:nonceSize], enc[nonceSize:]
	plaintext, err := aesGCM.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		var err = AESError{HaveErrors: true, ErrorCode: "PIPE_DECRYPT_FAILED", Error: err}
		return "", err
	}
	var decrypted = fmt.Sprintf("%s", plaintext)
	return decrypted, AESError{HaveErrors: false}
}
