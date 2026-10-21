package api

import (
	"github.com/LeUrok/DS-lab2/cars-service/internal/auth"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func NewRouter(h *CarHandler, mw *auth.Middleware) *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/manage/health", h.Health)
	r.Route("/api/v1", func(r chi.Router) {
		r.Use(mw.Handler)
		r.Get("/cars", h.GetAll)
		r.Get("/cars/{carUid}", h.GetByUID)
		r.Post("/cars/{carUid}/reserve", h.Reserve)
		r.Post("/cars/{carUid}/unreserve", h.Unreserve)
	})
	return r
}
