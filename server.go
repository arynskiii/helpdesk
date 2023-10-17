package smarthelpdesc

import (
	"context"
	"net/http"
	"time"
)

type Server struct {
	httpServer http.Server
}

func (s *Server) Run(port string, handler http.Handler) error {
	s.httpServer = http.Server{
		Addr:           ":" + port,
		MaxHeaderBytes: 1 << 20,
		Handler:        handler,
		WriteTimeout:   10 * time.Second,
		ReadTimeout:    10 * time.Second,
	}
	return s.httpServer.ListenAndServe()

}

func (s *Server) ShutDown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
