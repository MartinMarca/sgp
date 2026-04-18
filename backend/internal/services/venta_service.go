package services

import (
	"time"

	"github.com/martin/sgp/internal/models"
	"github.com/martin/sgp/internal/repositories"
	"gorm.io/gorm"
)

// VentaService maneja la lógica de negocio de ventas
type VentaService struct {
	db    *gorm.DB
	repos *repositories.RepositoryContainer
}

func NewVentaService(db *gorm.DB, repos *repositories.RepositoryContainer) *VentaService {
	return &VentaService{db: db, repos: repos}
}

// --- DTOs ---

type CrearVentaInput struct {
	GranjaID   uint    `json:"granja_id" binding:"required"`
	Fecha      string  `json:"fecha" binding:"required"`
	TipoAnimal string  `json:"tipo_animal" binding:"required"`
	Cantidad   int     `json:"cantidad" binding:"required,gte=1"`
	KgTotales  float64 `json:"kg_totales" binding:"required,gte=0"`
	Monto      float64 `json:"monto" binding:"required,gte=0"`
	Comprador  string  `json:"comprador" binding:"required"`
	LoteID     *uint   `json:"lote_id"`
	CorralID   *uint   `json:"corral_id"`
	Notas      string  `json:"notas"`
}

type ActualizarVentaInput struct {
	Fecha     *string  `json:"fecha"`
	TipoAnimal *string `json:"tipo_animal"`
	Cantidad  *int     `json:"cantidad"`
	KgTotales *float64 `json:"kg_totales"`
	Monto     *float64 `json:"monto"`
	Comprador *string  `json:"comprador"`
	LoteID    *uint    `json:"lote_id"`
	CorralID  *uint    `json:"corral_id"`
	Notas     *string  `json:"notas"`
}

var tiposAnimalValidos = map[string]bool{
	models.TipoAnimalVentaCerda:    true,
	models.TipoAnimalVentaPadrillo: true,
	models.TipoAnimalVentaLechon:   true,
}

// Crear registra una nueva venta
func (s *VentaService) Crear(input CrearVentaInput) (*models.Venta, error) {
	if !tiposAnimalValidos[input.TipoAnimal] {
		return nil, ErrTipoAnimalInvalido
	}

	fecha, err := time.Parse("2006-01-02", input.Fecha)
	if err != nil {
		return nil, err
	}

	var notas *string
	if input.Notas != "" {
		notas = &input.Notas
	}

	venta := &models.Venta{
		GranjaID:   input.GranjaID,
		Fecha:      fecha,
		TipoAnimal: input.TipoAnimal,
		Cantidad:   input.Cantidad,
		KgTotales:  input.KgTotales,
		Monto:      input.Monto,
		Comprador:  input.Comprador,
		LoteID:     input.LoteID,
		CorralID:   input.CorralID,
		Notas:      notas,
	}

	err = s.db.Transaction(func(tx *gorm.DB) error {
		// Si es venta de lechones con lote, validar y decrementar
		if input.TipoAnimal == models.TipoAnimalVentaLechon && input.LoteID != nil {
			var lote models.Lote
			if err := tx.First(&lote, *input.LoteID).Error; err != nil {
				return ErrNotFound
			}
			if lote.Estado != models.EstadoLoteActivo {
				return ErrLoteNoActivo
			}
			if input.Cantidad > lote.CantidadLechones {
				return ErrVentaExcedeLote
			}
			if err := tx.Model(&models.Lote{}).Where("id = ?", *input.LoteID).
				Update("cantidad_lechones", gorm.Expr("cantidad_lechones - ?", input.Cantidad)).Error; err != nil {
				return err
			}
		}

		return tx.Create(venta).Error
	})
	if err != nil {
		return nil, err
	}
	return venta, nil
}

// ObtenerPorID obtiene una venta por ID
func (s *VentaService) ObtenerPorID(id uint) (*models.Venta, error) {
	venta, err := s.repos.Venta.FindByID(id, "Granja", "Lote", "Lote.Corral", "Corral")
	if err != nil {
		return nil, ErrNotFound
	}
	return venta, nil
}

