package pipe

import "servus/core"

type _ctx_user string

const CtxUser _ctx_user = "ELVEN_USER_PIPE"
const (
	AccessTypeAdminOnly  = "ELVEN_ACCESS_ADMIN_ONLY"
	AccessTypeReadOnly   = "ELVEN_ACCESS_READ_ONLY"
	AccessTypeAuthorized = "ELVEN_ACCESS_AUTHORIZED"
	userRoleAdmin        = "admin"
	userRoleUser         = "user"
)

var call *core.Instance

func Boot(c *core.Instance) {
	call = c
}
