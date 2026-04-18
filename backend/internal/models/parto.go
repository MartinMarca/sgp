package models

import (
	"time"

	"gorm.io/gorm"
)

// Parto representa un parto de una cerda
type Parto struct {
	ID                    uint           `gorm:"primaryKey" json:"id"`
	CerdaID               uint           `gorm:"not null;index:idx_partos_cerda_fecha" json:"cerda_id" binding:"required"`
	ServicioID            *uint          `gorm:"index" json:"servicio_id,omitempty"`
	FechaParto            time.Time      `gorm:"type:date;not null;index:idx_partos_cerda_fecha" json:"fecha_parto" binding:"required"`
	LechonesNacidosVivos  int            `gorm:"not null;default:0" json:"lechones_nacidos_vivos" binding:"required,gte=0"`
	LechonesNacidosTotales int           `gorm:"not null;default:0" json:"lechones_nacidos_totales" binding:"required,gte=0"`
	LechonesHembras       int            `gorm:"not null;default:0" json:"lechones_hembras" binding:"required,gte=0"`
	LechonesMachos        int            `gorm:"not null;default:0" json:"lechones_machos" binding:"required,gte=0"`
	FechaEstimada         time.Time      `gorm:"type:date;not null;index" json:"fecha_estimada" binding:"required"`
	CreatedAt             time.Time      `json:"created_at"`
	UpdatedAt             time.Time      `json:"updated_at"`
	DeletedAt             gorm.DeletedAt `gorm:"index" json:"-"`

	// Relaciones
	Cerda    Cerda     `gorm:"foreignKey:CerdaID;constraint:OnDelete:RESTRICT" json:"cerda,omitempty"`
	Servicio *Servicio `gorm:"foreignKey:ServicioID;constraint:OnDelete:SET NULL" json:"servicio,omitempty"`
	Destetes       []Destete       `gorm:"foreignKey:PartoID" json:"destetes,omitempty"`
	MuerteLechones []MuerteLechon  `gorm:"foreignKey:PartoID" json:"muertes_lechones,omitempty"`
}

// TableName especifica el nombre de la tabla
func (Parto) TableName() string {
	return "partos"
}

// BeforeCreate hook de GORM para validaciones antes de crear
func (p *Parto) BeforeCreate(tx *gorm.DB) error {
	return p.validate()
}

// BeforeUpdate hook de GORM para validaciones antes de actualizar
func (p *Parto) BeforeUpdate(tx *gorm.DB) error {
	return p.validate()
}

// validate valida las reglas de negocio del parto
func (p *Parto) validate() error {
	// Validar: lechones_hembras + lechones_machos == lechones_nacidos_vivos
	if p.LechonesHembras+p.LechonesMachos != p.LechonesNacidosVivos {
		return gorm.ErrInvalidValue
	}
	// Validar: lechones_nacidos_totales >= lechones_nacidos_vivos
	if p.LechonesNacidosTotales < p.LechonesNacidosVivos {
		return gorm.ErrInvalidValue
	}
	return nil
}
