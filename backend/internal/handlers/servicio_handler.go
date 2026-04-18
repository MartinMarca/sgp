package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/martin/sgp/internal/services"
	"github.com/martin/sgp/internal/utils"
)

// ServicioHandler maneja los endpoints de servicios
type ServicioHandler struct {
	service *services.ServicioService
}

// NewServicioHandler crea una nueva instancia del handler
func NewServicioHandler(service *services.ServicioService) *ServicioHandler {
	return &ServicioHandler{service: service}
}

// Crear godoc
// POST /api/servicios
func (h *ServicioHandler) Crear(c *gin.Context) {
	var input services.CrearServicioInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	servicio, err := h.service.Crear(input)
	if err != nil {
		utils.ErrorResponse(c, mapErrorToStatus(err), err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "Servicio registrado exitosamente", servicio)
}

// ObtenerPorID godoc
// GET /api/servicios/:id
func (h *ServicioHandler) ObtenerPorID(c *gin.Context) {
	id, err := getIDParam(c, "id")
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "ID inválido")
		return
	}

	servicio, err := h.service.ObtenerPorID(id)
	if err != nil {
		utils.ErrorResponse(c, mapErrorToStatus(err), err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "", servicio)
}

// ListarPorCerda godoc
// GET /api/cerdas/:id/servicios
func (h *ServicioHandler) ListarPorCerda(c *gin.Context) {
	cerdaID, err := getIDParam(c, "id")
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "ID de cerda inválido")
		return
	}

	servicios, err := h.service.ListarPorCerda(cerdaID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "", servicios)
}

// ListarPorPeriodo godoc
// GET /api/servicios?mes=1&anio=2026&granja_id=1
func (h *ServicioHandler) ListarPorPeriodo(c *gin.Context) {
	mes := getIntQuery(c, "mes", 0)
	anio := getIntQuery(c, "anio", 0)
	granjaID := getOptionalUintQuery(c, "granja_id")

	servicios, err := h.service.ListarPorPeriodo(granjaID, mes, anio)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "", servicios)
}

// ListarPendientesConfirmacion godoc
// GET /api/servicios/pendientes?granja_id=1
func (h *ServicioHandler) ListarPendientesConfirmacion(c *gin.Context) {
	granjaID := getOptionalUintQuery(c, "granja_id")

	servicios, err := h.service.ListarPendientesConfirmacion(granjaID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "", servicios)
}

// Actualizar godoc
// PUT /api/servicios/:id
func (h *ServicioHandler) Actualizar(c *gin.Context) {
	id, err := getIDParam(c, "id")
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "ID inválido")
		return
	}

	var input services.ActualizarServicioInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	servicio, err := h.service.Actualizar(id, input)
	if err != nil {
		utils.ErrorResponse(c, mapErrorToStatus(err), err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Servicio actualizado exitosamente", servicio)
}

// ConfirmarPrenez godoc
// POST /api/servicios/:id/confirmar-prenez
func (h *ServicioHandler) ConfirmarPrenez(c *gin.Context) {
	id, err := getIDParam(c, "id")
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "ID inválido")
		return
	}

	var input services.ConfirmarPrenezInput
	// El body es opcional (puede confirmar con fecha actual)
	_ = c.ShouldBindJSON(&input)

	servicio, err := h.service.ConfirmarPrenez(id, input)
	if err != nil {
		utils.ErrorResponse(c, mapErrorToStatus(err), err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Preñez confirmada exitosamente", servicio)
}

// CancelarPrenez godoc
// POST /api/servicios/:id/cancelar-prenez
func (h *ServicioHandler) CancelarPrenez(c *gin.Context) {
	id, err := getIDParam(c, "id")
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "ID inválido")
		return
	}

	var input services.CancelarPrenezInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	servicio, err := h.service.CancelarPrenez(id, input)
	if err != nil {
		utils.ErrorResponse(c, mapErrorToStatus(err), err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Preñez cancelada exitosamente", servicio)
}
