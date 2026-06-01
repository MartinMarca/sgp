package services

import (
	"github.com/martin/sgp/internal/authz"
	"github.com/martin/sgp/internal/models"
	"github.com/martin/sgp/internal/repositories"
	"golang.org/x/crypto/bcrypt"
)

// UsuarioService maneja la lógica de negocio de usuarios.
type UsuarioService struct {
	repos *repositories.RepositoryContainer
}

// NewUsuarioService crea una nueva instancia del servicio.
func NewUsuarioService(repos *repositories.RepositoryContainer) *UsuarioService {
	return &UsuarioService{repos: repos}
}

// CrearUsuarioInput datos para crear un usuario.
type CrearUsuarioInput struct {
	Username        string `json:"username" binding:"required,min=3,max=50"`
	Email           string `json:"email" binding:"required,email"`
	Password        string `json:"password" binding:"required,min=6"`
	NombreCompleto  string `json:"nombre_completo"`
	Establecimiento string `json:"establecimiento"`
	Rol             string `json:"rol" binding:"required,oneof=admin propietario empleado"`
	PropietarioID   *uint  `json:"propietario_id"`
}

// ActualizarUsuarioInput datos para actualizar un usuario.
type ActualizarUsuarioInput struct {
	Email           string `json:"email"`
	NombreCompleto  string `json:"nombre_completo"`
	Establecimiento string `json:"establecimiento"`
	Rol             string `json:"rol"`
	PropietarioID   *uint  `json:"propietario_id"`
	Activo          *bool  `json:"activo"`
	Password        string `json:"password"`
}

// Listar devuelve usuarios según el alcance del actor.
func (s *UsuarioService) Listar(actor authz.Actor, activo *bool) ([]models.Usuario, error) {
	var usuarios []models.Usuario
	var err error

	switch {
	case actor.IsAdmin():
		usuarios, err = s.repos.Usuario.FindAll(activo)
	case actor.Rol == models.RolPropietario:
		usuarios, err = s.repos.Usuario.FindByPropietarioID(actor.ID, activo)
	default:
		return nil, ErrForbidden
	}
	if err != nil {
		return nil, err
	}

	for i := range usuarios {
		usuarios[i].PasswordHash = ""
	}
	return usuarios, nil
}

// ObtenerPorID obtiene un usuario por ID.
func (s *UsuarioService) ObtenerPorID(id uint) (*models.Usuario, error) {
	usuario, err := s.repos.Usuario.FindByID(id)
	if err != nil {
		return nil, ErrNotFound
	}
	usuario.PasswordHash = ""
	return usuario, nil
}

