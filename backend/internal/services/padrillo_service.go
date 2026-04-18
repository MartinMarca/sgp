package services

import (
	"database/sql"
	"time"

	"github.com/martin/sgp/internal/models"
	"github.com/martin/sgp/internal/repositories"
)

// PadrilloService maneja la lógica de negocio de padrillos
type PadrilloService struct {
	repos *repositories.RepositoryContainer
}

// NewPadrilloService crea una nueva instancia del servicio
func NewPadrilloService(repos *repositories.RepositoryContainer) *PadrilloService {
	return &PadrilloService{repos: repos}
}

// --- DTOs ---

// CrearPadrilloInput datos para crear un padrillo
type CrearPadrilloInput struct {
	GranjaID              uint   `json:"granja_id"`
	NumeroCaravana        string `json:"numero_caravana" binding:"required"`
	Nombre                string `json:"nombre" binding:"required"`
	Genetica              string `json:"genetica"`
	FechaUltimaVacunacion string `json:"fecha_ultima_vacunacion"` // formato YYYY-MM-DD
}

// ActualizarPadrilloInput datos para actualizar un padrillo
type ActualizarPadrilloInput struct {
	NumeroCaravana        string `json:"numero_caravana"`
	Nombre                string `json:"nombre"`
	Genetica              string `json:"genetica"`
	FechaUltimaVacunacion string `json:"fecha_ultima_vacunacion"`
}

// BajaPadrilloInput datos para dar de baja un padrillo
type BajaPadrilloInput struct {
	MotivoBaja string `json:"motivo_baja" binding:"required,oneof=muerte venta"`
}

// --- Métodos del servicio ---

// Crear registra un nuevo padrillo
func (s *PadrilloService) Crear(input CrearPadrilloInput) (*models.Padrillo, error) {
	// Validar caravana duplicada
	existe, err := s.repos.Padrillo.ExisteCaravana(input.GranjaID, input.NumeroCaravana, nil)
	if err != nil {
		return nil, err
	}
	if existe {
		return nil, ErrCaravanaDuplicada
	}

	padrillo := &models.Padrillo{
		GranjaID:       input.GranjaID,
		NumeroCaravana: input.NumeroCaravana,
		Nombre:         input.Nombre,
		Activo:         true,
	}

	if input.Genetica != "" {
		padrillo.Genetica = &input.Genetica
	}
	if input.FechaUltimaVacunacion != "" {
		fecha, err := time.Parse("2006-01-02", input.FechaUltimaVacunacion)
		if err != nil {
			return nil, err
		}
		padrillo.FechaUltimaVacunacion = &sql.NullTime{Time: fecha, Valid: true}
	}

	if err := s.repos.Padrillo.Create(padrillo); err != nil {
		return nil, err
	}

	return padrillo, nil
}

// ObtenerPorID obtiene un padrillo por ID
func (s *PadrilloService) ObtenerPorID(id uint) (*models.Padrillo, error) {
	padrillo, err := s.repos.Padrillo.FindByID(id, "Granja")
	if err != nil {
		return nil, ErrNotFound
	}
	return padrillo, nil
}

// ListarPorGranja lista padrillos por granja
func (s *PadrilloService) ListarPorGranja(granjaID uint, activo *bool) ([]models.Padrillo, error) {
	return s.repos.Padrillo.FindByGranjaID(granjaID, activo)
}

// Actualizar modifica los datos de un padrillo
func (s *PadrilloService) Actualizar(id uint, input ActualizarPadrilloInput) (*models.Padrillo, error) {
	padrillo, err := s.repos.Padrillo.FindByID(id)
	if err != nil {
		return nil, ErrNotFound
	}

	if !padrillo.Activo {
		return nil, ErrCerdaNoActiva // Reutilizamos el error genérico
	}

	// Validar caravana si cambió
	if input.NumeroCaravana != "" && input.NumeroCaravana != padrillo.NumeroCaravana {
		existe, err := s.repos.Padrillo.ExisteCaravana(padrillo.GranjaID, input.NumeroCaravana, &id)
		if err != nil {
			return nil, err
		}
		if existe {
			return nil, ErrCaravanaDuplicada
		}
		padrillo.NumeroCaravana = input.NumeroCaravana
	}

	if input.Nombre != "" {
		padrillo.Nombre = input.Nombre
	}
	if input.Genetica != "" {
		padrillo.Genetica = &input.Genetica
	}
	if input.FechaUltimaVacunacion != "" {
		fecha, err := time.Parse("2006-01-02", input.FechaUltimaVacunacion)
		if err != nil {
			return nil, err
		}
		padrillo.FechaUltimaVacunacion = &sql.NullTime{Time: fecha, Valid: true}
	}

	if err := s.repos.Padrillo.Update(padrillo); err != nil {
		return nil, err
	}

	return padrillo, nil
}

// DarDeBaja da de baja un padrillo por muerte o venta
func (s *PadrilloService) DarDeBaja(id uint, input BajaPadrilloInput) (*models.Padrillo, error) {
	padrillo, err := s.repos.Padrillo.FindByID(id)
	if err != nil {
		return nil, ErrNotFound
	}

	if !padrillo.Activo {
		return nil, ErrCerdaNoActiva
	}

	ahora := time.Now()
	padrillo.Activo = false
	padrillo.FechaBaja = &sql.NullTime{Time: ahora, Valid: true}
	padrillo.MotivoBaja = &input.MotivoBaja

	if err := s.repos.Padrillo.Update(padrillo); err != nil {
		return nil, err
	}

	return padrillo, nil
}

// GetEstadisticas obtiene estadísticas de un padrillo
func (s *PadrilloService) GetEstadisticas(id uint) (map[string]interface{}, error) {
	_, err := s.repos.Padrillo.FindByID(id)
	if err != nil {
		return nil, ErrNotFound
	}
	return s.repos.Padrillo.GetEstadisticas(id)
}
