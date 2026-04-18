package services

import (
	"time"

	"github.com/martin/sgp/internal/models"
	"github.com/martin/sgp/internal/repositories"
	"gorm.io/gorm"
)

// MuerteLechonService maneja la lógica de negocio de muertes de lechones
type MuerteLechonService struct {
	db    *gorm.DB
	repos *repositories.RepositoryContainer
}

// NewMuerteLechonService crea una nueva instancia del servicio
func NewMuerteLechonService(db *gorm.DB, repos *repositories.RepositoryContainer) *MuerteLechonService {
	return &MuerteLechonService{db: db, repos: repos}
}

// --- DTOs ---

// CrearMuerteLechonInput datos para registrar una muerte de animales
type CrearMuerteLechonInput struct {
	GranjaID uint   `json:"granja_id" binding:"required"`
	PartoID  *uint  `json:"parto_id"`
	CorralID *uint  `json:"corral_id"`
	Fecha    string `json:"fecha" binding:"required"`
	Cantidad int    `json:"cantidad" binding:"required,gte=1"`
	Causa    string `json:"causa" binding:"required"`
	Notas    string `json:"notas"`
}

// ActualizarMuerteLechonInput datos para actualizar un registro de muerte
type ActualizarMuerteLechonInput struct {
	Fecha    *string `json:"fecha"`
	Cantidad *int    `json:"cantidad"`
	Causa    *string `json:"causa"`
	Notas    *string `json:"notas"`
}

// --- Métodos del servicio ---

var causasValidas = map[string]bool{
	models.CausaMuerteAplastamiento: true,
	models.CausaMuerteEnfermedad:    true,
	models.CausaMuerteInanicion:     true,
	models.CausaMuerteOtro:          true,
}

// Crear registra una nueva muerte de lechones
func (s *MuerteLechonService) Crear(input CrearMuerteLechonInput) (*models.MuerteLechon, error) {
	if !causasValidas[input.Causa] {
		return nil, ErrCausaMuerteInvalida
	}

	tienePartoID := input.PartoID != nil && *input.PartoID > 0
	tieneCorralID := input.CorralID != nil && *input.CorralID > 0

	if (!tienePartoID && !tieneCorralID) || (tienePartoID && tieneCorralID) {
		return nil, ErrMuerteRequierePartoOLote
	}

	fecha, err := time.Parse("2006-01-02", input.Fecha)
	if err != nil {
		return nil, err
	}

	var notas *string
	if input.Notas != "" {
		notas = &input.Notas
	}

	var muerte *models.MuerteLechon

	err = s.db.Transaction(func(tx *gorm.DB) error {
		if tienePartoID {
			return s.crearMuerteLactancia(tx, input, fecha, notas, &muerte)
		}
		return s.crearMuerteEngorde(tx, input, fecha, notas, &muerte)
	})
	if err != nil {
		return nil, err
	}

	return muerte, nil
}

func (s *MuerteLechonService) crearMuerteLactancia(tx *gorm.DB, input CrearMuerteLechonInput, fecha time.Time, notas *string, out **models.MuerteLechon) error {
	var parto models.Parto
	if err := tx.First(&parto, *input.PartoID).Error; err != nil {
		return ErrNotFound
	}

	// Calcular muertes ya registradas para este parto
	var muertesExistentes int
	tx.Model(&models.MuerteLechon{}).
		Select("COALESCE(SUM(cantidad), 0)").
		Where("parto_id = ?", *input.PartoID).
		Scan(&muertesExistentes)

	if muertesExistentes+input.Cantidad > parto.LechonesNacidosVivos {
		return ErrMuertesExcedenVivos
	}

	muerte := &models.MuerteLechon{
		GranjaID: input.GranjaID,
		PartoID:  input.PartoID,
		Fecha:    fecha,
		Cantidad: input.Cantidad,
		Causa:    input.Causa,
		Notas:    notas,
	}
	if err := tx.Create(muerte).Error; err != nil {
		return err
	}

	*out = muerte
	return nil
}

