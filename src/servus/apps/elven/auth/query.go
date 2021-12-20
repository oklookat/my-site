package auth

// body - represents the body of the request that the user should send. Used in entityAuth login method.
type body struct {
	Username string
	Password string
	Type     string
}
