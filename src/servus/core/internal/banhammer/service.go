package banhammer

import (
	"net"
	"net/http"
	"strings"
)

type Service struct {
	db *SQLite
}

func (s *Service) New(db *SQLite) {
	s.db = db
}

func (s *Service) GetMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
			// get ip.
			var ip = s.GetIpByRequest(request)
			if ip == nil {
				next.ServeHTTP(response, request)
				return
			}

			// check ban.
			var ipString = ip.String()
			entry, err := s.db.GetEntry(ipString)
			var isEntryInvalid = entry == nil || err != nil
			if isEntryInvalid {
				next.ServeHTTP(response, request)
				return
			}

			// 403 if banned.
			if entry.IsBanned {
				response.WriteHeader(403)
				response.Write(nil)
				return
			}

			// continue if not.
			next.ServeHTTP(response, request)
			return
		})
	}
}

// get IP by request by X-REAL-IP / X-FORWARDED-FOR
//
// returns nil if failed to get IP.
func (s *Service) GetIpByRequest(request *http.Request) net.IP {
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
