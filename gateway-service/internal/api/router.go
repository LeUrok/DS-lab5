package api

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/LeUrok/DS-lab2/gateway-service/internal/auth"

)

func NewRouter(h *Handler, ah *AuthHandler, mw *auth.Middleware) *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/manage/health", h.Health)
	r.Route("/api/v1", func(r chi.Router) {
		r.Post("/authorize", ah.Authorize)
		r.Get("/callback", ah.Callback)

		r.Group(func(r chi.Router) {
			r.Use(mw.Handler)

			r.Get("/cars", h.GetCars)

			r.Get("/rental", h.GetRentals)
			r.Post("/rental", h.CreateRental)
			r.Get("/rental/{rentalUid}", h.GetRentalByUID)
			r.Post("/rental/{rentalUid}/finish", h.FinishRental)
			r.Delete("/rental/{rentalUid}", h.CancelRental)
		})
	})
	return r
}
