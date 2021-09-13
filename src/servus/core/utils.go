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
	"net/http"
	"os"
	"strings"
	"time"
	"unicode"
)

// RemoveSpaces remove spaces from string
func (u *Utils) RemoveSpaces(str string) string {
	return strings.Map(func(r rune) rune {
		if unicode.IsSpace(r) {
			return -1
		}
		return r
	}, str)
}

// GetExecuteDir get execution directory
func (u *Utils) GetExecuteDir() string {
	path, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	return path
}

// HashPassword hash password (bcrypt)
func (u *Utils) HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

// HashPasswordCheck check hashed password (bcrypt)
func (u *Utils) HashPasswordCheck(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// EncryptAES get text and return encrypted text (AES)
func (u *Utils) EncryptAES(text string) (encrypted string, err error) {
	key := []byte(servus.Config.Secret)
	plaintext := []byte(text)
	block, err := aes.NewCipher(key)
	if err != nil {
		servus.Logger.Error(err.Error())
		return "", errors.New("PIPE_MAKE_BLOCK_FAILED")
	}
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		servus.Logger.Error(err.Error())
		return "", errors.New("PIPE_MAKE_GCM_FAILED")
	}
	nonce := make([]byte, aesGCM.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		servus.Logger.Error(err.Error())
		return "", errors.New("PIPE_MAKE_NONCE_FAILED")
	}
	ciphertext := aesGCM.Seal(nonce, nonce, plaintext, nil)
	encrypted = fmt.Sprintf("%x", ciphertext)
	return encrypted, nil
}

// DecryptAES get encrypted and return decrypted text (AES)
func (u *Utils) DecryptAES(encrypted string) (text string, err error) {
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

// SetCookie set cookie
func (u *Utils) SetCookie(response *http.ResponseWriter, name string, value string) {
	var maxAge, err = u.ConvertTimeWord(servus.Config.Security.Cookie.MaxAge)
	if err != nil {
		servus.Logger.Panic(errors.New("Cookie wrong time. Check your settings."))
	}
	maxAgeSeconds := int(maxAge.Seconds())
	var domain = servus.Config.Security.Cookie.Domain
	var path = servus.Config.Security.Cookie.Path
	var httpOnly = servus.Config.Security.Cookie.HttpOnly
	var secure = servus.Config.Security.Cookie.Secure
	sameSite, err := convertCookieSameSite(servus.Config.Security.Cookie.SameSite)
	if err != nil {
		servus.Logger.Panic(err)
	}
	var cookie = http.Cookie{Name: name, Value: value, Path: path, Domain: domain, MaxAge: maxAgeSeconds, HttpOnly: httpOnly, Secure: secure, SameSite: sameSite}
	http.SetCookie(*response, &cookie)
}

// ConvertTimeWord convert time like "2h"; "2min"; "2sec" to duration
func (u *Utils) ConvertTimeWord(timeShortcut string) (time.Duration, error) {
	timeShortcut = strings.ToLower(timeShortcut)
	timeDuration, err := time.ParseDuration(timeShortcut)
	if err != nil {
		servus.Logger.Panic(errors.New("time converting failed. Is string with time correct?"))
	}
	return timeDuration, nil
}
