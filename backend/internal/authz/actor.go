package authz

import "github.com/martin/sgp/internal/models"

// Actor representa al usuario autenticado para checks de autorización.
type Actor struct {
	ID            uint
	Rol           string
	PropietarioID *uint
}

// ActorFromUsuario construye un Actor desde el modelo de usuario.
func ActorFromUsuario(u *models.Usuario) Actor {
	return Actor{
		ID:            u.ID,
		Rol:           u.Rol,
		PropietarioID: u.PropietarioID,
	}
}

// IsAdmin indica si el actor tiene rol admin.
func (a Actor) IsAdmin() bool {
	return a.Rol == models.RolAdmin
}

// OwnerID devuelve el ID del propietario dueño de las granjas visibles.
func (a Actor) OwnerID() uint {
	switch a.Rol {
	case models.RolAdmin:
		return 0
	case models.RolPropietario:
		return a.ID
	case models.RolEmpleado:
		if a.PropietarioID != nil {
			return *a.PropietarioID
		}
	}
	return 0
}
