package api

import (
	"github.com/LeUrok/DS-lab2/payment-service/internal/auth"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func NewRouter(h *PaymentHandler, mw *auth.Middleware) *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/manage/health", h.Health)
	r.Route("/api/v1", func(r chi.Router) {
		r.Use(mw.Handler)
		r.Post("/payment", h.Create)
		r.Get("/payment/{paymentUid}", h.GetByUID)
		r.Delete("/payment/{paymentUid}", h.Cancel)
	})
	return r
}
