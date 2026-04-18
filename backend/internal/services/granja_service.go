package services

import (
	"github.com/martin/sgp/internal/models"
	"github.com/martin/sgp/internal/repositories"
)

// GranjaService maneja la lógica de negocio de granjas
type GranjaService struct {
	repos *repositories.RepositoryContainer
}

// NewGranjaService crea una nueva instancia del servicio
func NewGranjaService(repos *repositories.RepositoryContainer) *GranjaService {
	return &GranjaService{repos: repos}
}

// --- DTOs ---

// CrearGranjaInput datos para crear una granja
type CrearGranjaInput struct {
	Nombre      string `json:"nombre" binding:"required"`
	Descripcion string `json:"descripcion"`
	Ubicacion   string `json:"ubicacion"`
}

// ActualizarGranjaInput datos para actualizar una granja
type ActualizarGranjaInput struct {
	Nombre      string `json:"nombre"`
	Descripcion string `json:"descripcion"`
	Ubicacion   string `json:"ubicacion"`
}

// AsignarUsuarioInput datos para asignar un usuario a una granja
type AsignarUsuarioInput struct {
	UsuarioID uint   `json:"usuario_id" binding:"required"`
	Rol       string `json:"rol" binding:"required,oneof=propietario administrador operador"`
}

// --- Métodos del servicio ---

// Crear registra una nueva granja
func (s *GranjaService) Crear(input CrearGranjaInput) (*models.Granja, error) {
	granja := &models.Granja{
		Nombre: input.Nombre,
		Activo: true,
	}

	if input.Descripcion != "" {
		granja.Descripcion = &input.Descripcion
	}
	if input.Ubicacion != "" {
		granja.Ubicacion = &input.Ubicacion
	}

	if err := s.repos.Granja.Create(granja); err != nil {
		return nil, err
	}

	return granja, nil
}

// ObtenerPorID obtiene una granja por ID
func (s *GranjaService) ObtenerPorID(id uint) (*models.Granja, error) {
	granja, err := s.repos.Granja.FindByID(id)
	if err != nil {
		return nil, ErrNotFound
	}
	return granja, nil
}

// ListarTodas lista todas las granjas, opcionalmente filtradas por activo
func (s *GranjaService) ListarTodas(activo *bool) ([]models.Granja, error) {
	return s.repos.Granja.FindAll(activo)
}

// ListarPorUsuario lista granjas a las que el usuario tiene acceso
func (s *GranjaService) ListarPorUsuario(usuarioID uint) ([]models.Granja, error) {
	return s.repos.Granja.FindByUsuarioID(usuarioID)
}

// Actualizar modifica una granja
func (s *GranjaService) Actualizar(id uint, input ActualizarGranjaInput) (*models.Granja, error) {
	granja, err := s.repos.Granja.FindByID(id)
	if err != nil {
		return nil, ErrNotFound
	}

	if input.Nombre != "" {
		granja.Nombre = input.Nombre
	}
	if input.Descripcion != "" {
		granja.Descripcion = &input.Descripcion
	}
	if input.Ubicacion != "" {
		granja.Ubicacion = &input.Ubicacion
	}

	if err := s.repos.Granja.Update(granja); err != nil {
		return nil, err
	}

	return granja, nil
}

// Eliminar da de baja una granja (solo si no tiene datos activos)
func (s *GranjaService) Eliminar(id uint) error {
	granja, err := s.repos.Granja.FindByID(id)
	if err != nil {
		return ErrNotFound
	}

	// Verificar que no tenga corrales activos
	corrales, err := s.repos.Corral.FindByGranjaID(id, boolPtr(true))
	if err != nil {
		return err
	}
	if len(corrales) > 0 {
		return ErrGranjaTieneDatosActivos
	}

	// Verificar que no tenga cerdas activas
	activas := true
	cerdas, err := s.repos.Cerda.FindByGranjaID(id, nil, &activas)
	if err != nil {
		return err
	}
	if len(cerdas) > 0 {
		return ErrGranjaTieneDatosActivos
	}

	// Verificar que no tenga padrillos activos
	padrillos, err := s.repos.Padrillo.FindByGranjaID(id, &activas)
	if err != nil {
		return err
	}
	if len(padrillos) > 0 {
		return ErrGranjaTieneDatosActivos
	}

	granja.Activo = false
	return s.repos.Granja.Update(granja)
}

// AsignarUsuario asigna un usuario a una granja con un rol
func (s *GranjaService) AsignarUsuario(granjaID uint, input AsignarUsuarioInput) error {
	// Verificar que la granja existe
	_, err := s.repos.Granja.FindByID(granjaID)
	if err != nil {
		return ErrNotFound
	}

	// Verificar que el usuario existe
	_, err = s.repos.Usuario.FindByID(input.UsuarioID)
	if err != nil {
		return ErrNotFound
	}

	return s.repos.Granja.AsignarUsuario(granjaID, input.UsuarioID, input.Rol)
}

// RemoverUsuario remueve un usuario de una granja
func (s *GranjaService) RemoverUsuario(granjaID, usuarioID uint) error {
	return s.repos.Granja.RemoverUsuario(granjaID, usuarioID)
}

// GetEstadisticas obtiene estadísticas generales de la granja
func (s *GranjaService) GetEstadisticas(granjaID uint) (map[string]interface{}, error) {
	_, err := s.repos.Granja.FindByID(granjaID)
	if err != nil {
		return nil, ErrNotFound
	}
	return s.repos.Granja.GetEstadisticas(granjaID)
}

// --- Helpers ---

func boolPtr(b bool) *bool {
	return &b
}
