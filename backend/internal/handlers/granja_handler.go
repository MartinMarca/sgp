package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/martin/sgp/internal/services"
	"github.com/martin/sgp/internal/utils"
)

// GranjaHandler maneja los endpoints de granjas
type GranjaHandler struct {
	service *services.GranjaService
}

// NewGranjaHandler crea una nueva instancia del handler
func NewGranjaHandler(service *services.GranjaService) *GranjaHandler {
	return &GranjaHandler{service: service}
}

// Crear godoc
// POST /api/granjas
func (h *GranjaHandler) Crear(c *gin.Context) {
	var input services.CrearGranjaInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	granja, err := h.service.Crear(input)
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
	activo := getOptionalBoolQuery(c, "activo")

	granjas, err := h.service.ListarTodas(activo)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "", granjas)
}

// ListarPorUsuario godoc
// GET /api/granjas/mis-granjas
func (h *GranjaHandler) ListarPorUsuario(c *gin.Context) {
	userID := getUserID(c)

	granjas, err := h.service.ListarPorUsuario(userID)
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

	if err := h.service.Eliminar(id); err != nil {
		utils.ErrorResponse(c, mapErrorToStatus(err), err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Granja dada de baja exitosamente", nil)
}

// AsignarUsuario godoc
// POST /api/granjas/:id/usuarios
func (h *GranjaHandler) AsignarUsuario(c *gin.Context) {
	id, err := getIDParam(c, "id")
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "ID inválido")
		return
	}

	var input services.AsignarUsuarioInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.service.AsignarUsuario(id, input); err != nil {
		utils.ErrorResponse(c, mapErrorToStatus(err), err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Usuario asignado exitosamente", nil)
}

// RemoverUsuario godoc
// DELETE /api/granjas/:id/usuarios/:usuario_id
func (h *GranjaHandler) RemoverUsuario(c *gin.Context) {
	granjaID, err := getIDParam(c, "id")
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "ID de granja inválido")
		return
	}

	usuarioID, err := getIDParam(c, "usuario_id")
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "ID de usuario inválido")
		return
	}

	if err := h.service.RemoverUsuario(granjaID, usuarioID); err != nil {
		utils.ErrorResponse(c, mapErrorToStatus(err), err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Usuario removido exitosamente", nil)
}

// GetEstadisticas godoc
// GET /api/granjas/:id/estadisticas
func (h *GranjaHandler) GetEstadisticas(c *gin.Context) {
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
