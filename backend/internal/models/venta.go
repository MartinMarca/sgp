package models

import (
	"time"

	"gorm.io/gorm"
)

// Tipos de animal vendido
const (
	TipoAnimalVentaCerda    = "cerda"
	TipoAnimalVentaPadrillo = "padrillo"
	TipoAnimalVentaLechon   = "lechon"
)

// Venta representa la venta de animales
type Venta struct {
	ID         uint           `gorm:"primaryKey" json:"id"`
	GranjaID   uint           `gorm:"not null;index" json:"granja_id"`
	Fecha      time.Time      `gorm:"type:date;not null;index" json:"fecha"`
	TipoAnimal string         `gorm:"type:enum('cerda','padrillo','lechon');not null" json:"tipo_animal"`
	Cantidad   int            `gorm:"not null" json:"cantidad"`
	KgTotales  float64        `gorm:"type:decimal(10,2);not null" json:"kg_totales"`
	Monto      float64        `gorm:"type:decimal(12,2);not null" json:"monto"`
	Comprador  string         `gorm:"type:varchar(200);not null" json:"comprador"`
	LoteID     *uint          `gorm:"index" json:"lote_id,omitempty"`
	CorralID   *uint          `gorm:"index" json:"corral_id,omitempty"`
	Notas      *string        `gorm:"type:text" json:"notas,omitempty"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`

	Granja Granja  `gorm:"foreignKey:GranjaID;constraint:OnDelete:RESTRICT" json:"granja,omitempty"`
	Lote   *Lote   `gorm:"foreignKey:LoteID;constraint:OnDelete:SET NULL" json:"lote,omitempty"`
	Corral *Corral `gorm:"foreignKey:CorralID;constraint:OnDelete:SET NULL" json:"corral,omitempty"`
}

func (Venta) TableName() string {
	return "ventas"
}
