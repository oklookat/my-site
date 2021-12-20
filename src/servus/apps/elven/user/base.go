package user

type Response struct {
	IsAdmin  bool   `json:"is_admin"`
	Username string `json:"username"`
}
