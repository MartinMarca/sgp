package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/martin/sgp/internal/services"
	"github.com/martin/sgp/internal/utils"
)

// PartoHandler maneja los endpoints de partos
type PartoHandler struct {
	service *services.PartoService
}

// NewPartoHandler crea una nueva instancia del handler
func NewPartoHandler(service *services.PartoService) *PartoHandler {
	return &PartoHandler{service: service}
}

// Crear godoc
// POST /api/partos
func (h *PartoHandler) Crear(c *gin.Context) {
	var input services.CrearPartoInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	parto, err := h.service.Crear(input)
	if err != nil {
		utils.ErrorResponse(c, mapErrorToStatus(err), err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "Parto registrado exitosamente", parto)
}

// ObtenerPorID godoc
// GET /api/partos/:id
func (h *PartoHandler) ObtenerPorID(c *gin.Context) {
	id, err := getIDParam(c, "id")
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "ID inválido")
		return
	}

	parto, err := h.service.ObtenerPorID(id)
	if err != nil {
		utils.ErrorResponse(c, mapErrorToStatus(err), err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "", parto)
}

// ListarPorCerda godoc
// GET /api/cerdas/:id/partos
func (h *PartoHandler) ListarPorCerda(c *gin.Context) {
	cerdaID, err := getIDParam(c, "id")
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "ID de cerda inválido")
		return
	}

	partos, err := h.service.ListarPorCerda(cerdaID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "", partos)
}

// ListarPorPeriodo godoc
// GET /api/partos?mes=1&anio=2026&granja_id=1
func (h *PartoHandler) ListarPorPeriodo(c *gin.Context) {
	mes := getIntQuery(c, "mes", 0)
	anio := getIntQuery(c, "anio", 0)
	granjaID := getOptionalUintQuery(c, "granja_id")

	partos, err := h.service.ListarPorPeriodo(granjaID, mes, anio)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "", partos)
}

// Actualizar godoc
// PUT /api/partos/:id
func (h *PartoHandler) Actualizar(c *gin.Context) {
	id, err := getIDParam(c, "id")
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "ID inválido")
		return
	}

	var input services.ActualizarPartoInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	parto, err := h.service.Actualizar(id, input)
	if err != nil {
		utils.ErrorResponse(c, mapErrorToStatus(err), err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Parto actualizado exitosamente", parto)
}

// GetEstadisticas godoc
// GET /api/partos/estadisticas?mes=1&anio=2026&granja_id=1
func (h *PartoHandler) GetEstadisticas(c *gin.Context) {
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
