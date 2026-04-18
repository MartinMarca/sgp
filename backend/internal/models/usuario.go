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
	PasswordHash    string         `gorm:"size:255;not null" json:"-"` // No se serializa en JSON
	NombreCompleto  *string        `gorm:"size:100" json:"nombre_completo"`
	Establecimiento *string        `gorm:"size:150" json:"establecimiento"`
	Rol             string         `gorm:"type:enum('admin','usuario','veterinario');default:'usuario'" json:"rol"`
	Activo          bool           `gorm:"default:true" json:"activo"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"-"` // Soft delete

	// Relaciones
	Granjas []Granja `gorm:"many2many:usuario_granja;" json:"granjas,omitempty"`
}

// TableName especifica el nombre de la tabla
func (Usuario) TableName() string {
	return "usuarios"
}
