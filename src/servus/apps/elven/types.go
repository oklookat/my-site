package elven

type ResponseContent struct {
	Meta struct {
		PerPage int    `json:"per_page"`
		Next    string `json:"next"`
	} `json:"meta"`
	Content interface{} `json:"content"`
}
