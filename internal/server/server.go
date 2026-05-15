package server

import (
	"github.com/alexedwards/scs/v2"

	"github.com/kaixianzheng1216-creator/go-fetch/internal/server/session"
	"github.com/kaixianzheng1216-creator/go-fetch/internal/store"
)

type App struct {
	store    *store.Store
	sessions *scs.SessionManager
}

func New(dataStore *store.Store) *App {
	return &App{
		store:    dataStore,
		sessions: session.NewManager(dataStore),
	}
}
