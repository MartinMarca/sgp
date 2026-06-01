package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/martin/sgp/internal/authz"
	"github.com/martin/sgp/internal/repositories"
	"github.com/martin/sgp/internal/services"
	"github.com/martin/sgp/internal/utils"
)

// GranjaHandler maneja los endpoints de granjas
type GranjaHandler struct {
	service *services.GranjaService
	authz   *authz.Service
	repos   *repositories.RepositoryContainer
}

// NewGranjaHandler crea una nueva instancia del handler
func NewGranjaHandler(service *services.GranjaService, authzSvc *authz.Service, repos *repositories.RepositoryContainer) *GranjaHandler {
	return &GranjaHandler{service: service, authz: authzSvc, repos: repos}
}

// Crear godoc
// POST /api/granjas
func (h *GranjaHandler) Crear(c *gin.Context) {
	var input services.CrearGranjaInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	actor, err := getActor(c, h.repos)
	if err != nil {
		utils.ErrorResponse(c, mapErrorToStatus(err), err.Error())
		return
	}

	propietarioID, err := h.authz.ResolvePropietarioIDForCreate(actor, input.PropietarioID)
	if err != nil {
		utils.ErrorResponse(c, mapErrorToStatus(err), err.Error())
		return
	}

	granja, err := h.service.Crear(input, propietarioID)
	if err != nil {
		utils.ErrorResponse(c, mapErrorToStatus(err), err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "Granja creada exitosamente", granja)
}

// ObtenerPorID godoc
// GET /api/granjas/:id
func (h *GranjaHandler) ObtenerPorID(c *gin.Context) {
	id, err := getIDParam(c, "id")
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "ID inválido")
		return
	}

	if !requireGranja(c, h.authz, h.repos, id, authz.PermGranjaRead) {
		return
	}

	granja, err := h.service.ObtenerPorID(id)
	if err != nil {
		utils.ErrorResponse(c, mapErrorToStatus(err), err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "", granja)
}

// Listar godoc
// GET /api/granjas
func (h *GranjaHandler) Listar(c *gin.Context) {
	actor, err := getActor(c, h.repos)
	if err != nil {
		utils.ErrorResponse(c, mapErrorToStatus(err), err.Error())
		return
	}

	activo := getOptionalBoolQuery(c, "activo")
	granjas, err := h.service.ListarAccesibles(actor, activo)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "", granjas)
}

// ListarPorUsuario godoc
// GET /api/granjas/mis-granjas
func (h *GranjaHandler) ListarPorUsuario(c *gin.Context) {
	actor, err := getActor(c, h.repos)
	if err != nil {
		utils.ErrorResponse(c, mapErrorToStatus(err), err.Error())
		return
	}

	granjas, err := h.service.ListarPorUsuario(actor)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "", granjas)
}

// Actualizar godoc
// PUT /api/granjas/:id
func (h *GranjaHandler) Actualizar(c *gin.Context) {
	id, err := getIDParam(c, "id")
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "ID inválido")
		return
	}

	if !requireGranja(c, h.authz, h.repos, id, authz.PermGranjaWrite) {
		return
	}

	var input services.ActualizarGranjaInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	granja, err := h.service.Actualizar(id, input)
	if err != nil {
		utils.ErrorResponse(c, mapErrorToStatus(err), err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Granja actualizada exitosamente", granja)
}

// Eliminar godoc
// DELETE /api/granjas/:id
func (h *GranjaHandler) Eliminar(c *gin.Context) {
	id, err := getIDParam(c, "id")
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "ID inválido")
		return
	}

	if !requireGranja(c, h.authz, h.repos, id, authz.PermGranjaDelete) {
		return
	}

	if err := h.service.Eliminar(id); err != nil {
		utils.ErrorResponse(c, mapErrorToStatus(err), err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Granja dada de baja exitosamente", nil)
}

// GetEstadisticas godoc
// GET /api/granjas/:id/estadisticas
func (h *GranjaHandler) GetEstadisticas(c *gin.Context) {
	id, err := getIDParam(c, "id")
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "ID inválido")
		return
	}

	if !requireGranja(c, h.authz, h.repos, id, authz.PermStatsRead) {
		return
	}

	stats, err := h.service.GetEstadisticas(id)
	if err != nil {
		utils.ErrorResponse(c, mapErrorToStatus(err), err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "", stats)
}
