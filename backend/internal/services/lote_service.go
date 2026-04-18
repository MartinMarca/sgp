package services

import (
	"database/sql"
	"time"

	"github.com/martin/sgp/internal/models"
	"github.com/martin/sgp/internal/repositories"
)

// LoteService maneja la lógica de negocio de lotes
type LoteService struct {
	repos *repositories.RepositoryContainer
}

// NewLoteService crea una nueva instancia del servicio
func NewLoteService(repos *repositories.RepositoryContainer) *LoteService {
	return &LoteService{repos: repos}
}

// --- DTOs ---

// CrearLoteInput datos para crear un lote anticipadamente (sin destete)
type CrearLoteInput struct {
	CorralID uint   `json:"corral_id"`
	Nombre   string `json:"nombre" binding:"required"`
	Fecha    string `json:"fecha"` // formato YYYY-MM-DD, si vacío = hoy
}

// ActualizarLoteInput datos para actualizar un lote
type ActualizarLoteInput struct {
	Nombre           string `json:"nombre"`
	CantidadLechones *int   `json:"cantidad_lechones"`
}

// CerrarLoteInput datos para cerrar un lote
type CerrarLoteInput struct {
	MotivoCierre string `json:"motivo_cierre" binding:"required"`
	Estado       string `json:"estado" binding:"required,oneof=cerrado vendido"`
}

// --- Métodos del servicio ---

// Crear registra un nuevo lote anticipadamente (puede tener 0 lechones)
func (s *LoteService) Crear(input CrearLoteInput) (*models.Lote, error) {
	// Verificar que el corral existe
	_, err := s.repos.Corral.FindByID(input.CorralID)
	if err != nil {
		return nil, ErrNotFound
	}

	fecha := time.Now()
	if input.Fecha != "" {
		fecha, err = time.Parse("2006-01-02", input.Fecha)
		if err != nil {
			return nil, err
		}
	}

	lote := &models.Lote{
		CorralID:         input.CorralID,
		Nombre:           input.Nombre,
		CantidadLechones: 0,
		FechaCreacion:    fecha,
		Estado:           models.EstadoLoteActivo,
	}

	if err := s.repos.Lote.Create(lote); err != nil {
		return nil, err
	}

	return lote, nil
}

// ObtenerPorID obtiene un lote por ID con relaciones
func (s *LoteService) ObtenerPorID(id uint) (*models.Lote, error) {
	lote, err := s.repos.Lote.FindByID(id, "Corral", "Destetes")
	if err != nil {
		return nil, ErrNotFound
	}
	return lote, nil
}

// ListarPorCorral lista lotes por corral, opcionalmente filtrados por estado
func (s *LoteService) ListarPorCorral(corralID uint, estado *string) ([]models.Lote, error) {
	return s.repos.Lote.FindByCorralID(corralID, estado)
}

// ListarPorEstado lista todos los lotes por estado
func (s *LoteService) ListarPorEstado(estado string) ([]models.Lote, error) {
	return s.repos.Lote.FindByEstado(estado)
}

// Actualizar modifica nombre y/o cantidad de lechones de un lote
func (s *LoteService) Actualizar(id uint, input ActualizarLoteInput) (*models.Lote, error) {
	lote, err := s.repos.Lote.FindByID(id)
	if err != nil {
		return nil, ErrNotFound
	}

	if lote.Estado != models.EstadoLoteActivo {
		return nil, ErrLoteNoActivo
	}

	if input.Nombre != "" {
		lote.Nombre = input.Nombre
	}
	if input.CantidadLechones != nil {
		lote.CantidadLechones = *input.CantidadLechones
	}

	if err := s.repos.Lote.Update(lote); err != nil {
		return nil, err
	}

	return lote, nil
}

// Cerrar cambia el estado de un lote a cerrado o vendido
func (s *LoteService) Cerrar(id uint, input CerrarLoteInput) (*models.Lote, error) {
	lote, err := s.repos.Lote.FindByID(id)
	if err != nil {
		return nil, ErrNotFound
	}

	if lote.Estado != models.EstadoLoteActivo {
		return nil, ErrLoteNoActivo
	}

	ahora := time.Now()
	lote.Estado = input.Estado
	lote.FechaCierre = &sql.NullTime{Time: ahora, Valid: true}
	lote.MotivoCierre = &input.MotivoCierre

	if err := s.repos.Lote.Update(lote); err != nil {
		return nil, err
	}

	return lote, nil
}

// GetDestetes obtiene los destetes asociados a un lote
func (s *LoteService) GetDestetes(id uint) ([]models.Destete, error) {
	_, err := s.repos.Lote.FindByID(id)
	if err != nil {
		return nil, ErrNotFound
	}
	return s.repos.Lote.GetDestetes(id)
}
