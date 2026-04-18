package repositories

import (
	"github.com/martin/sgp/internal/models"
	"gorm.io/gorm"
)

// GranjaRepository maneja las operaciones de base de datos para Granjas
type GranjaRepository struct {
	*BaseRepository
}

// NewGranjaRepository crea una nueva instancia de GranjaRepository
func NewGranjaRepository(db *gorm.DB) *GranjaRepository {
	return &GranjaRepository{
		BaseRepository: NewBaseRepository(db),
	}
}

// Create crea una nueva granja
func (r *GranjaRepository) Create(granja *models.Granja) error {
	return r.db.Create(granja).Error
}

// FindByID busca una granja por ID con todas sus relaciones
func (r *GranjaRepository) FindByID(id uint, preload ...string) (*models.Granja, error) {
	var granja models.Granja
	query := r.db
	
	// Preload relaciones si se especifican
	for _, rel := range preload {
		query = query.Preload(rel)
	}
	
	err := query.First(&granja, id).Error
	if err != nil {
		return nil, err
	}
	return &granja, nil
}

// FindAll obtiene todas las granjas activas
func (r *GranjaRepository) FindAll(activo *bool) ([]models.Granja, error) {
	var granjas []models.Granja
	query := r.db
	
	if activo != nil {
		query = query.Where("activo = ?", *activo)
	}
	
	err := query.Find(&granjas).Error
	return granjas, err
}

// Update actualiza una granja
func (r *GranjaRepository) Update(granja *models.Granja) error {
	return r.db.Save(granja).Error
}

// Delete elimina una granja (soft delete)
func (r *GranjaRepository) Delete(id uint) error {
	return r.db.Delete(&models.Granja{}, id).Error
}

// FindByUsuarioID obtiene todas las granjas de un usuario
func (r *GranjaRepository) FindByUsuarioID(usuarioID uint) ([]models.Granja, error) {
	var granjas []models.Granja
	err := r.db.
		Joins("JOIN usuario_granja ON usuario_granja.granja_id = granjas.id").
		Where("usuario_granja.usuario_id = ?", usuarioID).
		Find(&granjas).Error
	return granjas, err
}

// AsignarUsuario asigna un usuario a una granja con un rol
func (r *GranjaRepository) AsignarUsuario(granjaID, usuarioID uint, rol string) error {
	usuarioGranja := models.UsuarioGranja{
		GranjaID:  granjaID,
		UsuarioID: usuarioID,
		Rol:       rol,
	}
	return r.db.Create(&usuarioGranja).Error
}

// RemoverUsuario remueve un usuario de una granja
func (r *GranjaRepository) RemoverUsuario(granjaID, usuarioID uint) error {
	return r.db.
		Where("granja_id = ? AND usuario_id = ?", granjaID, usuarioID).
		Delete(&models.UsuarioGranja{}).Error
}

// GetEstadisticas obtiene estadísticas básicas de una granja
func (r *GranjaRepository) GetEstadisticas(granjaID uint) (map[string]interface{}, error) {
	stats := make(map[string]interface{})
	
	// Contar corrales
	var corralesCount int64
	r.db.Model(&models.Corral{}).Where("granja_id = ? AND activo = ?", granjaID, true).Count(&corralesCount)
	stats["corrales"] = corralesCount
	
	// Contar cerdas por estado
	var cerdas []struct {
		Estado string
		Count  int64
	}
	r.db.Model(&models.Cerda{}).
		Select("estado, COUNT(*) as count").
		Where("granja_id = ? AND activo = ?", granjaID, true).
		Group("estado").
		Scan(&cerdas)
	
	cerdasPorEstado := make(map[string]int64)
	for _, c := range cerdas {
		cerdasPorEstado[c.Estado] = c.Count
	}
	stats["cerdas_por_estado"] = cerdasPorEstado
	
	// Contar padrillos
	var padrillosCount int64
	r.db.Model(&models.Padrillo{}).Where("granja_id = ? AND activo = ?", granjaID, true).Count(&padrillosCount)
	stats["padrillos"] = padrillosCount
	
	return stats, nil
}
