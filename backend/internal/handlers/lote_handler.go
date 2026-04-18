package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/martin/sgp/internal/services"
	"github.com/martin/sgp/internal/utils"
)

// LoteHandler maneja los endpoints de lotes
type LoteHandler struct {
	service *services.LoteService
}

// NewLoteHandler crea una nueva instancia del handler
func NewLoteHandler(service *services.LoteService) *LoteHandler {
	return &LoteHandler{service: service}
}

// CrearEnCorral godoc
// POST /api/corrales/:id/lotes
func (h *LoteHandler) CrearEnCorral(c *gin.Context) {
	corralID, err := getIDParam(c, "id")
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "ID de corral inválido")
		return
	}

	var input services.CrearLoteInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	input.CorralID = corralID

	lote, err := h.service.Crear(input)
	if err != nil {
		utils.ErrorResponse(c, mapErrorToStatus(err), err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "Lote creado exitosamente", lote)
}

// ObtenerPorID godoc
// GET /api/lotes/:id
func (h *LoteHandler) ObtenerPorID(c *gin.Context) {
	id, err := getIDParam(c, "id")
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "ID inválido")
		return
	}

	lote, err := h.service.ObtenerPorID(id)
	if err != nil {
		utils.ErrorResponse(c, mapErrorToStatus(err), err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "", lote)
}

// ListarPorCorral godoc
// GET /api/corrales/:id/lotes
func (h *LoteHandler) ListarPorCorral(c *gin.Context) {
	corralID, err := getIDParam(c, "id")
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "ID de corral inválido")
		return
	}

	estado := getOptionalStringQuery(c, "estado")

	lotes, err := h.service.ListarPorCorral(corralID, estado)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "", lotes)
}

// ListarPorEstado godoc
// GET /api/lotes?estado=activo
func (h *LoteHandler) ListarPorEstado(c *gin.Context) {
	estado := c.DefaultQuery("estado", "activo")

	lotes, err := h.service.ListarPorEstado(estado)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "", lotes)
}

// Actualizar godoc
// PUT /api/lotes/:id
func (h *LoteHandler) Actualizar(c *gin.Context) {
	id, err := getIDParam(c, "id")
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "ID inválido")
		return
	}

	var input services.ActualizarLoteInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	lote, err := h.service.Actualizar(id, input)
	if err != nil {
		utils.ErrorResponse(c, mapErrorToStatus(err), err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Lote actualizado exitosamente", lote)
}

// Cerrar godoc
// POST /api/lotes/:id/cerrar
func (h *LoteHandler) Cerrar(c *gin.Context) {
	id, err := getIDParam(c, "id")
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "ID inválido")
		return
	}

	var input services.CerrarLoteInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	lote, err := h.service.Cerrar(id, input)
	if err != nil {
		utils.ErrorResponse(c, mapErrorToStatus(err), err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Lote cerrado exitosamente", lote)
}

// GetDestetes godoc
// GET /api/lotes/:id/destetes
func (h *LoteHandler) GetDestetes(c *gin.Context) {
	id, err := getIDParam(c, "id")
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "ID inválido")
		return
	}

	destetes, err := h.service.GetDestetes(id)
	if err != nil {
		utils.ErrorResponse(c, mapErrorToStatus(err), err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "", destetes)
}
