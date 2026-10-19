package server

import (
	"dle-sitemap/database"
	"dle-sitemap/sitemap"
	"fmt"
	"log"
	"net/http"
)

type Service struct {
	port           string
	server         http.Server
	sitemapService *sitemap.Service
	dbService      *database.Service
}

func (s *Service) Run() {
	addr := fmt.Sprintf(":%s", s.port)
	log.Println("Starting proxy server on", addr)
	err := s.server.ListenAndServe()
	if err != nil {
		log.Fatal("Error starting proxy server: ", err)
	}
}

func NewService(port string, sitemapService *sitemap.Service, dbService *database.Service) (s *Service, err error) {

	s = &Service{
		port:           port,
		sitemapService: sitemapService,
		dbService:      dbService,
	}

	s.server = http.Server{
		Addr: fmt.Sprintf(":%s", port),
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/health" {
				s.health(w, r)
				return
			}
			s.handler(w, r)
		}),
	}

	return
}
