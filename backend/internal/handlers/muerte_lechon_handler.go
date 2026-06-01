package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/martin/sgp/internal/authz"
	"github.com/martin/sgp/internal/repositories"
	"github.com/martin/sgp/internal/services"
	"github.com/martin/sgp/internal/utils"
)

// MuerteLechonHandler maneja los endpoints de muertes de lechones
type MuerteLechonHandler struct {
	service *services.MuerteLechonService
	authDeps
}

// NewMuerteLechonHandler crea una nueva instancia del handler
func NewMuerteLechonHandler(service *services.MuerteLechonService, authzSvc *authz.Service, repos *repositories.RepositoryContainer) *MuerteLechonHandler {
	return &MuerteLechonHandler{service: service, authDeps: newAuthDeps(authzSvc, repos)}
}

// Crear godoc
// POST /api/muertes-lechones
func (h *MuerteLechonHandler) Crear(c *gin.Context) {
	var input services.CrearMuerteLechonInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	if !requireGranjaForCreate(c, h.authz, h.repos, input.GranjaID, authz.PermRecursoWrite) {
		return
	}

	muerte, err := h.service.Crear(input)
	if err != nil {
		utils.ErrorResponse(c, mapErrorToStatus(err), err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "Muerte de lechón registrada exitosamente", muerte)
}

// ObtenerPorID godoc
// GET /api/muertes-lechones/:id
func (h *MuerteLechonHandler) ObtenerPorID(c *gin.Context) {
	id, err := getIDParam(c, "id")
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "ID inválido")
		return
	}
	if !requireOnMuerte(c, h.authz, h.repos, id, authz.PermGranjaRead) {
		return
	}

	muerte, err := h.service.ObtenerPorID(id)
	if err != nil {
		utils.ErrorResponse(c, mapErrorToStatus(err), err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "", muerte)
}

// ListarPorGranja godoc
// GET /api/granjas/:id/muertes-lechones
func (h *MuerteLechonHandler) ListarPorGranja(c *gin.Context) {
	granjaID, err := getIDParam(c, "id")
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "ID de granja inválido")
		return
	}
	if !requireGranja(c, h.authz, h.repos, granjaID, authz.PermGranjaRead) {
		return
	}

	muertes, err := h.service.ListarPorGranja(granjaID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "", muertes)
}

// ListarPorParto godoc
// GET /api/partos/:id/muertes-lechones
func (h *MuerteLechonHandler) ListarPorParto(c *gin.Context) {
	partoID, err := getIDParam(c, "id")
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "ID de parto inválido")
		return
	}
	if !requireOnParto(c, h.authz, h.repos, partoID, authz.PermGranjaRead) {
		return
	}

	muertes, err := h.service.ListarPorParto(partoID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "", muertes)
}

// ListarPorCorral godoc
// GET /api/corrales/:id/muertes
func (h *MuerteLechonHandler) ListarPorCorral(c *gin.Context) {
	corralID, err := getIDParam(c, "id")
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "ID de corral inválido")
		return
	}
	if !requireOnCorral(c, h.authz, h.repos, corralID, authz.PermGranjaRead) {
		return
	}

	muertes, err := h.service.ListarPorCorral(corralID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "", muertes)
}

// ListarPorPeriodo godoc
// GET /api/muertes-lechones?mes=1&anio=2026&granja_id=1
func (h *MuerteLechonHandler) ListarPorPeriodo(c *gin.Context) {
	if !requireGranjaQuery(c, h.authz, h.repos, authz.PermGranjaRead) {
		return
	}

	mes := getIntQuery(c, "mes", 0)
	anio := getIntQuery(c, "anio", 0)
	granjaID := getOptionalUintQuery(c, "granja_id")

	muertes, err := h.service.ListarPorPeriodo(granjaID, mes, anio)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "", muertes)
}

// Actualizar godoc
// PUT /api/muertes-lechones/:id
func (h *MuerteLechonHandler) Actualizar(c *gin.Context) {
	id, err := getIDParam(c, "id")
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "ID inválido")
		return
	}
	if !requireOnMuerte(c, h.authz, h.repos, id, authz.PermRecursoWrite) {
		return
	}

	var input services.ActualizarMuerteLechonInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	muerte, err := h.service.Actualizar(id, input)
	if err != nil {
		utils.ErrorResponse(c, mapErrorToStatus(err), err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Registro de muerte actualizado exitosamente", muerte)
}

// Eliminar godoc
// DELETE /api/muertes-lechones/:id
func (h *MuerteLechonHandler) Eliminar(c *gin.Context) {
	id, err := getIDParam(c, "id")
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "ID inválido")
		return
	}
	if !requireOnMuerte(c, h.authz, h.repos, id, authz.PermRecursoWrite) {
		return
	}

	if err := h.service.Eliminar(id); err != nil {
		utils.ErrorResponse(c, mapErrorToStatus(err), err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Registro de muerte eliminado exitosamente", nil)
}

// GetEstadisticas godoc
// GET /api/muertes-lechones/estadisticas?mes=1&anio=2026&granja_id=1
func (h *MuerteLechonHandler) GetEstadisticas(c *gin.Context) {
	if !requireGranjaQuery(c, h.authz, h.repos, authz.PermStatsRead) {
		return
	}

	mes := getIntQuery(c, "mes", 0)
	anio := getIntQuery(c, "anio", 0)
	granjaID := getOptionalUintQuery(c, "granja_id")

	stats, err := h.service.GetEstadisticas(granjaID, mes, anio)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "", stats)
}
