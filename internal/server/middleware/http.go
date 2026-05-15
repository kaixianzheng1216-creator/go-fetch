package middleware

import (
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
)

const httpRequestTimeout = 30 * time.Second

func UseHTTP(r chi.Router, sessions *scs.SessionManager) {
	r.Use(chimiddleware.RequestID)
	r.Use(chimiddleware.RealIP)
	r.Use(chimiddleware.Logger)
	r.Use(chimiddleware.Recoverer)
	r.Use(chimiddleware.Timeout(httpRequestTimeout))
	r.Use(sessions.LoadAndSave)
}
