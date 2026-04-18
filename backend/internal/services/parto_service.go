package services

import (
	"time"

	"github.com/martin/sgp/internal/models"
	"github.com/martin/sgp/internal/repositories"
	"gorm.io/gorm"
)

// Constantes de cría
const (
	DiasCria = 30 // Días estimados de cría hasta destete
)

// PartoService maneja la lógica de negocio de partos
type PartoService struct {
	db    *gorm.DB
	repos *repositories.RepositoryContainer
}

// NewPartoService crea una nueva instancia del servicio
func NewPartoService(db *gorm.DB, repos *repositories.RepositoryContainer) *PartoService {
	return &PartoService{db: db, repos: repos}
}

// --- DTOs ---

// CrearPartoInput datos para registrar un parto
// Flujo: el usuario selecciona una cerda en estado "gestación"
type CrearPartoInput struct {
	CerdaID                uint   `json:"cerda_id" binding:"required"`
	FechaParto             string `json:"fecha_parto" binding:"required"` // formato YYYY-MM-DD
	LechonesNacidosVivos   int    `json:"lechones_nacidos_vivos" binding:"min=0"`
	LechonesNacidosTotales int    `json:"lechones_nacidos_totales" binding:"min=0"`
	LechonesHembras        int    `json:"lechones_hembras" binding:"min=0"`
	LechonesMachos         int    `json:"lechones_machos" binding:"min=0"`
}

// ActualizarPartoInput datos para actualizar un parto
type ActualizarPartoInput struct {
	FechaParto             *string `json:"fecha_parto"`
	LechonesNacidosVivos   *int    `json:"lechones_nacidos_vivos"`
	LechonesNacidosTotales *int    `json:"lechones_nacidos_totales"`
	LechonesHembras        *int    `json:"lechones_hembras"`
	LechonesMachos         *int    `json:"lechones_machos"`
}

// --- Métodos del servicio ---

// Crear registra un nuevo parto y cambia el estado de la cerda a "cría"
func (s *PartoService) Crear(input CrearPartoInput) (*models.Parto, error) {
	// Validar cerda
	cerda, err := s.repos.Cerda.FindByID(input.CerdaID)
	if err != nil {
		return nil, ErrNotFound
	}
	if !cerda.Activo {
		return nil, ErrCerdaNoActiva
	}
	if cerda.Estado != models.EstadoCerdaGestacion {
		return nil, ErrCerdaNoEnGestacion
	}

	// Obtener el servicio con preñez confirmada de esta cerda
	servicio, err := s.repos.Servicio.GetServicioConPrenezConfirmada(input.CerdaID)
	if err != nil {
		return nil, ErrNoHayServicioConfirmado
	}

	// Validar lechones: hembras + machos == vivos
	if input.LechonesHembras+input.LechonesMachos != input.LechonesNacidosVivos {
		return nil, ErrLechonesInvalidos
	}

	// Validar: totales >= vivos
	if input.LechonesNacidosTotales < input.LechonesNacidosVivos {
		return nil, ErrTotalMenorQueVivos
	}

	// Parsear fecha
	fechaParto, err := time.Parse("2006-01-02", input.FechaParto)
	if err != nil {
		return nil, err
	}

	// Calcular fecha estimada de destete: fecha_parto + 30 días
	fechaEstimadaDestete := fechaParto.AddDate(0, 0, DiasCria)

	parto := &models.Parto{
		CerdaID:                input.CerdaID,
		ServicioID:             &servicio.ID,
		FechaParto:             fechaParto,
		LechonesNacidosVivos:   input.LechonesNacidosVivos,
		LechonesNacidosTotales: input.LechonesNacidosTotales,
		LechonesHembras:        input.LechonesHembras,
		LechonesMachos:         input.LechonesMachos,
		FechaEstimada:          fechaEstimadaDestete,
	}

	// Transacción: crear parto + cambiar estado de cerda a cría
	err = s.repos.Cerda.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(parto).Error; err != nil {
			return err
		}

		// Cambiar estado de la cerda a "cría"
		if err := tx.Model(&models.Cerda{}).Where("id = ?", input.CerdaID).
			Update("estado", models.EstadoCerdaCria).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return parto, nil
}

// ObtenerPorID obtiene un parto con sus relaciones
func (s *PartoService) ObtenerPorID(id uint) (*models.Parto, error) {
	parto, err := s.repos.Parto.FindByID(id, "Cerda", "Servicio", "Servicio.Padrillo")
	if err != nil {
		return nil, ErrNotFound
	}
	return parto, nil
}

// ListarPorCerda lista los partos de una cerda
func (s *PartoService) ListarPorCerda(cerdaID uint) ([]models.Parto, error) {
	return s.repos.Parto.FindByCerdaID(cerdaID)
}

// ListarPorPeriodo lista partos por mes/año, opcionalmente filtrados por granja
func (s *PartoService) ListarPorPeriodo(granjaID *uint, mes, anio int) ([]models.Parto, error) {
	return s.repos.Parto.FindByPeriodo(granjaID, mes, anio)
}

// Actualizar modifica la información de un parto
func (s *PartoService) Actualizar(id uint, input ActualizarPartoInput) (*models.Parto, error) {
	parto, err := s.repos.Parto.FindByID(id)
	if err != nil {
		return nil, ErrNotFound
	}

	if input.FechaParto != nil {
		fecha, err := time.Parse("2006-01-02", *input.FechaParto)
		if err != nil {
			return nil, err
		}
		parto.FechaParto = fecha
		// Recalcular fecha estimada de destete
		parto.FechaEstimada = fecha.AddDate(0, 0, DiasCria)
	}

	if input.LechonesNacidosVivos != nil {
		parto.LechonesNacidosVivos = *input.LechonesNacidosVivos
	}
	if input.LechonesNacidosTotales != nil {
		parto.LechonesNacidosTotales = *input.LechonesNacidosTotales
	}
	if input.LechonesHembras != nil {
		parto.LechonesHembras = *input.LechonesHembras
	}
	if input.LechonesMachos != nil {
		parto.LechonesMachos = *input.LechonesMachos
	}

	// Re-validar: hembras + machos == vivos
	if parto.LechonesHembras+parto.LechonesMachos != parto.LechonesNacidosVivos {
		return nil, ErrLechonesInvalidos
	}
	if parto.LechonesNacidosTotales < parto.LechonesNacidosVivos {
		return nil, ErrTotalMenorQueVivos
	}

	if err := s.repos.Parto.Update(parto); err != nil {
		return nil, err
	}

	return parto, nil
}

// GetEstadisticas obtiene estadísticas de partos
func (s *PartoService) GetEstadisticas(granjaID *uint, mes, anio int) (map[string]interface{}, error) {
	return s.repos.Parto.GetEstadisticas(granjaID, mes, anio)
}
