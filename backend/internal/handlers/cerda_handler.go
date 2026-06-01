package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/martin/sgp/internal/authz"
	"github.com/martin/sgp/internal/repositories"
	"github.com/martin/sgp/internal/services"
	"github.com/martin/sgp/internal/utils"
)

// CerdaHandler maneja los endpoints de cerdas
type CerdaHandler struct {
	service *services.CerdaService
	authDeps
}

// NewCerdaHandler crea una nueva instancia del handler
func NewCerdaHandler(service *services.CerdaService, authzSvc *authz.Service, repos *repositories.RepositoryContainer) *CerdaHandler {
	return &CerdaHandler{service: service, authDeps: newAuthDeps(authzSvc, repos)}
}

// CrearEnGranja godoc
// POST /api/granjas/:id/cerdas
func (h *CerdaHandler) CrearEnGranja(c *gin.Context) {
	granjaID, err := getIDParam(c, "id")
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "ID de granja inválido")
		return
	}
	if !requireGranja(c, h.authz, h.repos, granjaID, authz.PermCerdaWrite) {
		return
	}

	var input services.CrearCerdaInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	input.GranjaID = granjaID

	cerda, err := h.service.Crear(input)
	if err != nil {
		utils.ErrorResponse(c, mapErrorToStatus(err), err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "Cerda registrada exitosamente", cerda)
}

// ObtenerPorID godoc
// GET /api/cerdas/:id
func (h *CerdaHandler) ObtenerPorID(c *gin.Context) {
	id, err := getIDParam(c, "id")
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "ID inválido")
		return
	}
	if !requireOnCerda(c, h.authz, h.repos, id, authz.PermGranjaRead) {
		return
	}

	cerda, err := h.service.ObtenerPorID(id)
	if err != nil {
		utils.ErrorResponse(c, mapErrorToStatus(err), err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "", cerda)
}

// ListarPorGranja godoc
// GET /api/granjas/:id/cerdas?estado=disponible&activo=true
func (h *CerdaHandler) ListarPorGranja(c *gin.Context) {
	granjaID, err := getIDParam(c, "id")
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "ID de granja inválido")
		return
	}
	if !requireGranja(c, h.authz, h.repos, granjaID, authz.PermGranjaRead) {
		return
	}

	estado := getOptionalStringQuery(c, "estado")
	activo := getOptionalBoolQuery(c, "activo")

	cerdas, err := h.service.ListarPorGranja(granjaID, estado, activo)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "", cerdas)
}

// ListarPorEstado godoc
// GET /api/cerdas?estado=gestacion&granja_id=1
func (h *CerdaHandler) ListarPorEstado(c *gin.Context) {
	estado := c.Query("estado")
	if estado == "" {
		utils.ErrorResponse(c, http.StatusBadRequest, "El parámetro estado es requerido")
		return
	}
	if !requireGranjaQuery(c, h.authz, h.repos, authz.PermGranjaRead) {
		return
	}

	granjaID := getOptionalUintQuery(c, "granja_id")

	cerdas, err := h.service.ListarPorEstado(estado, granjaID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "", cerdas)
}

// Actualizar godoc
// PUT /api/cerdas/:id
func (h *CerdaHandler) Actualizar(c *gin.Context) {
	id, err := getIDParam(c, "id")
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "ID inválido")
		return
	}
	if !requireOnCerda(c, h.authz, h.repos, id, authz.PermCerdaWrite) {
		return
	}

	var input services.ActualizarCerdaInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	cerda, err := h.service.Actualizar(id, input)
	if err != nil {
		utils.ErrorResponse(c, mapErrorToStatus(err), err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Cerda actualizada exitosamente", cerda)
}

// DarDeBaja godoc
// POST /api/cerdas/:id/baja
func (h *CerdaHandler) DarDeBaja(c *gin.Context) {
	id, err := getIDParam(c, "id")
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "ID inválido")
		return
	}
	if !requireOnCerda(c, h.authz, h.repos, id, authz.PermCerdaDelete) {
		return
	}

	var input services.BajaCerdaInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	cerda, err := h.service.DarDeBaja(id, input)
	if err != nil {
		utils.ErrorResponse(c, mapErrorToStatus(err), err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Cerda dada de baja exitosamente", cerda)
}

// GetHistorial godoc
// GET /api/cerdas/:id/historial
func (h *CerdaHandler) GetHistorial(c *gin.Context) {
	id, err := getIDParam(c, "id")
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "ID inválido")
		return
	}
	if !requireOnCerda(c, h.authz, h.repos, id, authz.PermGranjaRead) {
		return
	}

	historial, err := h.service.GetHistorial(id)
	if err != nil {
		utils.ErrorResponse(c, mapErrorToStatus(err), err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "", historial)
}

// GetEstadisticas godoc
// GET /api/cerdas/:id/estadisticas
func (h *CerdaHandler) GetEstadisticas(c *gin.Context) {
	id, err := getIDParam(c, "id")
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "ID inválido")
		return
	}
	if !requireOnCerda(c, h.authz, h.repos, id, authz.PermStatsRead) {
		return
	}

	stats, err := h.service.GetEstadisticas(id)
	if err != nil {
		utils.ErrorResponse(c, mapErrorToStatus(err), err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "", stats)
}
