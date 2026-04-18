package repositories

import (
	"time"

	"github.com/martin/sgp/internal/models"
	"gorm.io/gorm"
)

// PartoRepository maneja las operaciones de base de datos para Partos
type PartoRepository struct {
	*BaseRepository
}

// NewPartoRepository crea una nueva instancia de PartoRepository
func NewPartoRepository(db *gorm.DB) *PartoRepository {
	return &PartoRepository{
		BaseRepository: NewBaseRepository(db),
	}
}

// Create crea un nuevo parto
func (r *PartoRepository) Create(parto *models.Parto) error {
	return r.db.Create(parto).Error
}

// FindByID busca un parto por ID
func (r *PartoRepository) FindByID(id uint, preload ...string) (*models.Parto, error) {
	var parto models.Parto
	query := r.db
	
	for _, rel := range preload {
		query = query.Preload(rel)
	}
	
	err := query.First(&parto, id).Error
	if err != nil {
		return nil, err
	}
	return &parto, nil
}

// FindByCerdaID obtiene todos los partos de una cerda
func (r *PartoRepository) FindByCerdaID(cerdaID uint) ([]models.Parto, error) {
	var partos []models.Parto
	err := r.db.
		Where("cerda_id = ?", cerdaID).
		Order("fecha_parto DESC").
		Find(&partos).Error
	return partos, err
}

// FindByPeriodo obtiene partos por período (mes/año)
func (r *PartoRepository) FindByPeriodo(granjaID *uint, mes, anio int) ([]models.Parto, error) {
	var partos []models.Parto
	query := r.db.Preload("Cerda")
	
	if granjaID != nil {
		query = query.Joins("JOIN cerdas ON cerdas.id = partos.cerda_id").
			Where("cerdas.granja_id = ?", *granjaID)
	}
	
	if mes > 0 && anio > 0 {
		query = query.Where("MONTH(fecha_parto) = ? AND YEAR(fecha_parto) = ?", mes, anio)
	} else if anio > 0 {
		query = query.Where("YEAR(fecha_parto) = ?", anio)
	}
	
	err := query.Order("fecha_parto DESC").Find(&partos).Error
	return partos, err
}

// Update actualiza un parto
func (r *PartoRepository) Update(parto *models.Parto) error {
	return r.db.Save(parto).Error
}

// Delete elimina un parto (soft delete)
func (r *PartoRepository) Delete(id uint) error {
	return r.db.Delete(&models.Parto{}, id).Error
}

// GetPartosFuturos obtiene partos futuros (fecha estimada >= hoy y sin fecha_parto)
func (r *PartoRepository) GetPartosFuturos(granjaID *uint, diasAntes int) ([]models.Parto, error) {
	var partos []models.Parto
	query := r.db.
		Preload("Cerda").
		Preload("Servicio").
		Where("fecha_parto IS NULL OR fecha_parto = '0000-00-00'")
	
	if diasAntes > 0 {
		// Partos estimados para los próximos N días
		fecha := time.Now().AddDate(0, 0, diasAntes)
		query = query.Where("fecha_estimada <= ?", fecha)
	}
	
	query = query.Where("fecha_estimada >= ?", time.Now())
	
	if granjaID != nil {
		query = query.Joins("JOIN cerdas ON cerdas.id = partos.cerda_id").
			Where("cerdas.granja_id = ?", *granjaID)
	}
	
	err := query.Order("fecha_estimada ASC").Find(&partos).Error
	return partos, err
}

// TieneDestete verifica si un parto ya tiene un destete asociado
func (r *PartoRepository) TieneDestete(partoID uint) (bool, error) {
	var count int64
	err := r.db.Model(&models.Destete{}).
		Where("parto_id = ?", partoID).
		Count(&count).Error
	return count > 0, err
}

// GetEstadisticas obtiene estadísticas de partos
func (r *PartoRepository) GetEstadisticas(granjaID *uint, mes, anio int) (map[string]interface{}, error) {
	stats := make(map[string]interface{})
	
	query := r.db.Model(&models.Parto{})
	
	if granjaID != nil {
		query = query.Joins("JOIN cerdas ON cerdas.id = partos.cerda_id").
			Where("cerdas.granja_id = ?", *granjaID)
	}
	
	if mes > 0 && anio > 0 {
		query = query.Where("MONTH(fecha_parto) = ? AND YEAR(fecha_parto) = ?", mes, anio)
	} else if anio > 0 {
		query = query.Where("YEAR(fecha_parto) = ?", anio)
	}
	
	// Total de partos
	var totalPartos int64
	query.Count(&totalPartos)
	stats["total_partos"] = totalPartos
	
	// Promedio de lechones nacidos vivos
	var promedioVivos float64
	query.Select("COALESCE(AVG(lechones_nacidos_vivos), 0)").Scan(&promedioVivos)
	stats["promedio_lechones_vivos"] = promedioVivos
	
	// Promedio de lechones totales
	var promedioTotales float64
	query.Select("COALESCE(AVG(lechones_nacidos_totales), 0)").Scan(&promedioTotales)
	stats["promedio_lechones_totales"] = promedioTotales
	
	// Total de lechones nacidos
	var totalLechones int
	query.Select("COALESCE(SUM(lechones_nacidos_vivos), 0)").Scan(&totalLechones)
	stats["total_lechones_nacidos"] = totalLechones
	
	return stats, nil
}
