package handlers

import (
	"github.com/martin/sgp/internal/authz"
	"github.com/martin/sgp/internal/repositories"
)

// authDeps agrupa dependencias de autorización inyectadas en handlers.
type authDeps struct {
	authz *authz.Service
	repos *repositories.RepositoryContainer
}

func newAuthDeps(authzSvc *authz.Service, repos *repositories.RepositoryContainer) authDeps {
	return authDeps{authz: authzSvc, repos: repos}
}
