package cryptor

// https://github.com/alexedwards/argon2id

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"golang.org/x/crypto/argon2"
	"strings"
)

var (
	// ErrArgonInvalidHash in returned by ComparePasswordAndHash if the provided
	// hash isn't in the expected format.
	ErrArgonInvalidHash = errors.New("cryptor (argon): hash is not in the correct format")

	// ErrArgonIncompatibleVersion in returned by ComparePasswordAndHash if the
	// provided hash was created using a different version of Argon2.
	ErrArgonIncompatibleVersion = errors.New("cryptor (argon): incompatible version of argon2")
)

// Argon - describes the input parameters used by the Argon2id algorithm.
type Argon struct {
	// Memory - the amount of memory used by the algorithm (in kibibytes).
	Memory uint32

	// Iterations - the number of iterations over the memory.
	Iterations uint32

	// Parallelism - the number of threads (or lanes) used by the algorithm.
	// Recommended value is between 1 and runtime.NumCPU().
	Parallelism uint8

	// SaltLength - length of the random salt. 16 bytes is recommended for password hashing.
	SaltLength uint32

	// KeyLength - length of the generated key. 16 bytes or more is recommended.
	KeyLength uint32
}

// Hash - returns an Argon2id hash of a plain-text password using the
// provided algorithm parameters.
func (a *Argon) Hash(data string) (hash string, err error) {
	salt, err := a.generateRandomBytes(a.SaltLength)
	if err != nil {
		return "", err
	}
	key := argon2.IDKey([]byte(data), salt, a.Iterations, a.Memory, a.Parallelism, a.KeyLength)
	b64Salt := base64.RawStdEncoding.EncodeToString(salt)
	b64Key := base64.RawStdEncoding.EncodeToString(key)
	hash = fmt.Sprintf("$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s", argon2.Version, a.Memory, a.Iterations, a.Parallelism, b64Salt, b64Key)
	return hash, nil
}

// Check - check argon2id hash.
func (a *Argon) Check(data, hash string) (match bool, err error) {
	match, _, err = a.matchHash(data, hash)
	return match, err
}

// matchHash - is like ComparePasswordAndHash, except it also returns the params that the hash was
// created with. This can be useful if you want to update your hash params over time (which you
// should).
func (a *Argon) matchHash(data, hash string) (match bool, config *Argon, err error) {
	params, salt, key, err := a.parseHash(hash)
	if err != nil {
		return false, nil, err
	}
	otherKey := argon2.IDKey([]byte(data), salt, a.Iterations, a.Memory, a.Parallelism, a.KeyLength)
	keyLen := int32(len(key))
	otherKeyLen := int32(len(otherKey))
	if subtle.ConstantTimeEq(keyLen, otherKeyLen) == 0 {
		return false, params, nil
	}
	if subtle.ConstantTimeCompare(key, otherKey) == 1 {
		return true, params, nil
	}
	return false, params, nil
}

// parseHash - expects a hash created from this package, and parses it to return the params used to
// create it, as well as the salt and key (password hash).
func (a *Argon) parseHash(hash string) (config *Argon, salt, key []byte, err error) {
	values := strings.Split(hash, "$")
	if len(values) != 6 {
		return nil, nil, nil, ErrArgonInvalidHash
	}
	var version int
	_, err = fmt.Sscanf(values[2], "v=%d", &version)
	if err != nil {
		return nil, nil, nil, err
	}
	if version != argon2.Version {
		return nil, nil, nil, ErrArgonIncompatibleVersion
	}
	config = &Argon{}
	_, err = fmt.Sscanf(values[3], "m=%d,t=%d,p=%d", &config.Memory, &config.Iterations, &config.Parallelism)
	if err != nil {
		return nil, nil, nil, err
	}
	salt, err = base64.RawStdEncoding.Strict().DecodeString(values[4])
	if err != nil {
		return nil, nil, nil, err
	}
	config.SaltLength = uint32(len(salt))
	key, err = base64.RawStdEncoding.Strict().DecodeString(values[5])
	if err != nil {
		return nil, nil, nil, err
	}
	config.KeyLength = uint32(len(key))

	return config, salt, key, nil
}

// generateRandomBytes - generate cryptographically random bytes.
func (a *Argon) generateRandomBytes(n uint32) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}
	return b, nil
}
