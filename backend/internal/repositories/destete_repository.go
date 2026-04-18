package repositories

import (
	"time"

	"github.com/martin/sgp/internal/models"
	"gorm.io/gorm"
)

// DesteteRepository maneja las operaciones de base de datos para Destetes
type DesteteRepository struct {
	*BaseRepository
}

// NewDesteteRepository crea una nueva instancia de DesteteRepository
func NewDesteteRepository(db *gorm.DB) *DesteteRepository {
	return &DesteteRepository{
		BaseRepository: NewBaseRepository(db),
	}
}

// Create crea un nuevo destete
func (r *DesteteRepository) Create(destete *models.Destete) error {
	return r.db.Create(destete).Error
}

// FindByID busca un destete por ID
func (r *DesteteRepository) FindByID(id uint, preload ...string) (*models.Destete, error) {
	var destete models.Destete
	query := r.db
	
	for _, rel := range preload {
		query = query.Preload(rel)
	}
	
	err := query.First(&destete, id).Error
	if err != nil {
		return nil, err
	}
	return &destete, nil
}

// FindByCerdaID obtiene todos los destetes de una cerda
func (r *DesteteRepository) FindByCerdaID(cerdaID uint) ([]models.Destete, error) {
	var destetes []models.Destete
	err := r.db.
		Preload("Lote").
		Preload("Parto").
		Where("cerda_id = ?", cerdaID).
		Order("fecha_destete DESC").
		Find(&destetes).Error
	return destetes, err
}

// FindByLoteID obtiene todos los destetes de un lote
func (r *DesteteRepository) FindByLoteID(loteID uint) ([]models.Destete, error) {
	var destetes []models.Destete
	err := r.db.
		Preload("Cerda").
		Preload("Parto").
		Where("lote_id = ?", loteID).
		Order("fecha_destete DESC").
		Find(&destetes).Error
	return destetes, err
}

// FindByPeriodo obtiene destetes por período (mes/año)
func (r *DesteteRepository) FindByPeriodo(granjaID *uint, mes, anio int) ([]models.Destete, error) {
	var destetes []models.Destete
	query := r.db.
		Preload("Cerda").
		Preload("Lote").
		Preload("Lote.Corral")
	
	if granjaID != nil {
		query = query.Joins("JOIN cerdas ON cerdas.id = destetes.cerda_id").
			Where("cerdas.granja_id = ?", *granjaID)
	}
	
	if mes > 0 && anio > 0 {
		query = query.Where("MONTH(fecha_destete) = ? AND YEAR(fecha_destete) = ?", mes, anio)
	} else if anio > 0 {
		query = query.Where("YEAR(fecha_destete) = ?", anio)
	}
	
	err := query.Order("fecha_destete DESC").Find(&destetes).Error
	return destetes, err
}

// Update actualiza un destete
func (r *DesteteRepository) Update(destete *models.Destete) error {
	return r.db.Save(destete).Error
}

// Delete elimina un destete (soft delete)
func (r *DesteteRepository) Delete(id uint) error {
	return r.db.Delete(&models.Destete{}, id).Error
}

// GetDestetesFuturos obtiene destetes futuros (fecha estimada >= hoy y sin fecha_destete)
func (r *DesteteRepository) GetDestetesFuturos(granjaID *uint, diasAntes int) ([]models.Destete, error) {
	var destetes []models.Destete
	query := r.db.
		Preload("Cerda").
		Preload("Parto").
		Where("fecha_destete IS NULL OR fecha_destete = '0000-00-00'")
	
	if diasAntes > 0 {
		// Destetes estimados para los próximos N días
		fecha := time.Now().AddDate(0, 0, diasAntes)
		query = query.Where("fecha_estimada <= ?", fecha)
	}
	
	query = query.Where("fecha_estimada >= ?", time.Now())
	
	if granjaID != nil {
		query = query.Joins("JOIN cerdas ON cerdas.id = destetes.cerda_id").
			Where("cerdas.granja_id = ?", *granjaID)
	}
	
	err := query.Order("fecha_estimada ASC").Find(&destetes).Error
	return destetes, err
}

// GetEstadisticas obtiene estadísticas de destetes
func (r *DesteteRepository) GetEstadisticas(granjaID *uint, mes, anio int) (map[string]interface{}, error) {
	stats := make(map[string]interface{})
	
	query := r.db.Model(&models.Destete{})
	
	if granjaID != nil {
		query = query.Joins("JOIN cerdas ON cerdas.id = destetes.cerda_id").
			Where("cerdas.granja_id = ?", *granjaID)
	}
	
	if mes > 0 && anio > 0 {
		query = query.Where("MONTH(fecha_destete) = ? AND YEAR(fecha_destete) = ?", mes, anio)
	} else if anio > 0 {
		query = query.Where("YEAR(fecha_destete) = ?", anio)
	}
	
	// Total de destetes
	var totalDestetes int64
	query.Count(&totalDestetes)
	stats["total_destetes"] = totalDestetes
	
	// Promedio de lechones destetados
	var promedioDestetados float64
	query.Select("COALESCE(AVG(cantidad_lechones_destetados), 0)").Scan(&promedioDestetados)
	stats["promedio_lechones_destetados"] = promedioDestetados
	
	// Total de lechones destetados
	var totalLechones int
	query.Select("COALESCE(SUM(cantidad_lechones_destetados), 0)").Scan(&totalLechones)
	stats["total_lechones_destetados"] = totalLechones
	
	return stats, nil
}
