package repositories

import (
	"github.com/martin/sgp/internal/models"
	"gorm.io/gorm"
)

// VentaRepository maneja las operaciones de base de datos para Ventas
type VentaRepository struct {
	*BaseRepository
}

// NewVentaRepository crea una nueva instancia de VentaRepository
func NewVentaRepository(db *gorm.DB) *VentaRepository {
	return &VentaRepository{BaseRepository: NewBaseRepository(db)}
}

func (r *VentaRepository) Create(venta *models.Venta) error {
	return r.db.Create(venta).Error
}

func (r *VentaRepository) FindByID(id uint, preload ...string) (*models.Venta, error) {
	var venta models.Venta
	q := r.db
	for _, rel := range preload {
		q = q.Preload(rel)
	}
	if err := q.First(&venta, id).Error; err != nil {
		return nil, err
	}
	return &venta, nil
}

func (r *VentaRepository) FindByPeriodo(granjaID *uint, mes, anio int) ([]models.Venta, error) {
	var ventas []models.Venta
	q := r.db.Preload("Lote").Preload("Corral")
	if granjaID != nil {
		q = q.Where("granja_id = ?", *granjaID)
	}
	if mes > 0 && anio > 0 {
		q = q.Where("MONTH(fecha) = ? AND YEAR(fecha) = ?", mes, anio)
	} else if anio > 0 {
		q = q.Where("YEAR(fecha) = ?", anio)
	}
	err := q.Order("fecha DESC").Find(&ventas).Error
	return ventas, err
}

func (r *VentaRepository) Update(venta *models.Venta) error {
	return r.db.Save(venta).Error
}

func (r *VentaRepository) Delete(id uint) error {
	return r.db.Delete(&models.Venta{}, id).Error
}

func (r *VentaRepository) GetEstadisticas(granjaID *uint, mes, anio int) (map[string]interface{}, error) {
	stats := make(map[string]interface{})

	q := func() *gorm.DB {
		base := r.db.Model(&models.Venta{})
		if granjaID != nil {
			base = base.Where("granja_id = ?", *granjaID)
		}
		if mes > 0 && anio > 0 {
			base = base.Where("MONTH(fecha) = ? AND YEAR(fecha) = ?", mes, anio)
		} else if anio > 0 {
			base = base.Where("YEAR(fecha) = ?", anio)
		}
		return base
	}

	var totalVentas int64
	q().Count(&totalVentas)
	stats["total_ventas"] = totalVentas

	var totalAnimales int
	q().Select("COALESCE(SUM(cantidad), 0)").Scan(&totalAnimales)
	stats["total_animales"] = totalAnimales

	var totalKg float64
	q().Select("COALESCE(SUM(kg_totales), 0)").Scan(&totalKg)
	stats["total_kg"] = totalKg

	var totalMonto float64
	q().Select("COALESCE(SUM(monto), 0)").Scan(&totalMonto)
	stats["total_monto"] = totalMonto

	// Desglose por tipo de animal
	type resRow struct {
		TipoAnimal string
		Cantidad   int
		KgTotales  float64
		Monto      float64
	}
	var rows []resRow
	q().Select("tipo_animal, SUM(cantidad) as cantidad, SUM(kg_totales) as kg_totales, SUM(monto) as monto").
		Group("tipo_animal").Scan(&rows)

	porTipo := make([]map[string]interface{}, 0)
	for _, r := range rows {
		porTipo = append(porTipo, map[string]interface{}{
			"tipo_animal": r.TipoAnimal,
			"cantidad":    r.Cantidad,
			"kg_totales":  r.KgTotales,
			"monto":       r.Monto,
		})
	}
	stats["por_tipo"] = porTipo

	return stats, nil
}
