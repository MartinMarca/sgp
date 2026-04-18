package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/martin/sgp/internal/services"
	"github.com/martin/sgp/internal/utils"
)

// DesteteHandler maneja los endpoints de destetes
type DesteteHandler struct {
	service *services.DesteteService
}

// NewDesteteHandler crea una nueva instancia del handler
func NewDesteteHandler(service *services.DesteteService) *DesteteHandler {
	return &DesteteHandler{service: service}
}

// Crear godoc
// POST /api/destetes
func (h *DesteteHandler) Crear(c *gin.Context) {
	var input services.CrearDesteteInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	destete, err := h.service.Crear(input)
	if err != nil {
		utils.ErrorResponse(c, mapErrorToStatus(err), err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "Destete registrado exitosamente", destete)
}

// ObtenerPorID godoc
// GET /api/destetes/:id
func (h *DesteteHandler) ObtenerPorID(c *gin.Context) {
	id, err := getIDParam(c, "id")
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "ID inválido")
		return
	}

	destete, err := h.service.ObtenerPorID(id)
	if err != nil {
		utils.ErrorResponse(c, mapErrorToStatus(err), err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "", destete)
}

// ListarPorCerda godoc
// GET /api/cerdas/:id/destetes
func (h *DesteteHandler) ListarPorCerda(c *gin.Context) {
	cerdaID, err := getIDParam(c, "id")
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "ID de cerda inválido")
		return
	}

	destetes, err := h.service.ListarPorCerda(cerdaID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "", destetes)
}

// ListarPorLote godoc
// GET /api/lotes/:id/destetes
func (h *DesteteHandler) ListarPorLote(c *gin.Context) {
	loteID, err := getIDParam(c, "id")
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "ID de lote inválido")
		return
	}

	destetes, err := h.service.ListarPorLote(loteID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "", destetes)
}

// ListarPorPeriodo godoc
// GET /api/destetes?mes=1&anio=2026&granja_id=1
func (h *DesteteHandler) ListarPorPeriodo(c *gin.Context) {
	mes := getIntQuery(c, "mes", 0)
	anio := getIntQuery(c, "anio", 0)
	granjaID := getOptionalUintQuery(c, "granja_id")

	destetes, err := h.service.ListarPorPeriodo(granjaID, mes, anio)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "", destetes)
}

// Actualizar godoc
// PUT /api/destetes/:id
func (h *DesteteHandler) Actualizar(c *gin.Context) {
	id, err := getIDParam(c, "id")
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "ID inválido")
		return
	}

	var input services.ActualizarDesteteInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	destete, err := h.service.Actualizar(id, input)
	if err != nil {
		utils.ErrorResponse(c, mapErrorToStatus(err), err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Destete actualizado exitosamente", destete)
}

// GetEstadisticas godoc
// GET /api/destetes/estadisticas?mes=1&anio=2026&granja_id=1
func (h *DesteteHandler) GetEstadisticas(c *gin.Context) {
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
