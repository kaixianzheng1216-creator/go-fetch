package middleware

import (
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
)

const httpRequestTimeout = 30 * time.Second

func UseHTTP(router chi.Router, sessionManager *scs.SessionManager) {
	router.Use(chimiddleware.RequestID)
	router.Use(chimiddleware.RealIP)
	router.Use(chimiddleware.Logger)
	router.Use(chimiddleware.Recoverer)
	router.Use(chimiddleware.Timeout(httpRequestTimeout))
	router.Use(sessionManager.LoadAndSave)
}
