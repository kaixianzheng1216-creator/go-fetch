package middleware

import (
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
)

func UseHTTP(router chi.Router, sessionManager *scs.SessionManager) {
	router.Use(chimiddleware.RealIP)
	router.Use(chimiddleware.RequestID)
	router.Use(chimiddleware.Recoverer)
	router.Use(chimiddleware.Logger)
	router.Use(chimiddleware.Timeout(60 * time.Second))
	router.Use(sessionManager.LoadAndSave)
}
