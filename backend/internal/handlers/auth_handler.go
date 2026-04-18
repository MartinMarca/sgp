package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/martin/sgp/internal/services"
	"github.com/martin/sgp/internal/utils"
)

// AuthHandler maneja los endpoints de autenticación
type AuthHandler struct {
	service *services.AuthService
}

// NewAuthHandler crea una nueva instancia del handler
func NewAuthHandler(service *services.AuthService) *AuthHandler {
	return &AuthHandler{service: service}
}

// Registrar godoc
// POST /api/auth/register
func (h *AuthHandler) Registrar(c *gin.Context) {
	var input services.RegistroInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	usuario, err := h.service.Registrar(input)
	if err != nil {
		status := mapErrorToStatus(err)
		utils.ErrorResponse(c, status, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "Usuario registrado exitosamente", usuario)
}

// Login godoc
// POST /api/auth/login
func (h *AuthHandler) Login(c *gin.Context) {
	var input services.LoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	response, err := h.service.Login(input)
	if err != nil {
		status := mapErrorToStatus(err)
		utils.ErrorResponse(c, status, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Login exitoso", response)
}
