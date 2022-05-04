package app

import (
	"context"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

type Server struct {
	router *mux.Router
	srv    *http.Server
}

func NewServer(address string) *Server {
	r := mux.NewRouter()
	return &Server{
		router: r,
		srv: &http.Server{
			Addr:    address,
			Handler: cors.Default().Handler(r),
		},
	}
}

func (s *Server) Start() error {
	//register the required paths
	s.router.HandleFunc("/v1/inventory", s.getInventory).Methods("POST")
	s.router.HandleFunc("/v1/catalogue", s.getCatalogue).Methods("GET")
	s.router.HandleFunc("/v1/purchase", s.purchase).Methods("POST")
	s.router.HandleFunc("/test", s.test)
	// start listening for calls
	return s.srv.ListenAndServe()
}

func (s *Server) Shutdown() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return s.srv.Shutdown(ctx)
}
