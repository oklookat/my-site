package cryptor

import (
	"errors"
	"strconv"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

type BCrypt struct {
	// Cost - see bcrypt.GenerateFromPassword.
	Cost int
}

type BCryptParsed struct {
	Algorithm string
	Cost      int
	Salt      string
	Hash      string
}

func (b *BCrypt) New(cost int) {
	b.Cost = cost
}

func (b *BCrypt) Hash(data string) (hashed string, err error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(data), b.Cost)
	return string(bytes), err
}

func (b *BCrypt) Compare(what string, with string) (match bool, err error) {
	err = bcrypt.CompareHashAndPassword([]byte(with), []byte(what))
	match = err == nil
	return
}

func (b *BCrypt) IsHash(data string) bool {
	_, err := b.parse(data)
	return err == nil
}

func (b *BCrypt) parse(hash string) (*BCryptParsed, error) {
	if len(hash) < 12 {
		return nil, errors.New("[bcrypt/parse]: too small hash length")
	}
	if !strings.HasPrefix(hash, "$") {
		return nil, errors.New("[bcrypt/parse]: hash not have dollar prefix")
	}
	// sliced output like: [ 2a 14 eyjC4iZWaXvn.yEw9gzrsO2qIEFLC0jTQwk1ttn.WGl0O/ogFSXJ6].
	var sliced = strings.Split(hash, "$")
	if len(sliced) != 4 {
		return nil, errors.New("[bcrypt/parse]: invalid hash slice length / dollar split err")
	}
	var parsed = &BCryptParsed{}
	// algorithm.
	var algorithm = sliced[1]
	if algorithm != "2a" {
		return nil, errors.New("[bcrypt/parse]: invalid algorithm")
	}
	parsed.Algorithm = algorithm
	// cost.
	var cost = sliced[2]
	costInt, err := strconv.Atoi(cost)
	if err != nil {
		return nil, errors.New("[bcrypt/parse]: cost NaN")
	}
	parsed.Cost = costInt
	// payload.
	var payload = sliced[3]
	if len(payload) != 53 {
		return nil, errors.New("[bcrypt/parse]: invalid payload length")
	}
	var charCounter = 0
	var saltLen = 22
	for _, r := range payload {
		if charCounter < saltLen {
			parsed.Salt = parsed.Salt + string(r)
		} else {
			parsed.Hash = parsed.Hash + string(r)
		}
		charCounter++
	}
	return parsed, err
}
