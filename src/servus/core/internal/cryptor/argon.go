package cryptor

// https://github.com/alexedwards/argon2id
// MIT License

// Copyright (c) 2018 Alex Edwards

// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:

// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"

	"golang.org/x/crypto/argon2"
)

var (
	// returned by ComparePasswordAndHash if the provided
	// hash isn't in the expected format.
	ErrArgonInvalidHash = errors.New("cryptor (argon): hash is not in the correct format")

	// returned by ComparePasswordAndHash if the
	// provided hash was created using a different version of Argon2.
	ErrArgonIncompatibleVersion = errors.New("cryptor (argon): incompatible version of argon2")
)

// describes the input parameters used by the Argon2id algorithm.
type Argon struct {
	// the amount of memory used by the algorithm (in kibibytes).
	Memory uint32
	// the number of iterations over the memory.
	Iterations uint32
	// the number of threads (or lanes) used by the algorithm.
	// Recommended value is between 1 and runtime.NumCPU().
	Parallelism uint8
	// length of the random salt. 16 bytes is recommended for password hashing.
	SaltLength uint32
	// length of the generated key. 16 bytes or more is recommended.
	KeyLength uint32
}

// returns an Argon2id hash of a plain-text password using the
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

// check argon2id hash.
func (a *Argon) Compare(what, with string) (match bool, err error) {
	match, _, err = a.matchHash(what, with)
	return match, err
}

// is data a hash.
func (a *Argon) IsHash(data string) bool {
	var conf, _, _, err = a.parseHash(data)
	return conf != nil && err == nil
}

// is like ComparePasswordAndHash, except it also returns the params that the hash was
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

// expects a hash created from this package, and parses it to return the params used to
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

// generate cryptographically random bytes.
func (a *Argon) generateRandomBytes(n uint32) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}
	return b, nil
}
