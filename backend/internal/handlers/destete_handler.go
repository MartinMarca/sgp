package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/martin/sgp/internal/authz"
	"github.com/martin/sgp/internal/repositories"
	"github.com/martin/sgp/internal/services"
	"github.com/martin/sgp/internal/utils"
)

// DesteteHandler maneja los endpoints de destetes
type DesteteHandler struct {
	service *services.DesteteService
	authDeps
}

// NewDesteteHandler crea una nueva instancia del handler
func NewDesteteHandler(service *services.DesteteService, authzSvc *authz.Service, repos *repositories.RepositoryContainer) *DesteteHandler {
	return &DesteteHandler{service: service, authDeps: newAuthDeps(authzSvc, repos)}
}

// Crear godoc
// POST /api/destetes
func (h *DesteteHandler) Crear(c *gin.Context) {
	var input services.CrearDesteteInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	if !requireCerdaForCreate(c, h.authz, h.repos, input.CerdaID, authz.PermRecursoWrite) {
		return
	}
	if input.LoteID != nil && !requireOnLote(c, h.authz, h.repos, *input.LoteID, authz.PermRecursoWrite) {
		return
	}

	destete, err := h.service.Crear(input)
	if err != nil {
		utils.ErrorResponse(c, mapErrorToStatus(err), err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "Destete registrado exitosamente", destete)
}

// ObtenerPorID godoc
// GET /api/destetes/:id
func (h *DesteteHandler) ObtenerPorID(c *gin.Context) {
	id, err := getIDParam(c, "id")
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "ID inválido")
		return
	}
	if !requireOnDestete(c, h.authz, h.repos, id, authz.PermGranjaRead) {
		return
	}

	destete, err := h.service.ObtenerPorID(id)
	if err != nil {
		utils.ErrorResponse(c, mapErrorToStatus(err), err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "", destete)
}

// ListarPorCerda godoc
// GET /api/cerdas/:id/destetes
func (h *DesteteHandler) ListarPorCerda(c *gin.Context) {
	cerdaID, err := getIDParam(c, "id")
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "ID de cerda inválido")
		return
	}
	if !requireOnCerda(c, h.authz, h.repos, cerdaID, authz.PermGranjaRead) {
		return
	}

	destetes, err := h.service.ListarPorCerda(cerdaID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "", destetes)
}

// ListarPorPeriodo godoc
// GET /api/destetes?mes=1&anio=2026&granja_id=1
func (h *DesteteHandler) ListarPorPeriodo(c *gin.Context) {
	if !requireGranjaQuery(c, h.authz, h.repos, authz.PermGranjaRead) {
		return
	}

	mes := getIntQuery(c, "mes", 0)
	anio := getIntQuery(c, "anio", 0)
	granjaID := getOptionalUintQuery(c, "granja_id")

	destetes, err := h.service.ListarPorPeriodo(granjaID, mes, anio)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "", destetes)
}

// Actualizar godoc
// PUT /api/destetes/:id
func (h *DesteteHandler) Actualizar(c *gin.Context) {
	id, err := getIDParam(c, "id")
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "ID inválido")
		return
	}
	if !requireOnDestete(c, h.authz, h.repos, id, authz.PermRecursoWrite) {
		return
	}

	var input services.ActualizarDesteteInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	destete, err := h.service.Actualizar(id, input)
	if err != nil {
		utils.ErrorResponse(c, mapErrorToStatus(err), err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Destete actualizado exitosamente", destete)
}

// GetEstadisticas godoc
// GET /api/destetes/estadisticas?mes=1&anio=2026&granja_id=1
func (h *DesteteHandler) GetEstadisticas(c *gin.Context) {
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
