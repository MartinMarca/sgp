package repositories

import (
	"github.com/martin/sgp/internal/models"
	"gorm.io/gorm"
)

// PadrilloRepository maneja las operaciones de base de datos para Padrillos
type PadrilloRepository struct {
	*BaseRepository
}

// NewPadrilloRepository crea una nueva instancia de PadrilloRepository
func NewPadrilloRepository(db *gorm.DB) *PadrilloRepository {
	return &PadrilloRepository{
		BaseRepository: NewBaseRepository(db),
	}
}

// Create crea un nuevo padrillo
func (r *PadrilloRepository) Create(padrillo *models.Padrillo) error {
	return r.db.Create(padrillo).Error
}

// FindByID busca un padrillo por ID
func (r *PadrilloRepository) FindByID(id uint, preload ...string) (*models.Padrillo, error) {
	var padrillo models.Padrillo
	query := r.db
	
	for _, rel := range preload {
		query = query.Preload(rel)
	}
	
	err := query.First(&padrillo, id).Error
	if err != nil {
		return nil, err
	}
	return &padrillo, nil
}

// FindByGranjaID obtiene todos los padrillos de una granja
func (r *PadrilloRepository) FindByGranjaID(granjaID uint, activo *bool) ([]models.Padrillo, error) {
	var padrillos []models.Padrillo
	query := r.db.Where("granja_id = ?", granjaID)
	
	if activo != nil {
		query = query.Where("activo = ?", *activo)
	}
	
	err := query.Order("nombre").Find(&padrillos).Error
	return padrillos, err
}

// Update actualiza un padrillo
func (r *PadrilloRepository) Update(padrillo *models.Padrillo) error {
	return r.db.Save(padrillo).Error
}

// Delete elimina un padrillo (soft delete)
func (r *PadrilloRepository) Delete(id uint) error {
	return r.db.Delete(&models.Padrillo{}, id).Error
}

// ExisteCaravana verifica si ya existe una caravana en la granja
func (r *PadrilloRepository) ExisteCaravana(granjaID uint, numeroCaravana string, excludeID *uint) (bool, error) {
	query := r.db.Model(&models.Padrillo{}).
		Where("granja_id = ? AND numero_caravana = ?", granjaID, numeroCaravana)
	
	if excludeID != nil {
		query = query.Where("id != ?", *excludeID)
	}
	
	var count int64
	err := query.Count(&count).Error
	return count > 0, err
}

// GetEstadisticas obtiene estadísticas de un padrillo
func (r *PadrilloRepository) GetEstadisticas(padrilloID uint) (map[string]interface{}, error) {
	stats := make(map[string]interface{})
	
	// Total de servicios
	var totalServicios int64
	r.db.Model(&models.Servicio{}).
		Where("padrillo_id = ? AND tipo_monta = ?", padrilloID, models.TipoMontaNatural).
		Count(&totalServicios)
	stats["total_servicios"] = totalServicios
	
	// Servicios exitosos (preñez confirmada)
	var serviciosExitosos int64
	r.db.Model(&models.Servicio{}).
		Where("padrillo_id = ? AND prenez_confirmada = ?", padrilloID, true).
		Count(&serviciosExitosos)
	stats["servicios_exitosos"] = serviciosExitosos
	
	// Tasa de éxito
	tasaExito := 0.0
	if totalServicios > 0 {
		tasaExito = float64(serviciosExitosos) / float64(totalServicios) * 100
	}
	stats["tasa_exito"] = tasaExito
	
	return stats, nil
}
