package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/martin/sgp/internal/authz"
	"github.com/martin/sgp/internal/repositories"
	"github.com/martin/sgp/internal/services"
	"github.com/martin/sgp/internal/utils"
)

// UsuarioHandler maneja los endpoints de usuarios.
type UsuarioHandler struct {
	service *services.UsuarioService
	authz   *authz.Service
	repos   *repositories.RepositoryContainer
}

// NewUsuarioHandler crea una nueva instancia del handler.
func NewUsuarioHandler(service *services.UsuarioService, authzSvc *authz.Service, repos *repositories.RepositoryContainer) *UsuarioHandler {
	return &UsuarioHandler{service: service, authz: authzSvc, repos: repos}
}

// Listar godoc
// GET /api/usuarios
func (h *UsuarioHandler) Listar(c *gin.Context) {
	if !requirePerm(c, h.authz, h.repos, authz.PermUsersManage) {
		return
	}

	actor, err := getActor(c, h.repos)
	if err != nil {
		utils.ErrorResponse(c, mapErrorToStatus(err), err.Error())
		return
	}

	activo := getOptionalBoolQuery(c, "activo")
	usuarios, err := h.service.Listar(actor, activo)
	if err != nil {
		utils.ErrorResponse(c, mapErrorToStatus(err), err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "", usuarios)
}

// Crear godoc
// POST /api/usuarios
func (h *UsuarioHandler) Crear(c *gin.Context) {
	var input services.CrearUsuarioInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	actor, err := getActor(c, h.repos)
	if err != nil {
		utils.ErrorResponse(c, mapErrorToStatus(err), err.Error())
		return
	}

	if err := h.authz.ValidateCreateRole(actor, input.Rol); err != nil {
		utils.ErrorResponse(c, mapErrorToStatus(err), err.Error())
		return
	}

	usuario, err := h.service.Crear(actor, input)
	if err != nil {
		utils.ErrorResponse(c, mapErrorToStatus(err), err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "Usuario creado exitosamente", usuario)
}

// ObtenerPorID godoc
// GET /api/usuarios/:id
func (h *UsuarioHandler) ObtenerPorID(c *gin.Context) {
	id, err := getIDParam(c, "id")
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "ID inválido")
		return
	}

	actor, err := getActor(c, h.repos)
	if err != nil {
		utils.ErrorResponse(c, mapErrorToStatus(err), err.Error())
		return
	}

	usuario, err := h.service.ObtenerPorID(id)
	if err != nil {
		utils.ErrorResponse(c, mapErrorToStatus(err), err.Error())
		return
	}

	if err := h.authz.CanViewUser(actor, usuario); err != nil {
		utils.ErrorResponse(c, mapErrorToStatus(err), err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "", usuario)
}

// Actualizar godoc
// PUT /api/usuarios/:id
func (h *UsuarioHandler) Actualizar(c *gin.Context) {
	id, err := getIDParam(c, "id")
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "ID inválido")
		return
	}

	actor, err := getActor(c, h.repos)
	if err != nil {
		utils.ErrorResponse(c, mapErrorToStatus(err), err.Error())
		return
	}

	target, err := h.service.ObtenerPorID(id)
	if err != nil {
		utils.ErrorResponse(c, mapErrorToStatus(err), err.Error())
		return
	}

	if err := h.authz.CanModifyUser(actor, target); err != nil {
		utils.ErrorResponse(c, mapErrorToStatus(err), err.Error())
		return
	}

	var input services.ActualizarUsuarioInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	usuario, err := h.service.Actualizar(actor, id, input)
	if err != nil {
		utils.ErrorResponse(c, mapErrorToStatus(err), err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Usuario actualizado exitosamente", usuario)
}

// Eliminar godoc
// DELETE /api/usuarios/:id
func (h *UsuarioHandler) Eliminar(c *gin.Context) {
	id, err := getIDParam(c, "id")
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "ID inválido")
		return
	}

	actor, err := getActor(c, h.repos)
	if err != nil {
		utils.ErrorResponse(c, mapErrorToStatus(err), err.Error())
		return
	}

	target, err := h.service.ObtenerPorID(id)
	if err != nil {
		utils.ErrorResponse(c, mapErrorToStatus(err), err.Error())
		return
	}

	if err := h.authz.CanDeleteUser(actor, target); err != nil {
		utils.ErrorResponse(c, mapErrorToStatus(err), err.Error())
		return
	}

	if err := h.service.Desactivar(id); err != nil {
		utils.ErrorResponse(c, mapErrorToStatus(err), err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Usuario desactivado exitosamente", nil)
}
