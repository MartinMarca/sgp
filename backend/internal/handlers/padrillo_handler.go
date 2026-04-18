package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/martin/sgp/internal/services"
	"github.com/martin/sgp/internal/utils"
)

// PadrilloHandler maneja los endpoints de padrillos
type PadrilloHandler struct {
	service *services.PadrilloService
}

// NewPadrilloHandler crea una nueva instancia del handler
func NewPadrilloHandler(service *services.PadrilloService) *PadrilloHandler {
	return &PadrilloHandler{service: service}
}

// CrearEnGranja godoc
// POST /api/granjas/:id/padrillos
func (h *PadrilloHandler) CrearEnGranja(c *gin.Context) {
	granjaID, err := getIDParam(c, "id")
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "ID de granja inválido")
		return
	}

	var input services.CrearPadrilloInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	input.GranjaID = granjaID

	padrillo, err := h.service.Crear(input)
	if err != nil {
		utils.ErrorResponse(c, mapErrorToStatus(err), err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "Padrillo registrado exitosamente", padrillo)
}

// ObtenerPorID godoc
// GET /api/padrillos/:id
func (h *PadrilloHandler) ObtenerPorID(c *gin.Context) {
	id, err := getIDParam(c, "id")
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "ID inválido")
		return
	}

	padrillo, err := h.service.ObtenerPorID(id)
	if err != nil {
		utils.ErrorResponse(c, mapErrorToStatus(err), err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "", padrillo)
}

// ListarPorGranja godoc
// GET /api/granjas/:id/padrillos?activo=true
func (h *PadrilloHandler) ListarPorGranja(c *gin.Context) {
	granjaID, err := getIDParam(c, "id")
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "ID de granja inválido")
		return
	}

	activo := getOptionalBoolQuery(c, "activo")

	padrillos, err := h.service.ListarPorGranja(granjaID, activo)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "", padrillos)
}

// Actualizar godoc
// PUT /api/padrillos/:id
func (h *PadrilloHandler) Actualizar(c *gin.Context) {
	id, err := getIDParam(c, "id")
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "ID inválido")
		return
	}

	var input services.ActualizarPadrilloInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	padrillo, err := h.service.Actualizar(id, input)
	if err != nil {
		utils.ErrorResponse(c, mapErrorToStatus(err), err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Padrillo actualizado exitosamente", padrillo)
}

// DarDeBaja godoc
// POST /api/padrillos/:id/baja
func (h *PadrilloHandler) DarDeBaja(c *gin.Context) {
	id, err := getIDParam(c, "id")
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "ID inválido")
		return
	}

	var input services.BajaPadrilloInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	padrillo, err := h.service.DarDeBaja(id, input)
	if err != nil {
		utils.ErrorResponse(c, mapErrorToStatus(err), err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Padrillo dado de baja exitosamente", padrillo)
}

// GetEstadisticas godoc
// GET /api/padrillos/:id/estadisticas
func (h *PadrilloHandler) GetEstadisticas(c *gin.Context) {
	id, err := getIDParam(c, "id")
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "ID inválido")
		return
	}

	stats, err := h.service.GetEstadisticas(id)
	if err != nil {
		utils.ErrorResponse(c, mapErrorToStatus(err), err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "", stats)
}
