package services

import (
	"time"

	"github.com/martin/sgp/internal/models"
	"github.com/martin/sgp/internal/repositories"
	"github.com/martin/sgp/internal/utils"
	"golang.org/x/crypto/bcrypt"
)

// AuthService maneja la lógica de autenticación
type AuthService struct {
	repos         *repositories.RepositoryContainer
	jwtSecret     string
	jwtExpiration time.Duration
}

// NewAuthService crea una nueva instancia del servicio
func NewAuthService(repos *repositories.RepositoryContainer, jwtSecret string, jwtExpiration time.Duration) *AuthService {
	return &AuthService{
		repos:         repos,
		jwtSecret:     jwtSecret,
		jwtExpiration: jwtExpiration,
	}
}

// --- DTOs ---

// RegistroInput datos para registrar un usuario
type RegistroInput struct {
	Username        string `json:"username" binding:"required,min=3,max=50"`
	Email           string `json:"email" binding:"required,email"`
	Password        string `json:"password" binding:"required,min=6"`
	NombreCompleto  string `json:"nombre_completo"`
	Establecimiento string `json:"establecimiento"`
}

// LoginInput datos para iniciar sesión
type LoginInput struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// LoginResponse respuesta del login
type LoginResponse struct {
	Token   string          `json:"token"`
	Usuario *models.Usuario `json:"usuario"`
}

// --- Métodos ---

// Registrar crea un nuevo usuario
func (s *AuthService) Registrar(input RegistroInput) (*models.Usuario, error) {
	// Verificar username duplicado
	existe, err := s.repos.Usuario.ExisteUsername(input.Username, nil)
	if err != nil {
		return nil, err
	}
	if existe {
		return nil, ErrDuplicateKey
	}

	// Verificar email duplicado
	existe, err = s.repos.Usuario.ExisteEmail(input.Email, nil)
	if err != nil {
		return nil, err
	}
	if existe {
		return nil, ErrDuplicateKey
	}

	// Hash de password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	usuario := &models.Usuario{
		Username:     input.Username,
		Email:        input.Email,
		PasswordHash: string(hashedPassword),
		Rol:          models.RolUsuario,
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

	// Limpiar password del response
	usuario.PasswordHash = ""

	return usuario, nil
}

// Login autentica un usuario y retorna un JWT
func (s *AuthService) Login(input LoginInput) (*LoginResponse, error) {
	// Buscar usuario por username
	usuario, err := s.repos.Usuario.FindByUsername(input.Username)
	if err != nil {
		return nil, ErrNotFound
	}

	if !usuario.Activo {
		return nil, ErrForbidden
	}

	// Verificar password
	if err := bcrypt.CompareHashAndPassword([]byte(usuario.PasswordHash), []byte(input.Password)); err != nil {
		return nil, ErrForbidden
	}

	// Generar JWT
	token, err := utils.GenerateJWT(usuario.ID, usuario.Email, usuario.Rol, s.jwtSecret, s.jwtExpiration)
	if err != nil {
		return nil, err
	}

	// Limpiar password del response
	usuario.PasswordHash = ""

	return &LoginResponse{
		Token:   token,
		Usuario: usuario,
	}, nil
}
