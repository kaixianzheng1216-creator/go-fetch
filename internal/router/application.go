package router

import (
	"github.com/alexedwards/scs/v2"

	"github.com/kaixianzheng1216-creator/go-fetch/internal/repository"
)

type Server struct {
	store    *repository.Store
	sessions *scs.SessionManager
}

func New(dataStore *repository.Store) *Server {
	return &Server{
		store:    dataStore,
		sessions: newSessionManager(dataStore),
	}
}
