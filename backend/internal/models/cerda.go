package models

import (
	"database/sql"
	"time"

	"gorm.io/gorm"
)

// Cerda representa una cerda (cerda hembra reproductora)
type Cerda struct {
	ID             uint           `gorm:"primaryKey" json:"id"`
	GranjaID       uint           `gorm:"not null;uniqueIndex:idx_caravana_granja_cerda" json:"granja_id" binding:"required"`
	NumeroCaravana string         `gorm:"size:50;not null;uniqueIndex:idx_caravana_granja_cerda" json:"numero_caravana" binding:"required,min=1,max=50"`
	DetallePelaje  *string        `gorm:"type:text" json:"detalle_pelaje"`
	Genetica       *string        `gorm:"size:100" json:"genetica"`
	Estado         string         `gorm:"type:enum('disponible','servicio','gestacion','cria');default:'disponible'" json:"estado"`
	Activo         bool           `gorm:"default:true" json:"activo"`
	FechaBaja      *sql.NullTime  `gorm:"type:date" json:"fecha_baja,omitempty"`
	MotivoBaja     *string        `gorm:"type:enum('muerte','venta')" json:"motivo_baja,omitempty"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`

	// Relaciones
	Granja    Granja     `gorm:"foreignKey:GranjaID;constraint:OnDelete:RESTRICT" json:"granja,omitempty"`
	Servicios []Servicio `gorm:"foreignKey:CerdaID" json:"servicios,omitempty"`
	Partos    []Parto    `gorm:"foreignKey:CerdaID" json:"partos,omitempty"`
	Destetes  []Destete  `gorm:"foreignKey:CerdaID" json:"destetes,omitempty"`
}

// TableName especifica el nombre de la tabla
func (Cerda) TableName() string {
	return "cerdas"
}
