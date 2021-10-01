package elven

type PaginatedResponse struct {
	Meta struct {
		Total           int         `json:"total"`
		PerPage         int         `json:"per_page"`
		CurrentPage     int         `json:"current_page"`
		LastPage        int         `json:"last_page"`
		FirstPage       int         `json:"first_page"`
		FirstPageURL    string      `json:"first_page_url"`
		LastPageURL     string      `json:"last_page_url"`
		NextPageURL     interface{} `json:"next_page_url"`
		PreviousPageURL interface{} `json:"previous_page_url"`
	} `json:"meta"`
	Data []interface{} `json:"data"`
}

type controllerAuthLoginBody struct {
	Username string
	Password string
	Type     string
}

type TokenAuthData struct {
	UserID string
	TokenID string
}
