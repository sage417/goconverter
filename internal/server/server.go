// internal/server/server.go
package server

import (
	"net/http"
)

type Server struct {
	router *http.ServeMux
}

func NewServer() *Server {
	s := &Server{
		router: http.NewServeMux(),
	}
	s.routes()
	return s
}

func (s *Server) routes() {
	s.router.HandleFunc("/convert", s.handleConvert())
}

func (s *Server) Run(addr string) error {
	return http.ListenAndServe(addr, s.router)
}
