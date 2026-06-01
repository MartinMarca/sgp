package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/martin/sgp/internal/authz"
	"github.com/martin/sgp/internal/repositories"
	"github.com/martin/sgp/internal/services"
	"github.com/martin/sgp/internal/utils"
)

// EstadisticasHandler maneja los endpoints de estadísticas
type EstadisticasHandler struct {
	service *services.EstadisticasService
	authDeps
}

// NewEstadisticasHandler crea una nueva instancia del handler
func NewEstadisticasHandler(service *services.EstadisticasService, authzSvc *authz.Service, repos *repositories.RepositoryContainer) *EstadisticasHandler {
	return &EstadisticasHandler{service: service, authDeps: newAuthDeps(authzSvc, repos)}
}

// GetResumenGranja godoc
// GET /api/estadisticas/granja/:id
func (h *EstadisticasHandler) GetResumenGranja(c *gin.Context) {
	granjaID, err := getIDParam(c, "id")
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "ID de granja inválido")
		return
	}
	if !requireGranja(c, h.authz, h.repos, granjaID, authz.PermStatsRead) {
		return
	}

	resumen, err := h.service.GetResumenGranja(granjaID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "", resumen)
}

// GetEstadisticasPeriodo godoc
// GET /api/estadisticas/periodo?mes=1&anio=2026&granja_id=1
func (h *EstadisticasHandler) GetEstadisticasPeriodo(c *gin.Context) {
	if !requireGranjaQuery(c, h.authz, h.repos, authz.PermStatsRead) {
		return
	}

	mes := getIntQuery(c, "mes", 0)
	anio := getIntQuery(c, "anio", 0)
	granjaID := getOptionalUintQuery(c, "granja_id")

	stats, err := h.service.GetEstadisticasPeriodo(granjaID, mes, anio)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "", stats)
}
