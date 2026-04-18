package repositories

import (
	"time"

	"github.com/martin/sgp/internal/models"
	"gorm.io/gorm"
)

// ServicioRepository maneja las operaciones de base de datos para Servicios
type ServicioRepository struct {
	*BaseRepository
}

// NewServicioRepository crea una nueva instancia de ServicioRepository
func NewServicioRepository(db *gorm.DB) *ServicioRepository {
	return &ServicioRepository{
		BaseRepository: NewBaseRepository(db),
	}
}

// Create crea un nuevo servicio
func (r *ServicioRepository) Create(servicio *models.Servicio) error {
	return r.db.Create(servicio).Error
}

// FindByID busca un servicio por ID
func (r *ServicioRepository) FindByID(id uint, preload ...string) (*models.Servicio, error) {
	var servicio models.Servicio
	query := r.db
	
	for _, rel := range preload {
		query = query.Preload(rel)
	}
	
	err := query.First(&servicio, id).Error
	if err != nil {
		return nil, err
	}
	return &servicio, nil
}

// FindByCerdaID obtiene todos los servicios de una cerda
func (r *ServicioRepository) FindByCerdaID(cerdaID uint) ([]models.Servicio, error) {
	var servicios []models.Servicio
	err := r.db.
		Preload("Padrillo").
		Where("cerda_id = ?", cerdaID).
		Order("fecha_servicio DESC").
		Find(&servicios).Error
	return servicios, err
}

// FindByPeriodo obtiene servicios por período (mes/año)
func (r *ServicioRepository) FindByPeriodo(granjaID *uint, mes, anio int) ([]models.Servicio, error) {
	var servicios []models.Servicio
	query := r.db.
		Preload("Cerda").
		Preload("Padrillo")
	
	if granjaID != nil {
		query = query.Joins("JOIN cerdas ON cerdas.id = servicios.cerda_id").
			Where("cerdas.granja_id = ?", *granjaID)
	}
	
	if mes > 0 && anio > 0 {
		query = query.Where("MONTH(fecha_servicio) = ? AND YEAR(fecha_servicio) = ?", mes, anio)
	} else if anio > 0 {
		query = query.Where("YEAR(fecha_servicio) = ?", anio)
	}
	
	err := query.Order("fecha_servicio DESC").Find(&servicios).Error
	return servicios, err
}

// Update actualiza un servicio
func (r *ServicioRepository) Update(servicio *models.Servicio) error {
	return r.db.Save(servicio).Error
}

// Delete elimina un servicio (soft delete)
func (r *ServicioRepository) Delete(id uint) error {
	return r.db.Delete(&models.Servicio{}, id).Error
}

// ConfirmarPrenez confirma la preñez de un servicio
func (r *ServicioRepository) ConfirmarPrenez(servicioID uint, fechaConfirmacion time.Time) error {
	return r.db.Model(&models.Servicio{}).
		Where("id = ?", servicioID).
		Updates(map[string]interface{}{
			"prenez_confirmada":          true,
			"fecha_confirmacion_prenez":  fechaConfirmacion,
		}).Error
}

// CancelarPrenez cancela la preñez de un servicio
func (r *ServicioRepository) CancelarPrenez(servicioID uint, fechaCancelacion time.Time, motivo string) error {
	return r.db.Model(&models.Servicio{}).
		Where("id = ?", servicioID).
		Updates(map[string]interface{}{
			"prenez_cancelada":         true,
			"fecha_cancelacion_prenez": fechaCancelacion,
			"motivo_cancelacion":       motivo,
		}).Error
}

// GetServicioConPrenezConfirmada obtiene el servicio con preñez confirmada de una cerda
func (r *ServicioRepository) GetServicioConPrenezConfirmada(cerdaID uint) (*models.Servicio, error) {
	var servicio models.Servicio
	err := r.db.
		Where("cerda_id = ? AND prenez_confirmada = ? AND prenez_cancelada = ?", 
			cerdaID, true, false).
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

// GetEstadisticas obtiene estadísticas de servicios para un período
func (r *ServicioRepository) GetEstadisticas(granjaID *uint, mes, anio int) (map[string]interface{}, error) {
	stats := make(map[string]interface{})

	baseQuery := func() *gorm.DB {
		q := r.db.Model(&models.Servicio{})
		if granjaID != nil {
			q = q.Joins("JOIN cerdas ON cerdas.id = servicios.cerda_id").
				Where("cerdas.granja_id = ?", *granjaID)
		}
		return q
	}

	// Total de servicios en el período (por fecha_servicio)
	servicioQ := baseQuery()
	if mes > 0 && anio > 0 {
		servicioQ = servicioQ.Where("MONTH(servicios.fecha_servicio) = ? AND YEAR(servicios.fecha_servicio) = ?", mes, anio)
	} else if anio > 0 {
		servicioQ = servicioQ.Where("YEAR(servicios.fecha_servicio) = ?", anio)
	}
	var totalServicios int64
	servicioQ.Count(&totalServicios)
	stats["total_servicios"] = totalServicios

	// Total de confirmaciones de preñez en el período (por fecha_confirmacion_prenez)
	confirmQ := baseQuery().Where("servicios.prenez_confirmada = ? AND servicios.prenez_cancelada = ?", true, false)
	if mes > 0 && anio > 0 {
		confirmQ = confirmQ.Where("MONTH(servicios.fecha_confirmacion_prenez) = ? AND YEAR(servicios.fecha_confirmacion_prenez) = ?", mes, anio)
	} else if anio > 0 {
		confirmQ = confirmQ.Where("YEAR(servicios.fecha_confirmacion_prenez) = ?", anio)
	}
	var totalConfirmaciones int64
	confirmQ.Count(&totalConfirmaciones)
	stats["total_confirmaciones"] = totalConfirmaciones

	return stats, nil
}

// GetServiciosPendientesConfirmacion obtiene servicios pendientes de confirmación
func (r *ServicioRepository) GetServiciosPendientesConfirmacion(granjaID *uint) ([]models.Servicio, error) {
	var servicios []models.Servicio
	query := r.db.
		Preload("Cerda").
		Preload("Padrillo").
		Where("prenez_confirmada = ? AND prenez_cancelada = ?", false, false)
	
	if granjaID != nil {
		query = query.Joins("JOIN cerdas ON cerdas.id = servicios.cerda_id").
			Where("cerdas.granja_id = ?", *granjaID)
	}
	
	err := query.Order("fecha_servicio DESC").Find(&servicios).Error
	return servicios, err
}
