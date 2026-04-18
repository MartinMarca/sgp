package models

import (
	"time"

	"gorm.io/gorm"
)

// Granja representa una granja porcina
type Granja struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	Nombre      string         `gorm:"size:100;not null" json:"nombre" binding:"required,min=3,max=100"`
	Descripcion *string        `gorm:"type:text" json:"descripcion"`
	Ubicacion   *string        `gorm:"size:200" json:"ubicacion"`
	Activo      bool           `gorm:"default:true" json:"activo"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`

	// Relaciones
	Usuarios  []Usuario  `gorm:"many2many:usuario_granja;" json:"usuarios,omitempty"`
	Corrales  []Corral   `gorm:"foreignKey:GranjaID" json:"corrales,omitempty"`
	Cerdas    []Cerda    `gorm:"foreignKey:GranjaID" json:"cerdas,omitempty"`
	Padrillos []Padrillo `gorm:"foreignKey:GranjaID" json:"padrillos,omitempty"`
}

// TableName especifica el nombre de la tabla
func (Granja) TableName() string {
	return "granjas"
}
