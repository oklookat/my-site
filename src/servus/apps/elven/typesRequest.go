package elven

// ResponseContent - template for response.
type ResponseContent struct {
	Meta struct {
		PerPage int    `json:"per_page"`
		Next    string `json:"next"`
	} `json:"meta"`
	Data interface{} `json:"data"`
}


// BodyArticle - represents the body of the request that the user should send. Used in create and update methods.
type BodyArticle struct {
	IsPublished *bool  `json:"is_published"`
	Title       string `json:"title"`
	Content     struct {
		Time   int64 `json:"time"`
		Blocks []struct {
			ID   string      `json:"id"`
			Type string      `json:"type"`
			Data interface{} `json:"data"`
		} `json:"blocks"`
		Version string `json:"version"`
	} `json:"content"`
}

// bodyAuth - represents the body of the request that the user should send. Used in entityAuth login method.
type bodyAuth struct {
	Username string
	Password string
	Type     string
}
