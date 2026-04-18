package models

import (
	"database/sql"
	"time"

	"gorm.io/gorm"
)

// Padrillo representa un padrillo (cerdo macho reproductor)
type Padrillo struct {
	ID                     uint           `gorm:"primaryKey" json:"id"`
	GranjaID               uint           `gorm:"not null;uniqueIndex:idx_caravana_granja_padrillo" json:"granja_id" binding:"required"`
	NumeroCaravana         string         `gorm:"size:50;not null;uniqueIndex:idx_caravana_granja_padrillo" json:"numero_caravana" binding:"required,min=1,max=50"`
	Nombre                 string         `gorm:"size:100;not null" json:"nombre" binding:"required,min=1,max=100"`
	Genetica               *string        `gorm:"size:100" json:"genetica"`
	FechaUltimaVacunacion  *sql.NullTime  `gorm:"type:date" json:"fecha_ultima_vacunacion,omitempty"`
	Activo                 bool           `gorm:"default:true" json:"activo"`
	FechaBaja              *sql.NullTime  `gorm:"type:date" json:"fecha_baja,omitempty"`
	MotivoBaja             *string        `gorm:"type:enum('muerte','venta')" json:"motivo_baja,omitempty"`
	CreatedAt              time.Time      `json:"created_at"`
	UpdatedAt              time.Time      `json:"updated_at"`
	DeletedAt              gorm.DeletedAt `gorm:"index" json:"-"`

	// Relaciones
	Granja    Granja     `gorm:"foreignKey:GranjaID;constraint:OnDelete:RESTRICT" json:"granja,omitempty"`
	Servicios []Servicio `gorm:"foreignKey:PadrilloID" json:"servicios,omitempty"`
}

// TableName especifica el nombre de la tabla
func (Padrillo) TableName() string {
	return "padrillos"
}
