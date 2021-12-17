package cryptor

import (
	"testing"
)

func hash(b *BCrypt) (string, error) {
	var dirty = "123"
	var hashed, err = b.Hash(dirty)
	if err != nil {
		return "", err
	}
	return hashed, nil
}

// Test_BCrypt_HashCompare - test hash and compare functions.
func Test_BCrypt_HashCompare(t *testing.T) {
	var cost = 14
	var dirty = "123"
	var dirtyWrong = "000"
	//
	var b = BCrypt{}
	b.New(cost)
	// hash.
	t.Log("testing hashing...")
	var hashed, err = hash(&b)
	if err != nil {
		t.Fatalf("failed to hash. Error: %v", err.Error())
	}
	t.Logf("hash: %v", hashed)
	// compare (true).
	t.Log("testing compare...")
	_, err = b.Compare(dirty, hashed)
	if err != nil {
		t.Fatalf("failed to compare. Error: %v", err.Error())
	}
	// compare (wrong).
	t.Log("testing compare (wrong)...")
	_, err = b.Compare(dirtyWrong, hashed)
	if err == nil {
		t.Fatalf("failed to compare wrong value. Error: %v", err.Error())
	}
}

// Test_BCrypt_HashCompare - make hash and get info.
func Test_BCrypt_Parse(t *testing.T) {
	var cost = 14
	var mockTrue = "$2a$14$/nQKx9M7VI/gFRgmaCiDh.LWco3PCva7Mk.Rv6mfnnAvzik1CvqyW"
	var mockWrongCases = []string{
		"123",
		"2a$14/nQKx9M7VI/gFRgmaCiDh.LWco3PCva7Mk.Rv6mfnnAvzik1CvqyW",
		"$2a$14$/nQKx9M7VI/gFRgmaCiDh.LWco3PCva7Mk.Rv6mfnnAvzik1CvqyW1234567890",
		"$2a/14//nQKx9M7VI/gFRgmaCiDh.LWco3PCva7Mk.Rv6mfnnAvzik1CvqyW",
	}
	//
	var b = BCrypt{}
	b.New(cost)
	// parse.
	t.Log("parsing (true)...")
	parsed, err := b.parse(mockTrue)
	if err != nil {
		t.Fatalf("failed to parse. Error: %v", err.Error())
	}
	if parsed.Algorithm != "2a" {
		t.Fatalf("wrong algorithm")
	}
	if parsed.Cost != cost {
		t.Fatalf("wrong cost")
	}
	if parsed.Salt != "/nQKx9M7VI/gFRgmaCiDh." {
		t.Fatalf("wrong salt")
	}
	if parsed.Hash != "LWco3PCva7Mk.Rv6mfnnAvzik1CvqyW" {
		t.Fatalf("wrong hash")
	}
	// wrong parsing.
	t.Log("parsing (wrong)...")
	for _, wrongCase := range mockWrongCases {
		parsed, err = b.parse(wrongCase)
		t.Log(err.Error())
		if parsed != nil || err == nil {
			t.Errorf("parsed: %v", parsed)
			t.Fatalf("wrong parsing failed")
		}
	}
}
