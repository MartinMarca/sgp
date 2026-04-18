package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/martin/sgp/internal/services"
	"github.com/martin/sgp/internal/utils"
)

// VentaHandler maneja los endpoints de ventas
type VentaHandler struct {
	service *services.VentaService
}

func NewVentaHandler(service *services.VentaService) *VentaHandler {
	return &VentaHandler{service: service}
}

// POST /api/ventas
func (h *VentaHandler) Crear(c *gin.Context) {
	var input services.CrearVentaInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	venta, err := h.service.Crear(input)
	if err != nil {
		utils.ErrorResponse(c, mapErrorToStatus(err), err.Error())
		return
	}
	utils.SuccessResponse(c, http.StatusCreated, "Venta registrada exitosamente", venta)
}

// GET /api/ventas/:id
func (h *VentaHandler) ObtenerPorID(c *gin.Context) {
	id, err := getIDParam(c, "id")
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "ID inválido")
		return
	}
	venta, err := h.service.ObtenerPorID(id)
	if err != nil {
		utils.ErrorResponse(c, mapErrorToStatus(err), err.Error())
		return
	}
	utils.SuccessResponse(c, http.StatusOK, "", venta)
}

// GET /api/ventas?granja_id=1&mes=3&anio=2026
func (h *VentaHandler) ListarPorPeriodo(c *gin.Context) {
	granjaID := getOptionalUintQuery(c, "granja_id")
	mes := getIntQuery(c, "mes", 0)
	anio := getIntQuery(c, "anio", 0)
	ventas, err := h.service.ListarPorPeriodo(granjaID, mes, anio)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	utils.SuccessResponse(c, http.StatusOK, "", ventas)
}

// PUT /api/ventas/:id
func (h *VentaHandler) Actualizar(c *gin.Context) {
	id, err := getIDParam(c, "id")
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "ID inválido")
		return
	}
	var input services.ActualizarVentaInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	venta, err := h.service.Actualizar(id, input)
	if err != nil {
		utils.ErrorResponse(c, mapErrorToStatus(err), err.Error())
		return
	}
	utils.SuccessResponse(c, http.StatusOK, "Venta actualizada exitosamente", venta)
}

// DELETE /api/ventas/:id
func (h *VentaHandler) Eliminar(c *gin.Context) {
	id, err := getIDParam(c, "id")
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "ID inválido")
		return
	}
	if err := h.service.Eliminar(id); err != nil {
		utils.ErrorResponse(c, mapErrorToStatus(err), err.Error())
		return
	}
	utils.SuccessResponse(c, http.StatusOK, "Venta eliminada exitosamente", nil)
}

// GET /api/ventas/estadisticas?granja_id=1&mes=3&anio=2026
func (h *VentaHandler) GetEstadisticas(c *gin.Context) {
	granjaID := getOptionalUintQuery(c, "granja_id")
	mes := getIntQuery(c, "mes", 0)
	anio := getIntQuery(c, "anio", 0)
	stats, err := h.service.GetEstadisticas(granjaID, mes, anio)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	utils.SuccessResponse(c, http.StatusOK, "", stats)
}