func (s *MuerteLechonService) crearMuerteEngorde(tx *gorm.DB, input CrearMuerteLechonInput, fecha time.Time, notas *string, out **models.MuerteLechon) error {
	var corral models.Corral
	if err := tx.First(&corral, *input.CorralID).Error; err != nil {
		return ErrNotFound
	}

	// Validar que la cantidad no exceda los animales activos en el corral
	var totalAnimales int
	tx.Model(&models.Lote{}).
		Select("COALESCE(SUM(cantidad_lechones), 0)").
		Where("corral_id = ? AND estado = ?", *input.CorralID, models.EstadoLoteActivo).
		Scan(&totalAnimales)

	if input.Cantidad > totalAnimales {
		return ErrMuertesExcedenLote
	}

	muerte := &models.MuerteLechon{
		GranjaID: input.GranjaID,
		CorralID: input.CorralID,
		Fecha:    fecha,
		Cantidad: input.Cantidad,
		Causa:    input.Causa,
		Notas:    notas,
	}
	if err := tx.Create(muerte).Error; err != nil {
		return err
	}

	*out = muerte
	return nil
}

// ObtenerPorID obtiene un registro de muerte con sus relaciones
func (s *MuerteLechonService) ObtenerPorID(id uint) (*models.MuerteLechon, error) {
	muerte, err := s.repos.MuerteLechon.FindByID(id, "Granja", "Parto", "Parto.Cerda", "Corral")
	if err != nil {
		return nil, ErrNotFound
	}
	return muerte, nil
}

// ListarPorGranja lista las muertes de una granja
func (s *MuerteLechonService) ListarPorGranja(granjaID uint) ([]models.MuerteLechon, error) {
	return s.repos.MuerteLechon.FindByGranjaID(granjaID)
}

// ListarPorParto lista las muertes asociadas a un parto (lactancia)
func (s *MuerteLechonService) ListarPorParto(partoID uint) ([]models.MuerteLechon, error) {
	return s.repos.MuerteLechon.FindByPartoID(partoID)
}

// ListarPorCorral lista las muertes asociadas a un corral (engorde)
func (s *MuerteLechonService) ListarPorCorral(corralID uint) ([]models.MuerteLechon, error) {
	return s.repos.MuerteLechon.FindByCorralID(corralID)
}

// ListarPorPeriodo lista muertes por mes/año
func (s *MuerteLechonService) ListarPorPeriodo(granjaID *uint, mes, anio int) ([]models.MuerteLechon, error) {
	return s.repos.MuerteLechon.FindByPeriodo(granjaID, mes, anio)
}

// Actualizar modifica un registro de muerte
func (s *MuerteLechonService) Actualizar(id uint, input ActualizarMuerteLechonInput) (*models.MuerteLechon, error) {
	muerte, err := s.repos.MuerteLechon.FindByID(id)
	if err != nil {
		return nil, ErrNotFound
	}

	if input.Fecha != nil {
		fecha, err := time.Parse("2006-01-02", *input.Fecha)
		if err != nil {
			return nil, err
		}
		muerte.Fecha = fecha
	}

	if input.Causa != nil {
		if !causasValidas[*input.Causa] {
			return nil, ErrCausaMuerteInvalida
		}
		muerte.Causa = *input.Causa
	}

	if input.Notas != nil {
		muerte.Notas = input.Notas
	}

	if input.Cantidad != nil {
		if muerte.PartoID != nil {
			// Validar que no exceda los vivos del parto
			var parto models.Parto
			if err := s.db.First(&parto, *muerte.PartoID).Error; err != nil {
				return nil, ErrNotFound
			}
			var otherMuertes int
			s.db.Model(&models.MuerteLechon{}).
				Select("COALESCE(SUM(cantidad), 0)").
				Where("parto_id = ? AND id != ?", *muerte.PartoID, muerte.ID).
				Scan(&otherMuertes)
			if otherMuertes+*input.Cantidad > parto.LechonesNacidosVivos {
				return nil, ErrMuertesExcedenVivos
			}
		}
		muerte.Cantidad = *input.Cantidad
	}

	if err := s.db.Save(muerte).Error; err != nil {
		return nil, err
	}

	return muerte, nil
}

// Eliminar borra un registro de muerte
func (s *MuerteLechonService) Eliminar(id uint) error {
	muerte, err := s.repos.MuerteLechon.FindByID(id)
	if err != nil {
		return ErrNotFound
	}

	return s.db.Delete(&models.MuerteLechon{}, muerte.ID).Error
}

// GetEstadisticas obtiene estadísticas de mortalidad
func (s *MuerteLechonService) GetEstadisticas(granjaID *uint, mes, anio int) (map[string]interface{}, error) {
	return s.repos.MuerteLechon.GetEstadisticas(granjaID, mes, anio)
}
