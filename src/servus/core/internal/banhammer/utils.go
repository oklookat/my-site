package banhammer

import (
	"errors"
	"net"
)

func createError(message string) error {
	return errors.New(message)
}

func IsIpValid(ip string) bool {
	var res = parseIP(ip)
	return res != nil
}

func parseIP(ip string) net.IP {
	return net.ParseIP(ip)
}
