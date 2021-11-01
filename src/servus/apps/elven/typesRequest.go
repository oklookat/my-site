package elven

// ResponseContent - template for response.
type ResponseContent struct {
	Meta struct {
		PerPage     int `json:"per_page"`
		TotalPages  int `json:"total_pages"`
		CurrentPage int `json:"current_page"`
	} `json:"meta"`
	Data interface{} `json:"data"`
}

// BodyArticle - represents the body of the request that the user should send. Used in create and update methods.
type BodyArticle struct {
	IsPublished *bool           `json:"is_published"`
	Title       *string         `json:"title"`
	Content     *ArticleContent `json:"content"`
}

// bodyAuth - represents the body of the request that the user should send. Used in entityAuth login method.
type bodyAuth struct {
	Username string
	Password string
	Type     string
}
