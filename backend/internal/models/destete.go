package models

import (
	"time"

	"gorm.io/gorm"
)

// Destete representa un destete de una cerda
type Destete struct {
	ID                        uint           `gorm:"primaryKey" json:"id"`
	CerdaID                   uint           `gorm:"not null;index:idx_destetes_cerda_fecha" json:"cerda_id" binding:"required"`
	PartoID                   *uint          `gorm:"index" json:"parto_id,omitempty"`
	FechaDestete              time.Time      `gorm:"type:date;not null;index:idx_destetes_cerda_fecha" json:"fecha_destete" binding:"required"`
	CantidadLechonesDestetados int           `gorm:"not null;default:0" json:"cantidad_lechones_destetados" binding:"required,gte=0"`
	FechaEstimada             time.Time      `gorm:"type:date;not null;index" json:"fecha_estimada" binding:"required"`
	LoteID                    uint           `gorm:"not null;index" json:"lote_id" binding:"required"` // Obligatorio: los lechones deben estar en un lote
	CreatedAt                 time.Time      `json:"created_at"`
	UpdatedAt                 time.Time      `json:"updated_at"`
	DeletedAt                 gorm.DeletedAt `gorm:"index" json:"-"`

	// Relaciones
	Cerda Cerda  `gorm:"foreignKey:CerdaID;constraint:OnDelete:RESTRICT" json:"cerda,omitempty"`
	Parto *Parto `gorm:"foreignKey:PartoID;constraint:OnDelete:SET NULL" json:"parto,omitempty"`
	Lote  Lote   `gorm:"foreignKey:LoteID;constraint:OnDelete:RESTRICT" json:"lote,omitempty"`
}

// TableName especifica el nombre de la tabla
func (Destete) TableName() string {
	return "destetes"
}

// BeforeCreate hook de GORM para validaciones antes de crear
func (d *Destete) BeforeCreate(tx *gorm.DB) error {
	return d.validate(tx)
}

// BeforeUpdate hook de GORM para validaciones antes de actualizar
func (d *Destete) BeforeUpdate(tx *gorm.DB) error {
	return d.validate(tx)
}

// validate valida las reglas de negocio del destete
func (d *Destete) validate(tx *gorm.DB) error {
	// Si hay parto asociado, validar que cantidad_lechones_destetados <= lechones_nacidos_vivos
	if d.PartoID != nil {
		var parto Parto
		if err := tx.First(&parto, *d.PartoID).Error; err != nil {
			return err
		}
		if d.CantidadLechonesDestetados > parto.LechonesNacidosVivos {
			return gorm.ErrInvalidValue
		}
	}
	return nil
}
