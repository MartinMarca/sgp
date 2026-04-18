package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/martin/sgp/internal/services"
)

// getIDParam extrae un parámetro uint de la URL
func getIDParam(c *gin.Context, param string) (uint, error) {
	idStr := c.Param(param)
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return 0, err
	}
	return uint(id), nil
}

// getUserID obtiene el user_id del contexto (set por middleware auth)
func getUserID(c *gin.Context) uint {
	userID, _ := c.Get("user_id")
	if id, ok := userID.(uint); ok {
		return id
	}
	return 0
}

// getOptionalIntQuery extrae un query param int opcional
func getOptionalIntQuery(c *gin.Context, param string) *int {
	val := c.Query(param)
	if val == "" {
		return nil
	}
	i, err := strconv.Atoi(val)
	if err != nil {
		return nil
	}
	return &i
}

// getOptionalBoolQuery extrae un query param bool opcional
func getOptionalBoolQuery(c *gin.Context, param string) *bool {
	val := c.Query(param)
	if val == "" {
		return nil
	}
	b, err := strconv.ParseBool(val)
	if err != nil {
		return nil
	}
	return &b
}

// getOptionalStringQuery extrae un query param string opcional
func getOptionalStringQuery(c *gin.Context, param string) *string {
	val := c.Query(param)
	if val == "" {
		return nil
	}
	return &val
}

// getIntQuery extrae un query param int con valor default
func getIntQuery(c *gin.Context, param string, defaultVal int) int {
	val := c.Query(param)
	if val == "" {
		return defaultVal
	}
	i, err := strconv.Atoi(val)
	if err != nil {
		return defaultVal
	}
	return i
}

// getOptionalUintQuery extrae un query param uint opcional
func getOptionalUintQuery(c *gin.Context, param string) *uint {
	val := c.Query(param)
	if val == "" {
		return nil
	}
	i, err := strconv.ParseUint(val, 10, 32)
	if err != nil {
		return nil
	}
	u := uint(i)
	return &u
}

// mapErrorToStatus convierte errores de servicio a status HTTP
func mapErrorToStatus(err error) int {
	switch {
	case errors.Is(err, services.ErrNotFound):
		return http.StatusNotFound
	case errors.Is(err, services.ErrDuplicateKey),
		errors.Is(err, services.ErrCaravanaDuplicada):
		return http.StatusConflict
	case errors.Is(err, services.ErrForbidden):
		return http.StatusForbidden
	case errors.Is(err, services.ErrCerdaNoDisponible),
		errors.Is(err, services.ErrCerdaNoEnServicio),
		errors.Is(err, services.ErrCerdaNoEnGestacion),
		errors.Is(err, services.ErrCerdaNoEnCria),
		errors.Is(err, services.ErrCerdaNoActiva),
		errors.Is(err, services.ErrCerdaTieneServicioActivo),
		errors.Is(err, services.ErrPrenezYaConfirmada),
		errors.Is(err, services.ErrPrenezYaCancelada),
		errors.Is(err, services.ErrPrenezNoConfirmada),
		errors.Is(err, services.ErrMotivoRequerido),
		errors.Is(err, services.ErrNoHayServicioConfirmado),
		errors.Is(err, services.ErrNoHayPartoSinDestete),
		errors.Is(err, services.ErrLechonesInvalidos),
		errors.Is(err, services.ErrTotalMenorQueVivos),
		errors.Is(err, services.ErrDestetadosExcedenVivos),
		errors.Is(err, services.ErrLoteRequerido),
		errors.Is(err, services.ErrLoteNoActivo),
		errors.Is(err, services.ErrCorralRequerido),
		errors.Is(err, services.ErrCorralTieneLotesActivos),
		errors.Is(err, services.ErrGranjaTieneDatosActivos),
		errors.Is(err, services.ErrServicioRequierePadrillo),
		errors.Is(err, services.ErrServicioRequierePajuela),
		errors.Is(err, services.ErrMuerteRequierePartoOLote),
		errors.Is(err, services.ErrMuertesExcedenVivos),
		errors.Is(err, services.ErrMuertesExcedenLote),
		errors.Is(err, services.ErrCausaMuerteInvalida),
		errors.Is(err, services.ErrVentaExcedeLote),
		errors.Is(err, services.ErrTipoAnimalInvalido):
		return http.StatusUnprocessableEntity
	default:
		return http.StatusInternalServerError
	}
}
