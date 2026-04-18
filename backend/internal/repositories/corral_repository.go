package repositories

import (
	"github.com/martin/sgp/internal/models"
	"gorm.io/gorm"
)

// CorralRepository maneja las operaciones de base de datos para Corrales
type CorralRepository struct {
	*BaseRepository
}

// NewCorralRepository crea una nueva instancia de CorralRepository
func NewCorralRepository(db *gorm.DB) *CorralRepository {
	return &CorralRepository{
		BaseRepository: NewBaseRepository(db),
	}
}

// Create crea un nuevo corral
func (r *CorralRepository) Create(corral *models.Corral) error {
	return r.db.Create(corral).Error
}

// FindByID busca un corral por ID
func (r *CorralRepository) FindByID(id uint, preload ...string) (*models.Corral, error) {
	var corral models.Corral
	query := r.db
	
	for _, rel := range preload {
		query = query.Preload(rel)
	}
	
	err := query.First(&corral, id).Error
	if err != nil {
		return nil, err
	}
	return &corral, nil
}

// FindByGranjaID obtiene todos los corrales de una granja con su ocupación actual
func (r *CorralRepository) FindByGranjaID(granjaID uint, activo *bool) ([]models.Corral, error) {
	var corrales []models.Corral
	query := r.db.Where("granja_id = ?", granjaID)

	if activo != nil {
		query = query.Where("activo = ?", *activo)
	}

	if err := query.Find(&corrales).Error; err != nil {
		return nil, err
	}

	if len(corrales) == 0 {
		return corrales, nil
	}

	corralIDs := make([]uint, len(corrales))
	for i, c := range corrales {
		corralIDs[i] = c.ID
	}

	type ocupRow struct {
		CorralID uint
		Total    int
	}
	var rows []ocupRow
	r.db.Model(&models.Lote{}).
		Select("corral_id, COALESCE(SUM(cantidad_lechones), 0) as total").
		Where("corral_id IN ? AND estado = ?", corralIDs, models.EstadoLoteActivo).
		Group("corral_id").
		Scan(&rows)

	ocupMap := make(map[uint]int, len(rows))
	for _, row := range rows {
		ocupMap[row.CorralID] = row.Total
	}

	type muerteRow struct {
		CorralID uint
		Total    int
	}
	var muerteRows []muerteRow
	r.db.Model(&models.MuerteLechon{}).
		Select("corral_id, COALESCE(SUM(cantidad), 0) as total").
		Where("corral_id IN ?", corralIDs).
		Group("corral_id").
		Scan(&muerteRows)

	muerteMap := make(map[uint]int, len(muerteRows))
	for _, row := range muerteRows {
		muerteMap[row.CorralID] = row.Total
	}

	for i, c := range corrales {
		net := ocupMap[c.ID] - muerteMap[c.ID]
		if net < 0 {
			net = 0
		}
		corrales[i].TotalAnimales = net
	}

	return corrales, nil
}

// Update actualiza un corral
func (r *CorralRepository) Update(corral *models.Corral) error {
	return r.db.Save(corral).Error
}

// Delete elimina un corral (soft delete)
func (r *CorralRepository) Delete(id uint) error {
	return r.db.Delete(&models.Corral{}, id).Error
}

// TieneLotesActivos verifica si un corral tiene lotes activos
func (r *CorralRepository) TieneLotesActivos(corralID uint) (bool, error) {
	var count int64
	err := r.db.Model(&models.Lote{}).
		Where("corral_id = ? AND estado = ?", corralID, models.EstadoLoteActivo).
		Count(&count).Error
	return count > 0, err
}

// GetOcupacion obtiene la ocupación actual del corral (animales en lotes activos menos muertes en engorde)
func (r *CorralRepository) GetOcupacion(corralID uint) (int, error) {
	var totalLotes int
	err := r.db.Model(&models.Lote{}).
		Select("COALESCE(SUM(cantidad_lechones), 0)").
		Where("corral_id = ? AND estado = ?", corralID, models.EstadoLoteActivo).
		Scan(&totalLotes).Error
	if err != nil {
		return 0, err
	}

	var totalMuertes int
	r.db.Model(&models.MuerteLechon{}).
		Select("COALESCE(SUM(cantidad), 0)").
		Where("corral_id = ?", corralID).
		Scan(&totalMuertes)

	total := totalLotes - totalMuertes
	if total < 0 {
		total = 0
	}
	return total, nil
}
