package api

import (
	"net/http"
	"path/filepath"
	"runtime"

	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
)

// RegisterSwagger adds Swagger UI routes to the router
func RegisterSwagger(r *mux.Router) {
	// Get the absolute path to the docs directory
	_, b, _, _ := runtime.Caller(0)
	basepath := filepath.Dir(filepath.Dir(filepath.Dir(b)))
	docsPath := filepath.Join(basepath, "docs")

	// Serve the Swagger UI
	r.PathPrefix("/swagger/").Handler(httpSwagger.Handler(
		httpSwagger.URL("/swagger/doc.json"), // The URL to the API definition
		httpSwagger.DeepLinking(true),
		httpSwagger.DocExpansion("none"),
		httpSwagger.DomID("swagger-ui"),
	))

	// Serve the OpenAPI specification file directly
	r.PathPrefix("/swagger/doc.json").Handler(http.StripPrefix("/swagger", http.FileServer(http.Dir(docsPath))))

	// Redirect /docs to /swagger/index.html for convenience
	r.HandleFunc("/docs", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/swagger/index.html", http.StatusMovedPermanently)
	})
}
