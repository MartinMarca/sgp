package services

import (
	"time"

	"github.com/martin/sgp/internal/models"
	"github.com/martin/sgp/internal/repositories"
	"gorm.io/gorm"
)

// DesteteService maneja la lógica de negocio de destetes
type DesteteService struct {
	db    *gorm.DB
	repos *repositories.RepositoryContainer
}

// NewDesteteService crea una nueva instancia del servicio
func NewDesteteService(db *gorm.DB, repos *repositories.RepositoryContainer) *DesteteService {
	return &DesteteService{db: db, repos: repos}
}

// --- DTOs ---

// CrearDesteteInput datos para registrar un destete
// Flujo: el usuario selecciona una cerda en estado "cría"
type CrearDesteteInput struct {
	CerdaID                     uint   `json:"cerda_id" binding:"required"`
	FechaDestete                string `json:"fecha_destete" binding:"required"` // formato YYYY-MM-DD
	CantidadLechonesDestetados  int    `json:"cantidad_lechones_destetados" binding:"min=0"`
	// Asignación a lote: usar lote existente o crear uno nuevo
	LoteID     *uint            `json:"lote_id"`     // Si se asigna a un lote existente
	NuevoLote  *NuevoLoteInput  `json:"nuevo_lote"`  // Si se crea un lote nuevo
}

// NuevoLoteInput datos para crear un lote junto con el destete
type NuevoLoteInput struct {
	CorralID uint   `json:"corral_id" binding:"required"`
	Nombre   string `json:"nombre" binding:"required"`
}

// ActualizarDesteteInput datos para actualizar un destete
type ActualizarDesteteInput struct {
	FechaDestete               *string `json:"fecha_destete"`
	CantidadLechonesDestetados *int    `json:"cantidad_lechones_destetados"`
}

// --- Métodos del servicio ---

