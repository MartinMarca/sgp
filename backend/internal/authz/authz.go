package authz

import (
	"github.com/martin/sgp/internal/models"
	"github.com/martin/sgp/internal/repositories"
)

// Service centraliza la autorización por rol y scope de granja.
type Service struct {
	repos *repositories.RepositoryContainer
}

// NewService crea una instancia del servicio de autorización.
func NewService(repos *repositories.RepositoryContainer) *Service {
	return &Service{repos: repos}
}

// Require verifica un permiso global del rol.
func (s *Service) Require(actor Actor, perm string) error {
	if !HasPermission(actor.Rol, perm) {
		return ErrForbidden
	}
	return nil
}

// RequireGranja verifica permiso y acceso a una granja concreta.
func (s *Service) RequireGranja(actor Actor, granjaID uint, perm string) error {
	if err := s.Require(actor, perm); err != nil {
		return err
	}
	if actor.IsAdmin() {
		return nil
	}
	granja, err := s.repos.Granja.FindByID(granjaID)
	if err != nil {
		return ErrNotFound
	}
	if !s.canAccessGranja(actor, granja) {
		return ErrForbidden
	}
	return nil
}

// GranjasAccesibles devuelve las granjas visibles para el actor.
func (s *Service) GranjasAccesibles(actor Actor, activo *bool) ([]models.Granja, error) {
	if actor.IsAdmin() {
		return s.repos.Granja.FindAll(activo)
	}
	ownerID := actor.OwnerID()
	if ownerID == 0 {
		return []models.Granja{}, nil
	}
	return s.repos.Granja.FindByPropietarioID(ownerID, activo)
}

func (s *Service) canAccessGranja(actor Actor, granja *models.Granja) bool {
	if actor.IsAdmin() {
		return true
	}
	ownerID := actor.OwnerID()
	return ownerID != 0 && granja.PropietarioID == ownerID
}

// ResolvePropietarioIDForCreate determina el propietario_id al crear una granja.
func (s *Service) ResolvePropietarioIDForCreate(actor Actor, requested *uint) (uint, error) {
	if actor.IsAdmin() {
		if requested != nil && *requested > 0 {
			if _, err := s.repos.Usuario.FindByID(*requested); err != nil {
				return 0, ErrNotFound
			}
			return *requested, nil
		}
		return 0, ErrForbidden
	}
	if actor.Rol == models.RolPropietario {
		if err := s.Require(actor, PermGranjaWrite); err != nil {
			return 0, err
		}
		return actor.ID, nil
	}
	return 0, ErrForbidden
}
