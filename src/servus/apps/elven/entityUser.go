package elven

import (
	"encoding/json"
	"net/http"
	"servus/core/external/errorMan"
)

// entityUser - manage users.
type entityUser struct {
	*entityBase
}

type ResponseUser struct {
	IsAdmin bool `json:"is_admin"`
	Username string `json:"username"`
	LastIP *string `json:"last_ip"`
	LastAgent *string `json:"last_agent"`
}

func (u *entityUser) controllerGetMe(response http.ResponseWriter, request *http.Request) {
	auth := PipeAuth{}
	auth.get(request)
	var resp = ResponseUser{}
	resp.IsAdmin = auth.IsAdmin
	resp.Username = auth.User.Username
	resp.LastIP = auth.Token.LastIP
	resp.LastAgent = auth.Token.LastAgent
	bytes, err := json.Marshal(resp)
	if err != nil {
		u.Send(response, errorMan.ThrowServer(), 500)
		return
	}
	u.Send(response, string(bytes), 200)
	return
}