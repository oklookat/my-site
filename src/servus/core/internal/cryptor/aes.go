package cryptor

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"

	"github.com/pkg/errors"
)

var (
	ErrAESMakeCipher          = errors.New("cryptor (AES): make cipher failed")
	ErrAESMakeGCM             = errors.New("cryptor (AES): make GCM failed")
	ErrAESOpenGCM             = errors.New("cryptor (AES): open GCM failed")
	ErrAESDecodeHEX           = errors.New("cryptor (AES): failed to decode hex")
	ErrAESReadNonce           = errors.New("cryptor (AES): read nonce failed")
	ErrAESEncryptedValidation = errors.New("cryptor (AES): encrypted must be min 4 bit length")
	ErrAESSecretValidation    = errors.New("cryptor (AES): secret must be 32 bit length")
)

type AES struct {
	Secret string
}

// get text and return encrypted text (AES).
func (a *AES) Encrypt(text string) (encrypted string, err error) {
	key := []byte(a.Secret)
	plaintext := []byte(text)
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", ErrAESMakeCipher
	}
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", ErrAESMakeGCM
	}
	nonce := make([]byte, aesGCM.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return "", ErrAESReadNonce
	}
	ciphertext := aesGCM.Seal(nonce, nonce, plaintext, nil)
	encrypted = fmt.Sprintf("%x", ciphertext)
	return encrypted, err
}

// get encrypted and return decrypted text (AES).
func (a *AES) Decrypt(encrypted string) (text string, err error) {
	if len(encrypted) < 4 {
		return "", ErrAESEncryptedValidation
	}
	if len(a.Secret) < 32 {
		return "", ErrAESSecretValidation
	}
	key := []byte(a.Secret)
	// get hex from string (this is encrypted data)
	hexed, err := hex.DecodeString(encrypted)
	if err != nil {
		// if presented string not a hex
		return "", ErrAESDecodeHEX
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", ErrAESMakeCipher
	}
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", ErrAESMakeGCM
	}
	nonceSize := aesGCM.NonceSize()
	nonce, ciphertext := hexed[:nonceSize], hexed[nonceSize:]
	plaintext, err := aesGCM.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", ErrAESOpenGCM
	}
	var decrypted = fmt.Sprintf("%s", plaintext)
	return decrypted, err
}
