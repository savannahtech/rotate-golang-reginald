package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *serverApplication) routes() *httprouter.Router {
	router := httprouter.New()

	router.HandlerFunc(http.MethodGet, "/v1/health", apiKeyMiddleware(app.healthCheckHandler))
	router.HandlerFunc(http.MethodGet, "/v1/stats", apiKeyMiddleware(app.logsHandler))
	router.HandlerFunc(http.MethodPost, "/v1/command", apiKeyMiddleware(app.cpuCommandHandler))
	return router
}

func apiKeyMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// to-do add this to envs
		apiKey := r.Header.Get("X-API-Key")
		if apiKey != "testing123" {
			http.Error(w, "Forbidden: Invalid API Key", http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, r)
	}
}
