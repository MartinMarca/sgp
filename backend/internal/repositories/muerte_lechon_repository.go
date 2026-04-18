package repositories

import (
	"github.com/martin/sgp/internal/models"
	"gorm.io/gorm"
)

// MuerteLechonRepository maneja las operaciones de base de datos para MuerteLechones
type MuerteLechonRepository struct {
	*BaseRepository
}

// NewMuerteLechonRepository crea una nueva instancia de MuerteLechonRepository
func NewMuerteLechonRepository(db *gorm.DB) *MuerteLechonRepository {
	return &MuerteLechonRepository{
		BaseRepository: NewBaseRepository(db),
	}
}

// Create crea un nuevo registro de muerte de lechón
func (r *MuerteLechonRepository) Create(muerte *models.MuerteLechon) error {
	return r.db.Create(muerte).Error
}

// FindByID busca una muerte de lechón por ID
func (r *MuerteLechonRepository) FindByID(id uint, preload ...string) (*models.MuerteLechon, error) {
	var muerte models.MuerteLechon
	query := r.db

	for _, rel := range preload {
		query = query.Preload(rel)
	}

	err := query.First(&muerte, id).Error
	if err != nil {
		return nil, err
	}
	return &muerte, nil
}

// FindByPartoID obtiene todas las muertes asociadas a un parto (lactancia)
func (r *MuerteLechonRepository) FindByPartoID(partoID uint) ([]models.MuerteLechon, error) {
	var muertes []models.MuerteLechon
	err := r.db.
		Where("parto_id = ?", partoID).
		Order("fecha DESC").
		Find(&muertes).Error
	return muertes, err
}

// FindByCorralID obtiene todas las muertes asociadas a un corral (engorde)
func (r *MuerteLechonRepository) FindByCorralID(corralID uint) ([]models.MuerteLechon, error) {
	var muertes []models.MuerteLechon
	err := r.db.
		Where("corral_id = ?", corralID).
		Order("fecha DESC").
		Find(&muertes).Error
	return muertes, err
}

// FindByGranjaID obtiene todas las muertes de una granja
func (r *MuerteLechonRepository) FindByGranjaID(granjaID uint) ([]models.MuerteLechon, error) {
	var muertes []models.MuerteLechon
	err := r.db.
		Preload("Parto").
		Preload("Parto.Cerda").
		Preload("Corral").
		Where("granja_id = ?", granjaID).
		Order("fecha DESC").
		Find(&muertes).Error
	return muertes, err
}

// FindByPeriodo obtiene muertes por período (mes/año) y opcionalmente por granja
func (r *MuerteLechonRepository) FindByPeriodo(granjaID *uint, mes, anio int) ([]models.MuerteLechon, error) {
	var muertes []models.MuerteLechon
	query := r.db.
		Preload("Parto").
		Preload("Parto.Cerda").
		Preload("Corral")

	if granjaID != nil {
		query = query.Where("granja_id = ?", *granjaID)
	}

	if mes > 0 && anio > 0 {
		query = query.Where("MONTH(fecha) = ? AND YEAR(fecha) = ?", mes, anio)
	} else if anio > 0 {
		query = query.Where("YEAR(fecha) = ?", anio)
	}

	err := query.Order("fecha DESC").Find(&muertes).Error
	return muertes, err
}

// Update actualiza un registro de muerte
func (r *MuerteLechonRepository) Update(muerte *models.MuerteLechon) error {
	return r.db.Save(muerte).Error
}

// Delete elimina un registro de muerte (soft delete)
func (r *MuerteLechonRepository) Delete(id uint) error {
	return r.db.Delete(&models.MuerteLechon{}, id).Error
}

// SumByPartoID suma la cantidad total de muertes para un parto
func (r *MuerteLechonRepository) SumByPartoID(partoID uint) (int, error) {
	var total int
	err := r.db.Model(&models.MuerteLechon{}).
		Select("COALESCE(SUM(cantidad), 0)").
		Where("parto_id = ?", partoID).
		Scan(&total).Error
	return total, err
}

// SumByCorralID suma la cantidad total de muertes para un corral
func (r *MuerteLechonRepository) SumByCorralID(corralID uint) (int, error) {
	var total int
	err := r.db.Model(&models.MuerteLechon{}).
		Select("COALESCE(SUM(cantidad), 0)").
		Where("corral_id = ?", corralID).
		Scan(&total).Error
	return total, err
}

// GetEstadisticas obtiene estadísticas de mortalidad
func (r *MuerteLechonRepository) GetEstadisticas(granjaID *uint, mes, anio int) (map[string]interface{}, error) {
	stats := make(map[string]interface{})

	query := r.db.Model(&models.MuerteLechon{})

	if granjaID != nil {
		query = query.Where("granja_id = ?", *granjaID)
	}

	if mes > 0 && anio > 0 {
		query = query.Where("MONTH(fecha) = ? AND YEAR(fecha) = ?", mes, anio)
	} else if anio > 0 {
		query = query.Where("YEAR(fecha) = ?", anio)
	}

	var totalRegistros int64
	query.Count(&totalRegistros)
	stats["total_registros"] = totalRegistros

	var totalMuertes int
	query.Select("COALESCE(SUM(cantidad), 0)").Scan(&totalMuertes)
	stats["total_muertes"] = totalMuertes

	// Muertes en lactancia
	var muertesLactancia int
	r.db.Model(&models.MuerteLechon{}).
		Where(r.buildFilters(granjaID, mes, anio)).
		Where("parto_id IS NOT NULL").
		Select("COALESCE(SUM(cantidad), 0)").
		Scan(&muertesLactancia)
	stats["muertes_lactancia"] = muertesLactancia

	// Muertes en engorde
	var muertesEngorde int
	r.db.Model(&models.MuerteLechon{}).
		Where(r.buildFilters(granjaID, mes, anio)).
		Where("corral_id IS NOT NULL").
		Select("COALESCE(SUM(cantidad), 0)").
		Scan(&muertesEngorde)
	stats["muertes_engorde"] = muertesEngorde

	// Muertes por causa
	type causaCount struct {
		Causa    string `json:"causa"`
		Cantidad int    `json:"cantidad"`
	}
	var porCausa []causaCount
	r.db.Model(&models.MuerteLechon{}).
		Where(r.buildFilters(granjaID, mes, anio)).
		Select("causa, COALESCE(SUM(cantidad), 0) as cantidad").
		Group("causa").
		Order("cantidad DESC").
		Scan(&porCausa)
	stats["muertes_por_causa"] = porCausa

	return stats, nil
}

// buildFilters construye las condiciones de filtro reutilizables
func (r *MuerteLechonRepository) buildFilters(granjaID *uint, mes, anio int) *gorm.DB {
	q := r.db.Where("1 = 1")

	if granjaID != nil {
		q = q.Where("granja_id = ?", *granjaID)
	}
	if mes > 0 && anio > 0 {
		q = q.Where("MONTH(fecha) = ? AND YEAR(fecha) = ?", mes, anio)
	} else if anio > 0 {
		q = q.Where("YEAR(fecha) = ?", anio)
	}

	return q
}
