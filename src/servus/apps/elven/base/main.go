package base

// ResponseContent - template for response.
type ResponseContent struct {
	Meta struct {
		PerPage     int `json:"per_page"`
		TotalPages  int `json:"total_pages"`
		CurrentPage int `json:"current_page"`
	} `json:"meta"`
	Data interface{} `json:"data"`
}
