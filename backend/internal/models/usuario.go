package models

import (
	"time"

	"gorm.io/gorm"
)

// Usuario representa un usuario del sistema
type Usuario struct {
	ID              uint           `gorm:"primaryKey" json:"id"`
	Username        string         `gorm:"uniqueIndex;size:50;not null" json:"username" binding:"required,min=3,max=50"`
	Email           string         `gorm:"uniqueIndex;size:100;not null" json:"email" binding:"required,email"`
	PasswordHash    string         `gorm:"size:255;not null" json:"-"`
	NombreCompleto  *string        `gorm:"size:100" json:"nombre_completo"`
	Establecimiento *string        `gorm:"size:150" json:"establecimiento"`
	Rol             string         `gorm:"type:enum('admin','propietario','empleado');default:'propietario'" json:"rol"`
	PropietarioID   *uint          `gorm:"index" json:"propietario_id,omitempty"`
	Activo          bool           `gorm:"default:true" json:"activo"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"-"`

	Propietario *Usuario `gorm:"foreignKey:PropietarioID" json:"propietario,omitempty"`
	Empleados   []Usuario `gorm:"foreignKey:PropietarioID" json:"empleados,omitempty"`
	Granjas     []Granja  `gorm:"foreignKey:PropietarioID" json:"granjas,omitempty"`
}

// TableName especifica el nombre de la tabla
func (Usuario) TableName() string {
	return "usuarios"
}
