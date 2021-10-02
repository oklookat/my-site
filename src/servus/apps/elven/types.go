package elven

type ResponseContent struct {
	Cursor struct {
		PerPage int    `json:"per_page"`
		Next    string `json:"next"`
	} `json:"cursor"`
	Content interface{} `json:"content"`
}

type controllerAuthLoginBody struct {
	Username string
	Password string
	Type     string
}

type TokenAuthData struct {
	UserID  string
	TokenID string
}
