package models

import (
	"database/sql"
	"time"

	"gorm.io/gorm"
)

// Servicio representa un servicio de monta (natural o inseminación)
type Servicio struct {
	ID                       uint           `gorm:"primaryKey" json:"id"`
	CerdaID                  uint           `gorm:"not null;index:idx_servicios_cerda_fecha" json:"cerda_id" binding:"required"`
	FechaServicio            time.Time      `gorm:"type:date;not null;index:idx_servicios_cerda_fecha" json:"fecha_servicio" binding:"required"`
	TieneRepeticiones        bool           `gorm:"default:false" json:"tiene_repeticiones"`
	CantidadRepeticiones     int            `gorm:"default:0" json:"cantidad_repeticiones" binding:"gte=0"`
	TipoMonta                string         `gorm:"type:enum('natural','inseminacion');not null" json:"tipo_monta" binding:"required,oneof=natural inseminacion"`
	
	// Campos para Monta Natural
	PadrilloID               *uint          `gorm:"index" json:"padrillo_id,omitempty"`
	CantidadSaltos           *int           `json:"cantidad_saltos,omitempty" binding:"omitempty,gte=0"`
	
	// Campos para Inseminación
	NumeroPajuela            *string        `gorm:"size:50" json:"numero_pajuela,omitempty"`
	
	// Control de Preñez
	PrenezConfirmada         bool           `gorm:"default:false;index" json:"prenez_confirmada"`
	FechaConfirmacionPrenez  *sql.NullTime  `gorm:"type:date" json:"fecha_confirmacion_prenez,omitempty"`
	FechaEstimadaParto       *time.Time     `gorm:"type:date" json:"fecha_estimada_parto,omitempty"`
	PrenezCancelada          bool           `gorm:"default:false" json:"prenez_cancelada"`
	FechaCancelacionPrenez   *sql.NullTime  `gorm:"type:date" json:"fecha_cancelacion_prenez,omitempty"`
	MotivoCancelacion        *string        `gorm:"type:text" json:"motivo_cancelacion,omitempty"`
	
	CreatedAt                time.Time      `json:"created_at"`
	UpdatedAt                time.Time      `json:"updated_at"`
	DeletedAt                gorm.DeletedAt `gorm:"index" json:"-"`

	// Relaciones
	Cerda    Cerda     `gorm:"foreignKey:CerdaID;constraint:OnDelete:RESTRICT" json:"cerda,omitempty"`
	Padrillo *Padrillo `gorm:"foreignKey:PadrilloID;constraint:OnDelete:SET NULL" json:"padrillo,omitempty"`
	Partos   []Parto   `gorm:"foreignKey:ServicioID" json:"partos,omitempty"`
}

// TableName especifica el nombre de la tabla
func (Servicio) TableName() string {
	return "servicios"
}

// BeforeCreate hook de GORM para validaciones antes de crear
func (s *Servicio) BeforeCreate(tx *gorm.DB) error {
	return s.validate()
}

// BeforeUpdate hook de GORM para validaciones antes de actualizar
func (s *Servicio) BeforeUpdate(tx *gorm.DB) error {
	return s.validate()
}

// validate valida que los campos sean consistentes según el tipo de monta
func (s *Servicio) validate() error {
	// Validar que si es monta natural, tenga padrillo
	if s.TipoMonta == "natural" && s.PadrilloID == nil {
		return gorm.ErrInvalidValue
	}
	// Validar que si es inseminación, tenga número de pajuela
	if s.TipoMonta == "inseminacion" && (s.NumeroPajuela == nil || *s.NumeroPajuela == "") {
		return gorm.ErrInvalidValue
	}
	return nil
}
