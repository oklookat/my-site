package elven

import "time"

type modelToken struct {
	id       string
	userID    string
	token     string
	lastIP    string
	lastAgent string
	authIP    string
	authAgent string
	createdAt time.Time
	updatedAt time.Time
}

type modelUser struct {
	id        string
	role      string
	username  string
	password  string
	regIP     string
	regAgent  string
	createdAt time.Time
	updatedAt time.Time
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
