package repositories

import (
	"github.com/martin/sgp/internal/models"
	"gorm.io/gorm"
)

// UsuarioRepository maneja las operaciones de base de datos para Usuarios
type UsuarioRepository struct {
	*BaseRepository
}

// NewUsuarioRepository crea una nueva instancia de UsuarioRepository
func NewUsuarioRepository(db *gorm.DB) *UsuarioRepository {
	return &UsuarioRepository{
		BaseRepository: NewBaseRepository(db),
	}
}

// Create crea un nuevo usuario
func (r *UsuarioRepository) Create(usuario *models.Usuario) error {
	return r.db.Create(usuario).Error
}

// FindByID busca un usuario por ID
func (r *UsuarioRepository) FindByID(id uint, preload ...string) (*models.Usuario, error) {
	var usuario models.Usuario
	query := r.db
	
	for _, rel := range preload {
		query = query.Preload(rel)
	}
	
	err := query.First(&usuario, id).Error
	if err != nil {
		return nil, err
	}
	return &usuario, nil
}

// FindByUsername busca un usuario por username
func (r *UsuarioRepository) FindByUsername(username string) (*models.Usuario, error) {
	var usuario models.Usuario
	err := r.db.Where("username = ?", username).First(&usuario).Error
	if err != nil {
		return nil, err
	}
	return &usuario, nil
}

// FindByEmail busca un usuario por email
func (r *UsuarioRepository) FindByEmail(email string) (*models.Usuario, error) {
	var usuario models.Usuario
	err := r.db.Where("email = ?", email).First(&usuario).Error
	if err != nil {
		return nil, err
	}
	return &usuario, nil
}

// FindAll obtiene todos los usuarios
func (r *UsuarioRepository) FindAll(activo *bool) ([]models.Usuario, error) {
	var usuarios []models.Usuario
	query := r.db
	
	if activo != nil {
		query = query.Where("activo = ?", *activo)
	}
	
	err := query.Order("username").Find(&usuarios).Error
	return usuarios, err
}

// Update actualiza un usuario
func (r *UsuarioRepository) Update(usuario *models.Usuario) error {
	return r.db.Save(usuario).Error
}

// Delete elimina un usuario (soft delete)
func (r *UsuarioRepository) Delete(id uint) error {
	return r.db.Delete(&models.Usuario{}, id).Error
}

// ExisteUsername verifica si ya existe un username
func (r *UsuarioRepository) ExisteUsername(username string, excludeID *uint) (bool, error) {
	query := r.db.Model(&models.Usuario{}).Where("username = ?", username)
	
	if excludeID != nil {
		query = query.Where("id != ?", *excludeID)
	}
	
	var count int64
	err := query.Count(&count).Error
	return count > 0, err
}

// ExisteEmail verifica si ya existe un email
func (r *UsuarioRepository) ExisteEmail(email string, excludeID *uint) (bool, error) {
	query := r.db.Model(&models.Usuario{}).Where("email = ?", email)
	
	if excludeID != nil {
		query = query.Where("id != ?", *excludeID)
	}
	
	var count int64
	err := query.Count(&count).Error
	return count > 0, err
}

// GetGranjas obtiene todas las granjas de un usuario
func (r *UsuarioRepository) GetGranjas(usuarioID uint) ([]models.Granja, error) {
	var granjas []models.Granja
	err := r.db.
		Joins("JOIN usuario_granja ON usuario_granja.granja_id = granjas.id").
		Where("usuario_granja.usuario_id = ?", usuarioID).
		Find(&granjas).Error
	return granjas, err
}

// GetRolEnGranja obtiene el rol de un usuario en una granja
func (r *UsuarioRepository) GetRolEnGranja(usuarioID, granjaID uint) (string, error) {
	var usuarioGranja models.UsuarioGranja
	err := r.db.
		Where("usuario_id = ? AND granja_id = ?", usuarioID, granjaID).
		First(&usuarioGranja).Error
	
	if err != nil {
		return "", err
	}
	
	return usuarioGranja.Rol, nil
}
