package cryptor

type AESError struct {
	HaveErrors bool
	ErrorCode string
	// error from crypto
	Error error
}
