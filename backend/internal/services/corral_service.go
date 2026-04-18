package services

import (
	"github.com/martin/sgp/internal/models"
	"github.com/martin/sgp/internal/repositories"
)

// CorralService maneja la lógica de negocio de corrales
type CorralService struct {
	repos *repositories.RepositoryContainer
}

// NewCorralService crea una nueva instancia del servicio
func NewCorralService(repos *repositories.RepositoryContainer) *CorralService {
	return &CorralService{repos: repos}
}

// --- DTOs ---

// CrearCorralInput datos para crear un corral
type CrearCorralInput struct {
	GranjaID         uint   `json:"granja_id"`
	Nombre           string `json:"nombre" binding:"required"`
	Descripcion      string `json:"descripcion"`
	CapacidadMaxima  *int   `json:"capacidad_maxima"`
}

// ActualizarCorralInput datos para actualizar un corral
type ActualizarCorralInput struct {
	Nombre          string `json:"nombre"`
	Descripcion     string `json:"descripcion"`
	CapacidadMaxima *int   `json:"capacidad_maxima"`
}

// --- Métodos del servicio ---

// Crear registra un nuevo corral
func (s *CorralService) Crear(input CrearCorralInput) (*models.Corral, error) {
	// Verificar que la granja existe
	_, err := s.repos.Granja.FindByID(input.GranjaID)
	if err != nil {
		return nil, ErrNotFound
	}

	corral := &models.Corral{
		GranjaID: input.GranjaID,
		Nombre:   input.Nombre,
		Activo:   true,
	}

	if input.Descripcion != "" {
		corral.Descripcion = &input.Descripcion
	}
	if input.CapacidadMaxima != nil {
		corral.CapacidadMaxima = input.CapacidadMaxima
	}

	if err := s.repos.Corral.Create(corral); err != nil {
		return nil, err
	}

	return corral, nil
}

// ObtenerPorID obtiene un corral por ID
func (s *CorralService) ObtenerPorID(id uint) (*models.Corral, error) {
	corral, err := s.repos.Corral.FindByID(id, "Granja", "Lotes")
	if err != nil {
		return nil, ErrNotFound
	}
	return corral, nil
}

// ListarPorGranja lista corrales por granja
func (s *CorralService) ListarPorGranja(granjaID uint, activo *bool) ([]models.Corral, error) {
	return s.repos.Corral.FindByGranjaID(granjaID, activo)
}

// Actualizar modifica un corral
func (s *CorralService) Actualizar(id uint, input ActualizarCorralInput) (*models.Corral, error) {
	corral, err := s.repos.Corral.FindByID(id)
	if err != nil {
		return nil, ErrNotFound
	}

	if input.Nombre != "" {
		corral.Nombre = input.Nombre
	}
	if input.Descripcion != "" {
		corral.Descripcion = &input.Descripcion
	}
	if input.CapacidadMaxima != nil {
		corral.CapacidadMaxima = input.CapacidadMaxima
	}

	if err := s.repos.Corral.Update(corral); err != nil {
		return nil, err
	}

	return corral, nil
}

// Eliminar da de baja un corral (solo si no tiene lotes activos)
func (s *CorralService) Eliminar(id uint) error {
	_, err := s.repos.Corral.FindByID(id)
	if err != nil {
		return ErrNotFound
	}

	tiene, err := s.repos.Corral.TieneLotesActivos(id)
	if err != nil {
		return err
	}
	if tiene {
		return ErrCorralTieneLotesActivos
	}

	return s.repos.Corral.Delete(id)
}

// GetOcupacion obtiene la cantidad total de animales en el corral
func (s *CorralService) GetOcupacion(id uint) (int, error) {
	_, err := s.repos.Corral.FindByID(id)
	if err != nil {
		return 0, ErrNotFound
	}
	return s.repos.Corral.GetOcupacion(id)
}
