package cryptor

import "github.com/pkg/errors"

var ErrMakeBlock = errors.New("AESEncrypt / AESDecrypt: failed to make cipher block from key")
var ErrMakeGCM = errors.New("AESEncrypt / AESDecrypt: failed to make GCM from cipher block")
var ErrMakeNonceFromGCM = errors.New("AESEncrypt: make nonce from GCM failed")
var ErrDecodeHEX = errors.New("AESDecrypt: get hex from string failed. Is string contains only hex?")
var ErrDecryption = errors.New("AESDecrypt: decryption failed")

type AESError struct {
	HasErrors bool
	AdditionalErr error
	OriginalErr error
}