package banhammer

import "net/http"

type Service struct {
	list Lister
}

func (s *Service) New(l Lister) {
	s.list = l
}

func (s *Service) Amnesty(ip string) error {
	var ips, err = s.list.GetList()
	if err != nil {
		return err
	}
	delete(ips.List, ip)
	return err
}

func (s *Service) GetMiddleware() func(http.Handler) http.Handler {

}
