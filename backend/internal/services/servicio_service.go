package services

import (
	"database/sql"
	"time"

	"github.com/martin/sgp/internal/models"
	"github.com/martin/sgp/internal/repositories"
	"gorm.io/gorm"
)

// Constantes de gestación
const (
	DiasGestacion = 114 // Días estimados de gestación porcina
)

// ServicioService maneja la lógica de negocio de servicios
type ServicioService struct {
	db    *gorm.DB
	repos *repositories.RepositoryContainer
}

// NewServicioService crea una nueva instancia del servicio
func NewServicioService(db *gorm.DB, repos *repositories.RepositoryContainer) *ServicioService {
	return &ServicioService{db: db, repos: repos}
}

// --- DTOs ---

// CrearServicioInput datos para registrar un servicio
type CrearServicioInput struct {
	CerdaID              uint   `json:"cerda_id" binding:"required"`
	FechaServicio        string `json:"fecha_servicio" binding:"required"` // formato YYYY-MM-DD
	TipoMonta            string `json:"tipo_monta" binding:"required,oneof=natural inseminacion"`
	PadrilloID           *uint  `json:"padrillo_id"`
	CantidadSaltos       *int   `json:"cantidad_saltos"`
	NumeroPajuela        *string `json:"numero_pajuela"`
	TieneRepeticiones    bool   `json:"tiene_repeticiones"`
	CantidadRepeticiones int    `json:"cantidad_repeticiones"`
}

// ConfirmarPrenezInput datos para confirmar una preñez
type ConfirmarPrenezInput struct {
	FechaConfirmacion string `json:"fecha_confirmacion"` // formato YYYY-MM-DD, si vacío = hoy
}

// CancelarPrenezInput datos para cancelar una preñez
type CancelarPrenezInput struct {
	Motivo            string `json:"motivo" binding:"required"`
	FechaCancelacion  string `json:"fecha_cancelacion"` // formato YYYY-MM-DD, si vacío = hoy
}

// ActualizarServicioInput datos para actualizar un servicio
type ActualizarServicioInput struct {
	FechaServicio        string  `json:"fecha_servicio"`
	TipoMonta            string  `json:"tipo_monta"`
	PadrilloID           *uint   `json:"padrillo_id"`
	CantidadSaltos       *int    `json:"cantidad_saltos"`
	NumeroPajuela        *string `json:"numero_pajuela"`
	TieneRepeticiones    *bool   `json:"tiene_repeticiones"`
	CantidadRepeticiones *int    `json:"cantidad_repeticiones"`
}

// --- Métodos del servicio ---

// Crear registra un nuevo servicio y cambia el estado de la cerda a "servicio"
func (s *ServicioService) Crear(input CrearServicioInput) (*models.Servicio, error) {
	// Validar cerda
	cerda, err := s.repos.Cerda.FindByID(input.CerdaID)
	if err != nil {
		return nil, ErrNotFound
	}
	if !cerda.Activo {
		return nil, ErrCerdaNoActiva
	}
	if cerda.Estado != models.EstadoCerdaDisponible {
		return nil, ErrCerdaNoDisponible
	}

	// Validar campos según tipo de monta
	if input.TipoMonta == models.TipoMontaNatural {
		if input.PadrilloID == nil {
			return nil, ErrServicioRequierePadrillo
		}
		// Verificar que el padrillo existe y está activo
		padrillo, err := s.repos.Padrillo.FindByID(*input.PadrilloID)
		if err != nil {
			return nil, ErrNotFound
		}
		if !padrillo.Activo {
			return nil, ErrForbidden
		}
	} else if input.TipoMonta == models.TipoMontaInseminacion {
		if input.NumeroPajuela == nil || *input.NumeroPajuela == "" {
			return nil, ErrServicioRequierePajuela
		}
	}

	// Parsear fecha
	fechaServicio, err := time.Parse("2006-01-02", input.FechaServicio)
	if err != nil {
		return nil, err
	}

	servicio := &models.Servicio{
		CerdaID:              input.CerdaID,
		FechaServicio:        fechaServicio,
		TipoMonta:            input.TipoMonta,
		PadrilloID:           input.PadrilloID,
		CantidadSaltos:       input.CantidadSaltos,
		NumeroPajuela:        input.NumeroPajuela,
		TieneRepeticiones:    input.TieneRepeticiones,
		CantidadRepeticiones: input.CantidadRepeticiones,
	}

	// Ejecutar en transacción: crear servicio + cambiar estado de cerda
	err = s.repos.Cerda.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(servicio).Error; err != nil {
			return err
		}
		if err := tx.Model(&models.Cerda{}).Where("id = ?", input.CerdaID).
			Update("estado", models.EstadoCerdaServicio).Error; err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return servicio, nil
}

