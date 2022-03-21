package banhammer

import (
	"errors"
	"net"
)

func createError(message string) error {
	return errors.New("[banhammer] " + message)
}

func parseIP(ip string) net.IP {
	return net.ParseIP(ip)
}
