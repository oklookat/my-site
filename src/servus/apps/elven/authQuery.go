package elven

// authBody - represents the body of the request that the user should send. Used in entityAuth login method.
type authBody struct {
	Username string
	Password string
	Type     string
}