// ObtenerPorID obtiene un servicio con sus relaciones
func (s *ServicioService) ObtenerPorID(id uint) (*models.Servicio, error) {
	servicio, err := s.repos.Servicio.FindByID(id, "Cerda", "Padrillo")
	if err != nil {
		return nil, ErrNotFound
	}
	return servicio, nil
}

// ListarPorCerda lista los servicios de una cerda
func (s *ServicioService) ListarPorCerda(cerdaID uint) ([]models.Servicio, error) {
	return s.repos.Servicio.FindByCerdaID(cerdaID)
}

// ListarPorPeriodo lista servicios por mes/año, opcionalmente filtrados por granja
func (s *ServicioService) ListarPorPeriodo(granjaID *uint, mes, anio int) ([]models.Servicio, error) {
	return s.repos.Servicio.FindByPeriodo(granjaID, mes, anio)
}

// ListarPendientesConfirmacion lista servicios con preñez pendiente de confirmar
func (s *ServicioService) ListarPendientesConfirmacion(granjaID *uint) ([]models.Servicio, error) {
	return s.repos.Servicio.GetServiciosPendientesConfirmacion(granjaID)
}

// Actualizar modifica un servicio antes de que se confirme la preñez
func (s *ServicioService) Actualizar(id uint, input ActualizarServicioInput) (*models.Servicio, error) {
	servicio, err := s.repos.Servicio.FindByID(id)
	if err != nil {
		return nil, ErrNotFound
	}

	// No se puede modificar si la preñez ya fue confirmada (campos críticos)
	if servicio.PrenezConfirmada {
		// Solo permitimos modificar campos no críticos: repeticiones
		if input.TieneRepeticiones != nil {
			servicio.TieneRepeticiones = *input.TieneRepeticiones
		}
		if input.CantidadRepeticiones != nil {
			servicio.CantidadRepeticiones = *input.CantidadRepeticiones
		}
	} else {
		// Se puede modificar todo
		if input.FechaServicio != "" {
			fecha, err := time.Parse("2006-01-02", input.FechaServicio)
			if err != nil {
				return nil, err
			}
			servicio.FechaServicio = fecha
		}
		if input.TipoMonta != "" {
			servicio.TipoMonta = input.TipoMonta
		}
		if input.PadrilloID != nil {
			servicio.PadrilloID = input.PadrilloID
		}
		if input.CantidadSaltos != nil {
			servicio.CantidadSaltos = input.CantidadSaltos
		}
		if input.NumeroPajuela != nil {
			servicio.NumeroPajuela = input.NumeroPajuela
		}
		if input.TieneRepeticiones != nil {
			servicio.TieneRepeticiones = *input.TieneRepeticiones
		}
		if input.CantidadRepeticiones != nil {
			servicio.CantidadRepeticiones = *input.CantidadRepeticiones
		}
	}

	if err := s.repos.Servicio.Update(servicio); err != nil {
		return nil, err
	}

	return servicio, nil
}

