package models

import "time"

// UsuarioGranja representa la relación N:M entre usuarios y granjas
type UsuarioGranja struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UsuarioID uint      `gorm:"not null;uniqueIndex:idx_usuario_granja" json:"usuario_id"`
	GranjaID  uint      `gorm:"not null;uniqueIndex:idx_usuario_granja" json:"granja_id"`
	Rol       string    `gorm:"type:enum('propietario','administrador','operador');default:'operador'" json:"rol"`
	CreatedAt time.Time `json:"created_at"`

	// Relaciones
	Usuario Usuario `gorm:"foreignKey:UsuarioID;constraint:OnDelete:CASCADE" json:"usuario,omitempty"`
	Granja  Granja  `gorm:"foreignKey:GranjaID;constraint:OnDelete:CASCADE" json:"granja,omitempty"`
}

// TableName especifica el nombre de la tabla
func (UsuarioGranja) TableName() string {
	return "usuario_granja"
}
