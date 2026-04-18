package repositories

import (
	"github.com/martin/sgp/internal/models"
	"gorm.io/gorm"
)

// CerdaRepository maneja las operaciones de base de datos para Cerdas
type CerdaRepository struct {
	*BaseRepository
}

// NewCerdaRepository crea una nueva instancia de CerdaRepository
func NewCerdaRepository(db *gorm.DB) *CerdaRepository {
	return &CerdaRepository{
		BaseRepository: NewBaseRepository(db),
	}
}

// Create crea una nueva cerda
func (r *CerdaRepository) Create(cerda *models.Cerda) error {
	return r.db.Create(cerda).Error
}

// FindByID busca una cerda por ID
func (r *CerdaRepository) FindByID(id uint, preload ...string) (*models.Cerda, error) {
	var cerda models.Cerda
	query := r.db
	
	for _, rel := range preload {
		query = query.Preload(rel)
	}
	
	err := query.First(&cerda, id).Error
	if err != nil {
		return nil, err
	}
	return &cerda, nil
}

// FindByGranjaID obtiene todas las cerdas de una granja
func (r *CerdaRepository) FindByGranjaID(granjaID uint, estado *string, activo *bool) ([]models.Cerda, error) {
	var cerdas []models.Cerda
	query := r.db.Where("granja_id = ?", granjaID)
	
	if estado != nil {
		query = query.Where("estado = ?", *estado)
	}
	
	if activo != nil {
		query = query.Where("activo = ?", *activo)
	}
	
	err := query.Order("numero_caravana").Find(&cerdas).Error
	return cerdas, err
}

// FindByEstado obtiene todas las cerdas por estado
func (r *CerdaRepository) FindByEstado(estado string, granjaID *uint) ([]models.Cerda, error) {
	var cerdas []models.Cerda
	query := r.db.Where("estado = ? AND activo = ?", estado, true)
	
	if granjaID != nil {
		query = query.Where("granja_id = ?", *granjaID)
	}
	
	err := query.Order("numero_caravana").Find(&cerdas).Error
	return cerdas, err
}

// Update actualiza una cerda
func (r *CerdaRepository) Update(cerda *models.Cerda) error {
	return r.db.Save(cerda).Error
}

// Delete elimina una cerda (soft delete)
func (r *CerdaRepository) Delete(id uint) error {
	return r.db.Delete(&models.Cerda{}, id).Error
}

// CambiarEstado cambia el estado de una cerda
func (r *CerdaRepository) CambiarEstado(cerdaID uint, nuevoEstado string) error {
	return r.db.Model(&models.Cerda{}).
		Where("id = ?", cerdaID).
		Update("estado", nuevoEstado).Error
}

// ExisteCaravana verifica si ya existe una caravana en la granja
func (r *CerdaRepository) ExisteCaravana(granjaID uint, numeroCaravana string, excludeID *uint) (bool, error) {
	query := r.db.Model(&models.Cerda{}).
		Where("granja_id = ? AND numero_caravana = ?", granjaID, numeroCaravana)
	
	if excludeID != nil {
		query = query.Where("id != ?", *excludeID)
	}
	
	var count int64
	err := query.Count(&count).Error
	return count > 0, err
}

// GetHistorial obtiene el historial completo de una cerda
func (r *CerdaRepository) GetHistorial(cerdaID uint) (map[string]interface{}, error) {
	historial := make(map[string]interface{})
	
	// Servicios
	var servicios []models.Servicio
	r.db.Where("cerda_id = ?", cerdaID).
		Preload("Padrillo").
		Order("fecha_servicio DESC").
		Find(&servicios)
	historial["servicios"] = servicios
	
	// Partos
	var partos []models.Parto
	r.db.Where("cerda_id = ?", cerdaID).
		Order("fecha_parto DESC").
		Find(&partos)
	historial["partos"] = partos
	
	// Destetes
	var destetes []models.Destete
	r.db.Where("cerda_id = ?", cerdaID).
		Preload("Lote").
		Order("fecha_destete DESC").
		Find(&destetes)
	historial["destetes"] = destetes
	
	return historial, nil
}

// GetEstadisticas obtiene estadísticas de una cerda
func (r *CerdaRepository) GetEstadisticas(cerdaID uint) (map[string]interface{}, error) {
	stats := make(map[string]interface{})
	
	// Total de servicios
	var totalServicios int64
	r.db.Model(&models.Servicio{}).Where("cerda_id = ?", cerdaID).Count(&totalServicios)
	stats["total_servicios"] = totalServicios
	
	// Servicios exitosos (preñez confirmada)
	var serviciosExitosos int64
	r.db.Model(&models.Servicio{}).
		Where("cerda_id = ? AND prenez_confirmada = ?", cerdaID, true).
		Count(&serviciosExitosos)
	stats["servicios_exitosos"] = serviciosExitosos
	
	// Total de partos
	var totalPartos int64
	r.db.Model(&models.Parto{}).Where("cerda_id = ?", cerdaID).Count(&totalPartos)
	stats["total_partos"] = totalPartos
	
	// Promedio de lechones por parto
	var promedioLechones float64
	r.db.Model(&models.Parto{}).
		Select("COALESCE(AVG(lechones_nacidos_vivos), 0)").
		Where("cerda_id = ?", cerdaID).
		Scan(&promedioLechones)
	stats["promedio_lechones"] = promedioLechones
	
	// Total de destetes
	var totalDestetes int64
	r.db.Model(&models.Destete{}).Where("cerda_id = ?", cerdaID).Count(&totalDestetes)
	stats["total_destetes"] = totalDestetes
	
	return stats, nil
}

// GetUltimoServicio obtiene el último servicio de una cerda
func (r *CerdaRepository) GetUltimoServicio(cerdaID uint) (*models.Servicio, error) {
	var servicio models.Servicio
	err := r.db.Where("cerda_id = ?", cerdaID).
		Order("fecha_servicio DESC").
		First(&servicio).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &servicio, nil
}

// GetUltimoParto obtiene el último parto de una cerda
func (r *CerdaRepository) GetUltimoParto(cerdaID uint) (*models.Parto, error) {
	var parto models.Parto
	err := r.db.Where("cerda_id = ?", cerdaID).
		Order("fecha_parto DESC").
		First(&parto).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &parto, nil
}

// GetPartoSinDestete obtiene el parto sin destete de una cerda en cría
func (r *CerdaRepository) GetPartoSinDestete(cerdaID uint) (*models.Parto, error) {
	var parto models.Parto
	err := r.db.Raw(`
		SELECT p.* FROM partos p
		WHERE p.cerda_id = ?
		AND NOT EXISTS (
			SELECT 1 FROM destetes d WHERE d.parto_id = p.id
		)
		ORDER BY p.fecha_parto DESC
		LIMIT 1
	`, cerdaID).Scan(&parto).Error
	
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	
	if parto.ID == 0 {
		return nil, nil
	}
	
	return &parto, nil
}
