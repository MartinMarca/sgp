package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/martin/sgp/internal/authz"
	"github.com/martin/sgp/internal/repositories"
	"github.com/martin/sgp/internal/services"
	"github.com/martin/sgp/internal/utils"
)

// CalendarioHandler maneja los endpoints del calendario
type CalendarioHandler struct {
	service *services.CalendarioService
	authDeps
}

// NewCalendarioHandler crea una nueva instancia del handler
func NewCalendarioHandler(service *services.CalendarioService, authzSvc *authz.Service, repos *repositories.RepositoryContainer) *CalendarioHandler {
	return &CalendarioHandler{service: service, authDeps: newAuthDeps(authzSvc, repos)}
}

// GetEventosFuturos godoc
// GET /api/calendario?granja_id=1&dias=30
func (h *CalendarioHandler) GetEventosFuturos(c *gin.Context) {
	if !requireGranjaQuery(c, h.authz, h.repos, authz.PermGranjaRead) {
		return
	}

	granjaID := getOptionalUintQuery(c, "granja_id")
	dias := getIntQuery(c, "dias", 30)

	eventos, err := h.service.GetEventosFuturos(granjaID, dias)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "", eventos)
}
