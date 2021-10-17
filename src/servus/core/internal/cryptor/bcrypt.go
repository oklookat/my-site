package cryptor

import "golang.org/x/crypto/bcrypt"

type Bcrypt struct {
	// Cost - see bcrypt.GenerateFromPassword.
	Cost int
}

// Hash - hash (bcrypt).
func (b *Bcrypt) Hash(data string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(data), b.Cost)
	return string(bytes), err
}

// Check - check hash (bcrypt).
func (b *Bcrypt) Check(data string, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(data))
	return err == nil
}
