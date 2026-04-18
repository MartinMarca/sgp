package repositories

import (
	"github.com/martin/sgp/internal/models"
	"gorm.io/gorm"
)

// LoteRepository maneja las operaciones de base de datos para Lotes
type LoteRepository struct {
	*BaseRepository
}

// NewLoteRepository crea una nueva instancia de LoteRepository
func NewLoteRepository(db *gorm.DB) *LoteRepository {
	return &LoteRepository{
		BaseRepository: NewBaseRepository(db),
	}
}

// Create crea un nuevo lote
func (r *LoteRepository) Create(lote *models.Lote) error {
	return r.db.Create(lote).Error
}

// FindByID busca un lote por ID
func (r *LoteRepository) FindByID(id uint, preload ...string) (*models.Lote, error) {
	var lote models.Lote
	query := r.db
	
	for _, rel := range preload {
		query = query.Preload(rel)
	}
	
	err := query.First(&lote, id).Error
	if err != nil {
		return nil, err
	}
	return &lote, nil
}

// FindByCorralID obtiene todos los lotes de un corral
func (r *LoteRepository) FindByCorralID(corralID uint, estado *string) ([]models.Lote, error) {
	var lotes []models.Lote
	query := r.db.Preload("Corral").Where("corral_id = ?", corralID)

	if estado != nil {
		query = query.Where("estado = ?", *estado)
	}

	err := query.Order("fecha_creacion DESC").Find(&lotes).Error
	return lotes, err
}

// FindByEstado obtiene todos los lotes por estado
func (r *LoteRepository) FindByEstado(estado string) ([]models.Lote, error) {
	var lotes []models.Lote
	err := r.db.
		Preload("Corral").
		Preload("Corral.Granja").
		Where("estado = ?", estado).
		Order("fecha_creacion DESC").
		Find(&lotes).Error
	return lotes, err
}

// Update actualiza un lote
func (r *LoteRepository) Update(lote *models.Lote) error {
	return r.db.Save(lote).Error
}

// Delete elimina un lote (soft delete)
func (r *LoteRepository) Delete(id uint) error {
	return r.db.Delete(&models.Lote{}, id).Error
}

// GetDestetes obtiene todos los destetes asociados a un lote
func (r *LoteRepository) GetDestetes(loteID uint) ([]models.Destete, error) {
	var destetes []models.Destete
	err := r.db.
		Preload("Cerda").
		Preload("Parto").
		Where("lote_id = ?", loteID).
		Order("fecha_destete DESC").
		Find(&destetes).Error
	return destetes, err
}

// GetCantidadTotalLechones suma la cantidad de lechones de todos los destetes del lote
func (r *LoteRepository) GetCantidadTotalLechones(loteID uint) (int, error) {
	var total int
	err := r.db.Model(&models.Destete{}).
		Select("COALESCE(SUM(cantidad_lechones_destetados), 0)").
		Where("lote_id = ?", loteID).
		Scan(&total).Error
	return total, err
}