// ConfirmarPrenez confirma la preñez de un servicio y cambia el estado a "gestación"
func (s *ServicioService) ConfirmarPrenez(servicioID uint, input ConfirmarPrenezInput) (*models.Servicio, error) {
	servicio, err := s.repos.Servicio.FindByID(servicioID, "Cerda")
	if err != nil {
		return nil, ErrNotFound
	}

	// Validar que la cerda está en estado "servicio"
	if servicio.Cerda.Estado != models.EstadoCerdaServicio {
		return nil, ErrCerdaNoEnServicio
	}

	if servicio.PrenezConfirmada {
		return nil, ErrPrenezYaConfirmada
	}
	if servicio.PrenezCancelada {
		return nil, ErrPrenezYaCancelada
	}

	// Parsear fecha de confirmación
	fechaConfirmacion := time.Now()
	if input.FechaConfirmacion != "" {
		fechaConfirmacion, err = time.Parse("2006-01-02", input.FechaConfirmacion)
		if err != nil {
			return nil, err
		}
	}

	// Calcular fecha estimada de parto: fecha_servicio + 114 días
	fechaEstimadaParto := servicio.FechaServicio.AddDate(0, 0, DiasGestacion)

	// Transacción: confirmar preñez + cambiar estado cerda
	err = s.repos.Cerda.Transaction(func(tx *gorm.DB) error {
		// Actualizar servicio
		if err := tx.Model(&models.Servicio{}).Where("id = ?", servicioID).Updates(map[string]interface{}{
			"prenez_confirmada":        true,
			"fecha_confirmacion_prenez": sql.NullTime{Time: fechaConfirmacion, Valid: true},
			"fecha_estimada_parto":     fechaEstimadaParto,
		}).Error; err != nil {
			return err
		}

		// Cambiar estado de la cerda a "gestación"
		if err := tx.Model(&models.Cerda{}).Where("id = ?", servicio.CerdaID).
			Update("estado", models.EstadoCerdaGestacion).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	// Recargar servicio con datos actualizados
	servicio, _ = s.repos.Servicio.FindByID(servicioID, "Cerda", "Padrillo")
	return servicio, nil
}

// CancelarPrenez cancela una preñez y devuelve la cerda a "disponible"
func (s *ServicioService) CancelarPrenez(servicioID uint, input CancelarPrenezInput) (*models.Servicio, error) {
	servicio, err := s.repos.Servicio.FindByID(servicioID, "Cerda")
	if err != nil {
		return nil, ErrNotFound
	}

	if servicio.PrenezCancelada {
		return nil, ErrPrenezYaCancelada
	}

	// Se puede cancelar desde estado "servicio" (antes de confirmar) o "gestación" (después de confirmar)
	if servicio.Cerda.Estado != models.EstadoCerdaServicio && servicio.Cerda.Estado != models.EstadoCerdaGestacion {
		return nil, ErrForbidden
	}

	if input.Motivo == "" {
		return nil, ErrMotivoRequerido
	}

	// Parsear fecha de cancelación
	fechaCancelacion := time.Now()
	if input.FechaCancelacion != "" {
		fechaCancelacion, err = time.Parse("2006-01-02", input.FechaCancelacion)
		if err != nil {
			return nil, err
		}
	}

	// Transacción: cancelar preñez + devolver cerda a disponible
	err = s.repos.Cerda.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&models.Servicio{}).Where("id = ?", servicioID).Updates(map[string]interface{}{
			"prenez_cancelada":         true,
			"fecha_cancelacion_prenez": sql.NullTime{Time: fechaCancelacion, Valid: true},
			"motivo_cancelacion":       input.Motivo,
		}).Error; err != nil {
			return err
		}

		// Devolver cerda a estado disponible
		if err := tx.Model(&models.Cerda{}).Where("id = ?", servicio.CerdaID).
			Update("estado", models.EstadoCerdaDisponible).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	servicio, _ = s.repos.Servicio.FindByID(servicioID, "Cerda", "Padrillo")
	return servicio, nil
}