// ListarPorPeriodo lista ventas filtrando por granja/mes/año
func (s *VentaService) ListarPorPeriodo(granjaID *uint, mes, anio int) ([]models.Venta, error) {
	return s.repos.Venta.FindByPeriodo(granjaID, mes, anio)
}

// Actualizar modifica una venta y ajusta el lote si corresponde
func (s *VentaService) Actualizar(id uint, input ActualizarVentaInput) (*models.Venta, error) {
	venta, err := s.repos.Venta.FindByID(id)
	if err != nil {
		return nil, ErrNotFound
	}

	cantidadAnterior := venta.Cantidad
	loteAnterior := venta.LoteID

	if input.Fecha != nil {
		fecha, err := time.Parse("2006-01-02", *input.Fecha)
		if err != nil {
			return nil, err
		}
		venta.Fecha = fecha
	}
	if input.TipoAnimal != nil {
		if !tiposAnimalValidos[*input.TipoAnimal] {
			return nil, ErrTipoAnimalInvalido
		}
		venta.TipoAnimal = *input.TipoAnimal
	}
	if input.Cantidad != nil {
		venta.Cantidad = *input.Cantidad
	}
	if input.KgTotales != nil {
		venta.KgTotales = *input.KgTotales
	}
	if input.Monto != nil {
		venta.Monto = *input.Monto
	}
	if input.Comprador != nil {
		venta.Comprador = *input.Comprador
	}
	if input.LoteID != nil {
		venta.LoteID = input.LoteID
	}
	if input.CorralID != nil {
		venta.CorralID = input.CorralID
	}
	if input.Notas != nil {
		venta.Notas = input.Notas
	}

	err = s.db.Transaction(func(tx *gorm.DB) error {
		// Ajustar lote si era/es una venta de lechones con lote
		eraLechonConLote := venta.TipoAnimal == models.TipoAnimalVentaLechon && loteAnterior != nil
		esLechonConLote := venta.TipoAnimal == models.TipoAnimalVentaLechon && venta.LoteID != nil

		if eraLechonConLote && loteAnterior != nil {
			// Revertir cantidad anterior en el lote anterior
			if err := tx.Model(&models.Lote{}).Where("id = ?", *loteAnterior).
				Update("cantidad_lechones", gorm.Expr("cantidad_lechones + ?", cantidadAnterior)).Error; err != nil {
				return err
			}
		}
		if esLechonConLote && venta.LoteID != nil {
			var lote models.Lote
			if err := tx.First(&lote, *venta.LoteID).Error; err != nil {
				return ErrNotFound
			}
			disponibles := lote.CantidadLechones
			if eraLechonConLote && loteAnterior != nil && *loteAnterior == *venta.LoteID {
				// Ya se revirtió arriba, el valor en memoria es el original + cantidadAnterior
				disponibles = lote.CantidadLechones + cantidadAnterior
			}
			if venta.Cantidad > disponibles {
				return ErrVentaExcedeLote
			}
			if err := tx.Model(&models.Lote{}).Where("id = ?", *venta.LoteID).
				Update("cantidad_lechones", gorm.Expr("cantidad_lechones - ?", venta.Cantidad)).Error; err != nil {
				return err
			}
		}

		return tx.Save(venta).Error
	})
	if err != nil {
		return nil, err
	}
	return venta, nil
}

// Eliminar borra una venta y revierte el efecto en el lote si aplica
func (s *VentaService) Eliminar(id uint) error {
	venta, err := s.repos.Venta.FindByID(id)
	if err != nil {
		return ErrNotFound
	}

	return s.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Delete(&models.Venta{}, id).Error; err != nil {
			return err
		}
		if venta.TipoAnimal == models.TipoAnimalVentaLechon && venta.LoteID != nil {
			if err := tx.Model(&models.Lote{}).Where("id = ?", *venta.LoteID).
				Update("cantidad_lechones", gorm.Expr("cantidad_lechones + ?", venta.Cantidad)).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

// GetEstadisticas obtiene estadísticas de ventas
func (s *VentaService) GetEstadisticas(granjaID *uint, mes, anio int) (map[string]interface{}, error) {
	return s.repos.Venta.GetEstadisticas(granjaID, mes, anio)
}