// Crear registra un nuevo usuario (validaciones de negocio; authz en handler).
func (s *UsuarioService) Crear(actor authz.Actor, input CrearUsuarioInput) (*models.Usuario, error) {
	propietarioID, err := s.resolvePropietarioIDForCreate(actor, input)
	if err != nil {
		return nil, err
	}

	existe, err := s.repos.Usuario.ExisteUsername(input.Username, nil)
	if err != nil {
		return nil, err
	}
	if existe {
		return nil, ErrDuplicateKey
	}

	existe, err = s.repos.Usuario.ExisteEmail(input.Email, nil)
	if err != nil {
		return nil, err
	}
	if existe {
		return nil, ErrDuplicateKey
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	usuario := &models.Usuario{
		Username:     input.Username,
		Email:        input.Email,
		PasswordHash: string(hashedPassword),
		Rol:          input.Rol,
		PropietarioID: propietarioID,
		Activo:       true,
	}

	if input.NombreCompleto != "" {
		usuario.NombreCompleto = &input.NombreCompleto
	}
	if input.Establecimiento != "" {
		usuario.Establecimiento = &input.Establecimiento
	}

	if err := s.repos.Usuario.Create(usuario); err != nil {
		return nil, err
	}

	usuario.PasswordHash = ""
	return usuario, nil
}

// Actualizar modifica un usuario existente.
func (s *UsuarioService) Actualizar(actor authz.Actor, id uint, input ActualizarUsuarioInput) (*models.Usuario, error) {
	usuario, err := s.repos.Usuario.FindByID(id)
	if err != nil {
		return nil, ErrNotFound
	}

	if input.Email != "" && input.Email != usuario.Email {
		existe, err := s.repos.Usuario.ExisteEmail(input.Email, &id)
		if err != nil {
			return nil, err
		}
		if existe {
			return nil, ErrDuplicateKey
		}
		usuario.Email = input.Email
	}

	if input.NombreCompleto != "" {
		usuario.NombreCompleto = &input.NombreCompleto
	}
	if input.Establecimiento != "" {
		usuario.Establecimiento = &input.Establecimiento
	}

	if input.Password != "" {
		hashed, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
		if err != nil {
			return nil, err
		}
		usuario.PasswordHash = string(hashed)
	}

	if err := s.applyRoleAndActivo(actor, usuario, input); err != nil {
		return nil, err
	}

	if err := s.repos.Usuario.Update(usuario); err != nil {
		return nil, err
	}

	usuario.PasswordHash = ""
	return usuario, nil
}

// Desactivar marca un usuario como inactivo (soft delete lógico).
func (s *UsuarioService) Desactivar(id uint) error {
	usuario, err := s.repos.Usuario.FindByID(id)
	if err != nil {
		return ErrNotFound
	}

	if usuario.Rol == models.RolAdmin {
		activo := true
		count, err := s.repos.Usuario.CountByRol(models.RolAdmin, &activo)
		if err != nil {
			return err
		}
		if count <= 1 {
			return ErrUltimoAdmin
		}
	}

	usuario.Activo = false
	return s.repos.Usuario.Update(usuario)
}

func (s *UsuarioService) resolvePropietarioIDForCreate(actor authz.Actor, input CrearUsuarioInput) (*uint, error) {
	switch input.Rol {
	case models.RolAdmin, models.RolPropietario:
		return nil, nil
	case models.RolEmpleado:
		if actor.Rol == models.RolPropietario {
			id := actor.ID
			return &id, nil
		}
		if actor.IsAdmin() {
			if input.PropietarioID == nil || *input.PropietarioID == 0 {
				return nil, ErrPropietarioRequerido
			}
			prop, err := s.repos.Usuario.FindByID(*input.PropietarioID)
			if err != nil {
				return nil, ErrNotFound
			}
			if prop.Rol != models.RolPropietario {
				return nil, ErrPropietarioRequerido
			}
			return input.PropietarioID, nil
		}
		return nil, ErrForbidden
	default:
		return nil, ErrRolInvalido
	}
}

func (s *UsuarioService) applyRoleAndActivo(actor authz.Actor, usuario *models.Usuario, input ActualizarUsuarioInput) error {
	isSelf := actor.ID == usuario.ID

	if actor.IsAdmin() {
		if input.Rol != "" {
			if input.Rol != models.RolAdmin && input.Rol != models.RolPropietario && input.Rol != models.RolEmpleado {
				return ErrRolInvalido
			}
			usuario.Rol = input.Rol
			switch input.Rol {
			case models.RolEmpleado:
				if input.PropietarioID == nil || *input.PropietarioID == 0 {
					return ErrPropietarioRequerido
				}
				prop, err := s.repos.Usuario.FindByID(*input.PropietarioID)
				if err != nil || prop.Rol != models.RolPropietario {
					return ErrPropietarioRequerido
				}
				usuario.PropietarioID = input.PropietarioID
			default:
				usuario.PropietarioID = nil
			}
		}
		if input.Activo != nil {
			if !*input.Activo && usuario.Rol == models.RolAdmin {
				activo := true
				count, err := s.repos.Usuario.CountByRol(models.RolAdmin, &activo)
				if err != nil {
					return err
				}
				if count <= 1 {
					return ErrUltimoAdmin
				}
			}
			usuario.Activo = *input.Activo
		}
		return nil
	}

	if actor.Rol == models.RolPropietario && !isSelf {
		if input.Rol != "" {
			if input.Rol != models.RolEmpleado {
				return ErrRolInvalido
			}
			usuario.Rol = models.RolEmpleado
			id := actor.ID
			usuario.PropietarioID = &id
		}
		if input.Activo != nil {
			usuario.Activo = *input.Activo
		}
		return nil
	}

	if isSelf && actor.Rol == models.RolEmpleado {
		return nil
	}

	if input.Rol != "" || input.Activo != nil {
		return ErrForbidden
	}
	return nil
}
