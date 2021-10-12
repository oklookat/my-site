package elven

import (
	"github.com/oklog/ulid/v2"
	"math/rand"
	"net/http"
	"strings"
	"time"
)

// objectUtils - useful utilities.
type objectUtils struct {
}

// getIP - get IP by request.
func (u *objectUtils) getIP(request *http.Request) (ip string) {
	ip = ""
	var ips = strings.Split(request.Header.Get("X-FORWARDED-FOR"), ", ")
	for _, theIP := range ips {
		if theIP != "" {
			ip = theIP
			break
		}
	}
	if ip == "" {
		ip = request.RemoteAddr
	}
	return
}

// generateULID - returns unique string like 1GFGVSSRTHYWW52GVXZ.
func (u *objectUtils) generateULID() (ul string, err error) {
	current := time.Now()
	entropy := ulid.Monotonic(rand.New(rand.NewSource(current.UnixNano())), 0)
	ulType, err := ulid.New(ulid.Timestamp(current), entropy)
	if err != nil {
		return "", err
	}
	ul = ulType.String()
	return
}
