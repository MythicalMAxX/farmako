package swagger

import (
	"net/http"
	"os"
	"path/filepath"

	"github.com/go-chi/chi/v5"
)

// Register adds the Swagger routes to the router
func Register(r chi.Router) {
	// Get the working directory for Swagger files
	workDir, _ := filepath.Abs("./swagger")
	swaggerSpecPath := filepath.Join(workDir, "doc.json")

	// Ensure the swagger doc exists
	if _, err := os.Stat(swaggerSpecPath); os.IsNotExist(err) {
		// If not found, log a warning - in a real app, we would handle this better
		println("Warning: Swagger documentation not found at", swaggerSpecPath)
	}

	// Handle specific routes first (in order of precedence)

	// Serve custom index.html file to fix the 'list is not defined' error
	r.Get("/swagger/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		http.ServeFile(w, r, filepath.Join(workDir, "index.html"))
	})

	// Specifically handle index.html to make sure our fixed version is used
	r.Get("/swagger/index.html", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		http.ServeFile(w, r, filepath.Join(workDir, "index.html"))
	})

	// Serve swagger spec at /api/v1/swagger/doc.json
	r.Get("/swagger/doc.json", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		http.ServeFile(w, r, swaggerSpecPath)
	})

	// Serve static assets (JS and CSS)
	r.Get("/swagger/swagger-ui.css", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/css")
		http.ServeFile(w, r, filepath.Join(workDir, "swagger-ui.css"))
	})

	r.Get("/swagger/swagger-ui-bundle.js", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/javascript")
		http.ServeFile(w, r, filepath.Join(workDir, "swagger-ui-bundle.js"))
	})

	r.Get("/swagger/swagger-ui-standalone-preset.js", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/javascript")
		http.ServeFile(w, r, filepath.Join(workDir, "swagger-ui-standalone-preset.js"))
	})

	// Serve favicon files
	r.Get("/swagger/favicon-16x16.png", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "image/png")
		http.ServeFile(w, r, filepath.Join(workDir, "favicon-16x16.png"))
	})

	r.Get("/swagger/favicon-32x32.png", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "image/png")
		http.ServeFile(w, r, filepath.Join(workDir, "favicon-32x32.png"))
	})

	// We no longer need the default handler as we've explicitly defined all necessary routes
	// This avoids any conflicts between our custom handlers and the default Swagger UI
}
