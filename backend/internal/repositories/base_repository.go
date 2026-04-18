package repositories

import (
	"errors"

	"gorm.io/gorm"
)

// BaseRepository proporciona operaciones CRUD genéricas
type BaseRepository struct {
	db *gorm.DB
}

// NewBaseRepository crea una nueva instancia de BaseRepository
func NewBaseRepository(db *gorm.DB) *BaseRepository {
	return &BaseRepository{db: db}
}

// Create crea un nuevo registro
func (r *BaseRepository) Create(model interface{}) error {
	return r.db.Create(model).Error
}

// FindByID busca un registro por ID
func (r *BaseRepository) FindByID(model interface{}, id uint) error {
	return r.db.First(model, id).Error
}

// Update actualiza un registro
func (r *BaseRepository) Update(model interface{}) error {
	return r.db.Save(model).Error
}

// Delete elimina un registro (soft delete)
func (r *BaseRepository) Delete(model interface{}, id uint) error {
	return r.db.Delete(model, id).Error
}

// FindAll obtiene todos los registros
func (r *BaseRepository) FindAll(models interface{}, conditions ...interface{}) error {
	query := r.db
	if len(conditions) > 0 {
		query = query.Where(conditions[0], conditions[1:]...)
	}
	return query.Find(models).Error
}

// Count cuenta los registros que cumplen una condición
func (r *BaseRepository) Count(model interface{}, conditions ...interface{}) (int64, error) {
	var count int64
	query := r.db.Model(model)
	if len(conditions) > 0 {
		query = query.Where(conditions[0], conditions[1:]...)
	}
	err := query.Count(&count).Error
	return count, err
}

// Exists verifica si existe un registro
func (r *BaseRepository) Exists(model interface{}, conditions ...interface{}) (bool, error) {
	count, err := r.Count(model, conditions...)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// Transaction ejecuta una función dentro de una transacción
func (r *BaseRepository) Transaction(fn func(*gorm.DB) error) error {
	return r.db.Transaction(fn)
}

// Common errors
var (
	ErrNotFound      = gorm.ErrRecordNotFound
	ErrDuplicateKey  = errors.New("duplicate key")
	ErrInvalidValue  = errors.New("invalid value")
)
