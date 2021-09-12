package core

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
	"io"
	"os"
)

type utils interface {
	GetExecuteDir() string
	HashPassword(password string) (string, error)
	HashPasswordCheck(password, hash string) bool
	EncryptAES(text string) (hashed string, encrypted string, err error)
	DecryptAES(hashed string, encrypted string) (string, error)
}

type Utils struct {
	utils
}

func (u Utils) GetExecuteDir() string {
	path, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	return path
}

func (u Utils) HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func (u Utils) HashPasswordCheck(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// EncryptAES get text and return AES-hashed and non-hashed version of text
func (u Utils) EncryptAES(text string) (hashed string, encrypted string, err error) {
	key := []byte(servus.Config.Secret)
	plaintext := []byte(text)
	block, err := aes.NewCipher(key)
	if err != nil {
		servus.Logger.Error(err.Error())
		return "", "", errors.New("PIPE_MAKE_BLOCK_FAILED")
	}
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		servus.Logger.Error(err.Error())
		return "", "", errors.New("PIPE_MAKE_GCM_FAILED")
	}
	nonce := make([]byte, aesGCM.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		servus.Logger.Error(err.Error())
		return "", "", errors.New("PIPE_MAKE_NONCE_FAILED")
	}
	ciphertext := aesGCM.Seal(nonce, nonce, plaintext, nil)
	encrypted = fmt.Sprintf("%x", ciphertext)
	hashed, _ = u.HashPassword(encrypted)
	return hashed, encrypted, nil
}

// DecryptAES get AES hashed and non-hashed data and return encrypted data
func (u Utils) DecryptAES(hashed string, encrypted string) (string, error) {
	var isEqual = u.HashPasswordCheck(hashed, encrypted)
	if !isEqual {
		return "", errors.New("PIPE_DECRYPT_HASH_NOT_SAME")
	}
	key := []byte(servus.Config.Secret)
	enc, _ := hex.DecodeString(encrypted)
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err.Error())
	}
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		panic(err.Error())
	}
	nonceSize := aesGCM.NonceSize()
	nonce, ciphertext := enc[:nonceSize], enc[nonceSize:]
	plaintext, err := aesGCM.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		panic(err.Error())
	}
	var decrypted = fmt.Sprintf("%s", plaintext)
	return decrypted, nil
}
