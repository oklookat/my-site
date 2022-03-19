package banhammer

import (
	"net/http"
)

type Service struct {
	list Lister
}

func (s *Service) New(l Lister) {
	s.list = l
}

func (s *Service) GetMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
			// get ip.
			var ip = getIpByRequest(request)
			if ip == nil {
				next.ServeHTTP(response, request)
				return
			}

			// get list.
			var ips, err = s.list.GetList()
			var isListInvalid = ips == nil || ips.List == nil || err != nil
			if isListInvalid {
				next.ServeHTTP(response, request)
				return
			}

			// check ban.
			var ipString = ip.String()
			entry, err := s.list.GetEntry(ipString)
			var isEntryInvalid = entry == nil || err != nil
			if isEntryInvalid {
				next.ServeHTTP(response, request)
				return
			}
			var isBanned = entry.IsBanned || entry.WarnsCount == 3
			if isBanned {
				response.WriteHeader(403)
				response.Write(nil)
				return
			}
			next.ServeHTTP(response, request)
			return
		})
	}
}
