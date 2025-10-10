package http

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/yourname/energy-usage-api/internal/http/handlers"
)

func NewRouter(ah *handlers.AuthHandler, mh *handlers.MetersHandler, rh *handlers.ReadingsHandler, secret string) http.Handler {
	r := chi.NewRouter()
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"Authorization", "Content-Type"},
	}))

	r.Post("/v1/auth/signup", ah.Signup)
	r.Post("/v1/auth/login", ah.Login)

	r.Group(func(pr chi.Router) {
		pr.Use(AuthMiddleware(secret))
		pr.Route("/v1", func(api chi.Router) {
			api.Get("/meters", mh.List)
			api.Post("/meters", mh.Create)
			api.Post("/meters/{meterID}/readings", rh.Add)
			api.Get("/meters/{meterID}/usage", rh.Usage)
		})
	})

	return r
}
