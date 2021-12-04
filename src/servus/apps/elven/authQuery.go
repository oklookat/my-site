package elven

// bodyAuth - represents the body of the request that the user should send. Used in entityAuth login method.
type bodyAuth struct {
	Username string
	Password string
	Type     string
}
