package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/martin/sgp/internal/services"
	"github.com/martin/sgp/internal/utils"
)

// MuerteLechonHandler maneja los endpoints de muertes de lechones
type MuerteLechonHandler struct {
	service *services.MuerteLechonService
}

// NewMuerteLechonHandler crea una nueva instancia del handler
func NewMuerteLechonHandler(service *services.MuerteLechonService) *MuerteLechonHandler {
	return &MuerteLechonHandler{service: service}
}

// Crear godoc
// POST /api/muertes-lechones
func (h *MuerteLechonHandler) Crear(c *gin.Context) {
	var input services.CrearMuerteLechonInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
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

	if err := h.service.Eliminar(id); err != nil {
		utils.ErrorResponse(c, mapErrorToStatus(err), err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Registro de muerte eliminado exitosamente", nil)
}

// GetEstadisticas godoc
// GET /api/muertes-lechones/estadisticas?mes=1&anio=2026&granja_id=1
func (h *MuerteLechonHandler) GetEstadisticas(c *gin.Context) {
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
