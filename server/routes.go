package server

import (
	"OrderServer/services/Health"
	"github.com/go-chi/chi"
)

func (srv *Server) InjectRoutes() *chi.Mux {
	r := chi.NewRouter()
	r.Route("/", func(order chi.Router) {
		order.Get("/health", Health.Health)
		order.Route("/order", srv.OrderHandler.Serve)
	})
	return r
}
