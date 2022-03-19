package banhammer

import (
	"errors"
	"net"
	"net/http"
	"strings"
)

func createError(message string) error {
	return errors.New("[banhammer] " + message)
}

func parseIP(ip string) net.IP {
	return net.ParseIP(ip)
}

func getIpByRequest(request *http.Request) net.IP {
	var ip string
	// try get real.
	var real = request.Header.Get("X-REAL-IP")
	var isRealEmpty = len(real) < 4
	if !isRealEmpty {
		ip = real
	} else {
		// try get forwarded.
		var forwarded = request.Header.Get("X-FORWARDED-FOR")
		var isForwardedEmpty = len(forwarded) < 4
		if isForwardedEmpty {
			return nil
		}
		// get first ip.
		ips := strings.Split(forwarded, ", ")
		for _, ipItem := range ips {
			ip = ipItem
			break
		}
	}
	return parseIP(ip)
}
