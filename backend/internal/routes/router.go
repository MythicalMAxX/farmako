package routes

import (
	"net/http"
	"time"

	"github.com/farmako/coupon-system/internal/handlers"
	"github.com/farmako/coupon-system/internal/services"
	"github.com/farmako/coupon-system/internal/swagger"
	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

// SetupRouter configures and returns a router for the API
func SetupRouter(couponService *services.CouponService) chi.Router {
	r := chi.NewRouter()

	// Set up middleware
	r.Use(chimiddleware.RequestID)
	r.Use(chimiddleware.RealIP)
	r.Use(chimiddleware.Logger)
	r.Use(chimiddleware.Recoverer)
	r.Use(chimiddleware.Timeout(60 * time.Second))

	// Set up CORS
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// Health check endpoint
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// API version prefix
	r.Route("/api/v1", func(r chi.Router) {
		// Initialize handlers
		couponHandler := handlers.NewCouponHandler(couponService)

		// Register routes
		couponHandler.RegisterRoutes(r)

		// Register Swagger routes
		swagger.Register(r)
	})

	return r
}
