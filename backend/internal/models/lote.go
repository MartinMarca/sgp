package models

import (
	"database/sql"
	"time"

	"gorm.io/gorm"
)

// Lote representa un lote de lechones destetados
type Lote struct {
	ID               uint           `gorm:"primaryKey" json:"id"`
	CorralID         uint           `gorm:"not null;index" json:"corral_id" binding:"required"` // Obligatorio
	Nombre           string         `gorm:"size:100;not null" json:"nombre" binding:"required,min=1,max=100"`
	CantidadLechones int            `gorm:"not null;default:0" json:"cantidad_lechones" binding:"gte=0"`
	FechaCreacion    time.Time      `gorm:"type:date;not null" json:"fecha_creacion" binding:"required"`
	Estado           string         `gorm:"type:enum('activo','cerrado','vendido');default:'activo'" json:"estado"`
	FechaCierre      *sql.NullTime  `gorm:"type:date" json:"fecha_cierre,omitempty"`
	MotivoCierre     *string        `gorm:"type:text" json:"motivo_cierre,omitempty"`
	CreatedAt        time.Time      `json:"created_at"`
	UpdatedAt        time.Time      `json:"updated_at"`
	DeletedAt        gorm.DeletedAt `gorm:"index" json:"-"`

	// Relaciones
	Corral   Corral    `gorm:"foreignKey:CorralID;constraint:OnDelete:RESTRICT" json:"corral,omitempty"`
	Destetes []Destete `gorm:"foreignKey:LoteID" json:"destetes,omitempty"`
}

// TableName especifica el nombre de la tabla
func (Lote) TableName() string {
	return "lotes"
}
