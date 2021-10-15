package cryptor

import "golang.org/x/crypto/bcrypt"

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
