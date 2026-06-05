package services

import (
	"database/sql"
	"time"

	"github.com/martin/sgp/internal/models"
	"github.com/martin/sgp/internal/repositories"
)

type CerdaService struct {
	repos *repositories.RepositoryContainer
}

func NewCerdaService(repos *repositories.RepositoryContainer) *CerdaService {
	return &CerdaService{repos: repos}
}

type CrearCerdaInput struct {
	GranjaID       uint   `json:"granja_id"`
	NumeroCaravana string `json:"numero_caravana" binding:"required"`
	DetallePelaje  string `json:"detalle_pelaje"`
	Genetica       string `json:"genetica"`
	Estado         string `json:"estado"`
}

type ActualizarCerdaInput struct {
	NumeroCaravana string `json:"numero_caravana"`
	DetallePelaje  string `json:"detalle_pelaje"`
	Genetica       string `json:"genetica"`
}

type BajaCerdaInput struct {
	MotivoBaja string `json:"motivo_baja" binding:"required,oneof=muerte venta"`
}

func (s *CerdaService) Crear(input CrearCerdaInput) (*models.Cerda, error) {
	existe, err := s.repos.Cerda.ExisteCaravana(input.GranjaID, input.NumeroCaravana, nil)
	if err != nil {
		return nil, err
	}
	if existe {
		return nil, ErrCaravanaDuplicada
	}

	estado := input.Estado
	if estado == "" {
		estado = models.EstadoCerdaDisponible
	}

	if !esEstadoCerdaValido(estado) {
		return nil, ErrCerdaNoDisponible
	}

	cerda := &models.Cerda{
		GranjaID:       input.GranjaID,
		NumeroCaravana: input.NumeroCaravana,
		Estado:         estado,
		Activo:         true,
	}

	if input.DetallePelaje != "" {
		cerda.DetallePelaje = &input.DetallePelaje
	}
	if input.Genetica != "" {
		cerda.Genetica = &input.Genetica
	}

	if err := s.repos.Cerda.Create(cerda); err != nil {
		return nil, err
	}

	return cerda, nil
}

func (s *CerdaService) ObtenerPorID(id uint) (*models.Cerda, error) {
	cerda, err := s.repos.Cerda.FindByID(id, "Granja")
	if err != nil {
		return nil, ErrNotFound
	}
	return cerda, nil
}

func (s *CerdaService) ListarPorGranja(granjaID uint, estado *string, activo *bool) ([]models.Cerda, error) {
	return s.repos.Cerda.FindByGranjaID(granjaID, estado, activo)
}

func (s *CerdaService) ListarPorEstado(estado string, granjaID *uint) ([]models.Cerda, error) {
	return s.repos.Cerda.FindByEstado(estado, granjaID)
}

func (s *CerdaService) Actualizar(id uint, input ActualizarCerdaInput) (*models.Cerda, error) {
	cerda, err := s.repos.Cerda.FindByID(id)
	if err != nil {
		return nil, ErrNotFound
	}

	if !cerda.Activo {
		return nil, ErrCerdaNoActiva
	}

	if input.NumeroCaravana != "" && input.NumeroCaravana != cerda.NumeroCaravana {
		existe, err := s.repos.Cerda.ExisteCaravana(cerda.GranjaID, input.NumeroCaravana, &id)
		if err != nil {
			return nil, err
		}
		if existe {
			return nil, ErrCaravanaDuplicada
		}
		cerda.NumeroCaravana = input.NumeroCaravana
	}

	if input.DetallePelaje != "" {
		cerda.DetallePelaje = &input.DetallePelaje
	}
	if input.Genetica != "" {
		cerda.Genetica = &input.Genetica
	}

	if err := s.repos.Cerda.Update(cerda); err != nil {
		return nil, err
	}

	return cerda, nil
}

func (s *CerdaService) DarDeBaja(id uint, input BajaCerdaInput) (*models.Cerda, error) {
	cerda, err := s.repos.Cerda.FindByID(id)
	if err != nil {
		return nil, ErrNotFound
	}

	if !cerda.Activo {
		return nil, ErrCerdaNoActiva
	}

	// Si la cerda está en servicio o gestación, no se puede dar de baja directamente
	if cerda.Estado == models.EstadoCerdaServicio || cerda.Estado == models.EstadoCerdaGestacion {
		return nil, ErrCerdaTieneServicioActivo
	}

	ahora := time.Now()
	cerda.Activo = false
	cerda.FechaBaja = &sql.NullTime{Time: ahora, Valid: true}
	cerda.MotivoBaja = &input.MotivoBaja

	if err := s.repos.Cerda.Update(cerda); err != nil {
		return nil, err
	}

	return cerda, nil
}

// GetHistorial obtiene el historial completo de una cerda
func (s *CerdaService) GetHistorial(id uint) (map[string]interface{}, error) {
	// Verificar que la cerda existe
	_, err := s.repos.Cerda.FindByID(id)
	if err != nil {
		return nil, ErrNotFound
	}
	return s.repos.Cerda.GetHistorial(id)
}

// GetEstadisticas obtiene estadísticas de una cerda
func (s *CerdaService) GetEstadisticas(id uint) (map[string]interface{}, error) {
	_, err := s.repos.Cerda.FindByID(id)
	if err != nil {
		return nil, ErrNotFound
	}
	return s.repos.Cerda.GetEstadisticas(id)
}

func esEstadoCerdaValido(estado string) bool {
	switch estado {
	case models.EstadoCerdaDisponible,
		models.EstadoCerdaServicio,
		models.EstadoCerdaGestacion,
		models.EstadoCerdaCria:
		return true
	}
	return false
}
