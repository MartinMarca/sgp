package authz

import "github.com/martin/sgp/internal/models"

// CanViewUser indica si el actor puede ver el usuario objetivo.
func (s *Service) CanViewUser(actor Actor, target *models.Usuario) error {
	if actor.ID == target.ID {
		return nil
	}
	if actor.IsAdmin() {
		return nil
	}
	if actor.Rol == models.RolPropietario {
		if err := s.Require(actor, PermUsersManage); err != nil {
			return err
		}
		if target.Rol == models.RolEmpleado && target.PropietarioID != nil && *target.PropietarioID == actor.ID {
			return nil
		}
		return ErrForbidden
	}
	return ErrForbidden
}

// CanModifyUser indica si el actor puede editar el usuario objetivo.
func (s *Service) CanModifyUser(actor Actor, target *models.Usuario) error {
	return s.CanViewUser(actor, target)
}

// CanDeleteUser indica si el actor puede desactivar el usuario objetivo.
func (s *Service) CanDeleteUser(actor Actor, target *models.Usuario) error {
	if actor.ID == target.ID {
		return ErrForbidden
	}
	if actor.IsAdmin() {
		return nil
	}
	if actor.Rol == models.RolPropietario {
		if err := s.Require(actor, PermUsersManage); err != nil {
			return err
		}
		if target.Rol == models.RolEmpleado && target.PropietarioID != nil && *target.PropietarioID == actor.ID {
			return nil
		}
		return ErrForbidden
	}
	return ErrForbidden
}

// ValidateCreateRole verifica si el actor puede crear un usuario con el rol indicado.
func (s *Service) ValidateCreateRole(actor Actor, rol string) error {
	if err := s.Require(actor, PermUsersManage); err != nil {
		return err
	}
	if actor.IsAdmin() {
		return nil
	}
	if actor.Rol == models.RolPropietario {
		if rol == models.RolPropietario || rol == models.RolEmpleado {
			return nil
		}
		return ErrForbidden
	}
	return ErrForbidden
}
