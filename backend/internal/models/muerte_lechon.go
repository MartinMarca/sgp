package models

import (
	"time"

	"gorm.io/gorm"
)

// MuerteLechon representa una muerte de animales, ya sea en lactancia (asociada a parto) o en engorde (asociada a corral)
type MuerteLechon struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	GranjaID  uint           `gorm:"not null;index" json:"granja_id" binding:"required"`
	PartoID   *uint          `gorm:"index" json:"parto_id,omitempty"`
	CorralID  *uint          `gorm:"index" json:"corral_id,omitempty"`
	Fecha     time.Time      `gorm:"type:date;not null;index" json:"fecha" binding:"required"`
	Cantidad  int            `gorm:"not null" json:"cantidad" binding:"required,gte=1"`
	Causa     string         `gorm:"type:enum('aplastamiento','enfermedad','inanicion','otro');not null" json:"causa" binding:"required"`
	Notas     *string        `gorm:"type:text" json:"notas,omitempty"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	// Relaciones
	Granja Granja  `gorm:"foreignKey:GranjaID;constraint:OnDelete:RESTRICT" json:"granja,omitempty"`
	Parto  *Parto  `gorm:"foreignKey:PartoID;constraint:OnDelete:SET NULL" json:"parto,omitempty"`
	Corral *Corral `gorm:"foreignKey:CorralID;constraint:OnDelete:SET NULL" json:"corral,omitempty"`
}

// TableName especifica el nombre de la tabla
func (MuerteLechon) TableName() string {
	return "muertes_lechones"
}

// BeforeCreate hook de GORM para validaciones antes de crear
func (m *MuerteLechon) BeforeCreate(tx *gorm.DB) error {
	return m.validate()
}

// BeforeUpdate hook de GORM para validaciones antes de actualizar
func (m *MuerteLechon) BeforeUpdate(tx *gorm.DB) error {
	return m.validate()
}

// validate valida que se asocie a un parto o a un corral (no ambos, no ninguno)
func (m *MuerteLechon) validate() error {
	tienePartoID := m.PartoID != nil && *m.PartoID > 0
	tieneCorralID := m.CorralID != nil && *m.CorralID > 0

	if !tienePartoID && !tieneCorralID {
		return gorm.ErrInvalidValue
	}
	if tienePartoID && tieneCorralID {
		return gorm.ErrInvalidValue
	}
	if m.Cantidad < 1 {
		return gorm.ErrInvalidValue
	}
	return nil
}
