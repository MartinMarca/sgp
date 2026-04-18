package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/martin/sgp/internal/services"
	"github.com/martin/sgp/internal/utils"
)

// CorralHandler maneja los endpoints de corrales
type CorralHandler struct {
	service *services.CorralService
}

// NewCorralHandler crea una nueva instancia del handler
func NewCorralHandler(service *services.CorralService) *CorralHandler {
	return &CorralHandler{service: service}
}

// CrearEnGranja godoc
// POST /api/granjas/:id/corrales
func (h *CorralHandler) CrearEnGranja(c *gin.Context) {
	granjaID, err := getIDParam(c, "id")
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "ID de granja inválido")
		return
	}

	var input services.CrearCorralInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	input.GranjaID = granjaID

	corral, err := h.service.Crear(input)
	if err != nil {
		utils.ErrorResponse(c, mapErrorToStatus(err), err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "Corral creado exitosamente", corral)
}

// ObtenerPorID godoc
// GET /api/corrales/:id
func (h *CorralHandler) ObtenerPorID(c *gin.Context) {
	id, err := getIDParam(c, "id")
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "ID inválido")
		return
	}

	corral, err := h.service.ObtenerPorID(id)
	if err != nil {
		utils.ErrorResponse(c, mapErrorToStatus(err), err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "", corral)
}

// ListarPorGranja godoc
// GET /api/granjas/:id/corrales
func (h *CorralHandler) ListarPorGranja(c *gin.Context) {
	granjaID, err := getIDParam(c, "id")
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "ID de granja inválido")
		return
	}

	activo := getOptionalBoolQuery(c, "activo")

	corrales, err := h.service.ListarPorGranja(granjaID, activo)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "", corrales)
}

// Actualizar godoc
// PUT /api/corrales/:id
func (h *CorralHandler) Actualizar(c *gin.Context) {
	id, err := getIDParam(c, "id")
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "ID inválido")
		return
	}

	var input services.ActualizarCorralInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	corral, err := h.service.Actualizar(id, input)
	if err != nil {
		utils.ErrorResponse(c, mapErrorToStatus(err), err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Corral actualizado exitosamente", corral)
}

// Eliminar godoc
// DELETE /api/corrales/:id
func (h *CorralHandler) Eliminar(c *gin.Context) {
	id, err := getIDParam(c, "id")
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "ID inválido")
		return
	}

	if err := h.service.Eliminar(id); err != nil {
		utils.ErrorResponse(c, mapErrorToStatus(err), err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Corral eliminado exitosamente", nil)
}

// GetOcupacion godoc
// GET /api/corrales/:id/ocupacion
func (h *CorralHandler) GetOcupacion(c *gin.Context) {
	id, err := getIDParam(c, "id")
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "ID inválido")
		return
	}

	ocupacion, err := h.service.GetOcupacion(id)
	if err != nil {
		utils.ErrorResponse(c, mapErrorToStatus(err), err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "", gin.H{"ocupacion": ocupacion})
}
