package models

import (
	"time"

	"gorm.io/gorm"
)

// Corral representa un corral dentro de una granja
type Corral struct {
	ID              uint           `gorm:"primaryKey" json:"id"`
	GranjaID        uint           `gorm:"not null;index" json:"granja_id" binding:"required"`
	Nombre          string         `gorm:"size:100;not null" json:"nombre" binding:"required,min=1,max=100"`
	Descripcion     *string        `gorm:"type:text" json:"descripcion"`
	CapacidadMaxima *int           `json:"capacidad_maxima"` // Opcional
	Activo          bool           `gorm:"default:true" json:"activo"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"-"`

	TotalAnimales int `gorm:"-" json:"total_animales"`

	// Relaciones
	Granja Granja `gorm:"foreignKey:GranjaID;constraint:OnDelete:RESTRICT" json:"granja,omitempty"`
	Lotes  []Lote `gorm:"foreignKey:CorralID" json:"lotes,omitempty"`
}

// TableName especifica el nombre de la tabla
func (Corral) TableName() string {
	return "corrales"
}