// Crear registra un nuevo destete, asigna a lote y devuelve la cerda a "disponible"
func (s *DesteteService) Crear(input CrearDesteteInput) (*models.Destete, error) {
	// Validar cerda
	cerda, err := s.repos.Cerda.FindByID(input.CerdaID)
	if err != nil {
		return nil, ErrNotFound
	}
	if !cerda.Activo {
		return nil, ErrCerdaNoActiva
	}
	if cerda.Estado != models.EstadoCerdaCria {
		return nil, ErrCerdaNoEnCria
	}

	// Obtener el parto sin destete de esta cerda
	parto, err := s.repos.Cerda.GetPartoSinDestete(input.CerdaID)
	if err != nil {
		return nil, ErrNoHayPartoSinDestete
	}

	// Validar: cantidad destetados <= nacidos vivos del parto
	if input.CantidadLechonesDestetados > parto.LechonesNacidosVivos {
		return nil, ErrDestetadosExcedenVivos
	}

	// Validar que se asigne un lote (existente o nuevo)
	if input.LoteID == nil && input.NuevoLote == nil {
		return nil, ErrLoteRequerido
	}

	// Parsear fecha
	fechaDestete, err := time.Parse("2006-01-02", input.FechaDestete)
	if err != nil {
		return nil, err
	}

	var destete *models.Destete

	// Ejecutar todo en transacción
	err = s.repos.Cerda.Transaction(func(tx *gorm.DB) error {
		var loteID uint

		if input.LoteID != nil {
			// Usar lote existente: validar que esté activo
			var lote models.Lote
			if err := tx.First(&lote, *input.LoteID).Error; err != nil {
				return ErrNotFound
			}
			if lote.Estado != models.EstadoLoteActivo {
				return ErrLoteNoActivo
			}
			loteID = lote.ID

			// Sumar lechones al lote existente
			if err := tx.Model(&models.Lote{}).Where("id = ?", loteID).
				Update("cantidad_lechones", gorm.Expr("cantidad_lechones + ?", input.CantidadLechonesDestetados)).Error; err != nil {
				return err
			}
		} else {
			// Crear lote nuevo
			nuevoLote := &models.Lote{
				CorralID:          input.NuevoLote.CorralID,
				Nombre:            input.NuevoLote.Nombre,
				CantidadLechones:  input.CantidadLechonesDestetados,
				FechaCreacion:     fechaDestete,
				Estado:            models.EstadoLoteActivo,
			}
			if err := tx.Create(nuevoLote).Error; err != nil {
				return err
			}
			loteID = nuevoLote.ID
		}

		// Crear destete
		destete = &models.Destete{
			CerdaID:                     input.CerdaID,
			PartoID:                     &parto.ID,
			FechaDestete:                fechaDestete,
			CantidadLechonesDestetados:  input.CantidadLechonesDestetados,
			FechaEstimada:               parto.FechaEstimada, // Copiamos la fecha estimada del parto
			LoteID:                      loteID,
		}
		if err := tx.Create(destete).Error; err != nil {
			return err
		}

		// Devolver cerda a estado "disponible"
		if err := tx.Model(&models.Cerda{}).Where("id = ?", input.CerdaID).
			Update("estado", models.EstadoCerdaDisponible).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return destete, nil
}

// ObtenerPorID obtiene un destete con sus relaciones
func (s *DesteteService) ObtenerPorID(id uint) (*models.Destete, error) {
	destete, err := s.repos.Destete.FindByID(id, "Cerda", "Parto", "Lote")
	if err != nil {
		return nil, ErrNotFound
	}
	return destete, nil
}

// ListarPorCerda lista los destetes de una cerda
func (s *DesteteService) ListarPorCerda(cerdaID uint) ([]models.Destete, error) {
	return s.repos.Destete.FindByCerdaID(cerdaID)
}

// ListarPorLote lista los destetes asignados a un lote
func (s *DesteteService) ListarPorLote(loteID uint) ([]models.Destete, error) {
	return s.repos.Destete.FindByLoteID(loteID)
}

// ListarPorPeriodo lista destetes por mes/año
func (s *DesteteService) ListarPorPeriodo(granjaID *uint, mes, anio int) ([]models.Destete, error) {
	return s.repos.Destete.FindByPeriodo(granjaID, mes, anio)
}

// Actualizar modifica la información de un destete
// Si cambia la cantidad de destetados, se actualiza también el lote
func (s *DesteteService) Actualizar(id uint, input ActualizarDesteteInput) (*models.Destete, error) {
	destete, err := s.repos.Destete.FindByID(id, "Parto")
	if err != nil {
		return nil, ErrNotFound
	}

	cantidadAnterior := destete.CantidadLechonesDestetados

	if input.FechaDestete != nil {
		fecha, err := time.Parse("2006-01-02", *input.FechaDestete)
		if err != nil {
			return nil, err
		}
		destete.FechaDestete = fecha
	}

	if input.CantidadLechonesDestetados != nil {
		nuevaCantidad := *input.CantidadLechonesDestetados

		// Validar contra nacidos vivos del parto
		if destete.Parto != nil && nuevaCantidad > destete.Parto.LechonesNacidosVivos {
			return nil, ErrDestetadosExcedenVivos
		}

		destete.CantidadLechonesDestetados = nuevaCantidad
	}

	// Transacción: actualizar destete + ajustar lote si cambió la cantidad
	err = s.repos.Cerda.Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(destete).Error; err != nil {
			return err
		}

		// Si la cantidad cambió, ajustar el lote
		if input.CantidadLechonesDestetados != nil {
			diferencia := destete.CantidadLechonesDestetados - cantidadAnterior
			if diferencia != 0 {
				if err := tx.Model(&models.Lote{}).Where("id = ?", destete.LoteID).
					Update("cantidad_lechones", gorm.Expr("cantidad_lechones + ?", diferencia)).Error; err != nil {
					return err
				}
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return destete, nil
}

// GetEstadisticas obtiene estadísticas de destetes
func (s *DesteteService) GetEstadisticas(granjaID *uint, mes, anio int) (map[string]interface{}, error) {
	return s.repos.Destete.GetEstadisticas(granjaID, mes, anio)
}
