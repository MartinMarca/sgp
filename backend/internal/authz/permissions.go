package authz

import "github.com/martin/sgp/internal/models"

// Permisos atómicos del sistema.
const (
	PermUsersManage     = "users:manage"
	PermUsersManageAll  = "users:manage_all"
	PermGranjaRead      = "granja:read"
	PermGranjaWrite     = "granja:write"
	PermGranjaDelete    = "granja:delete"
	PermCerdaWrite      = "cerda:write"
	PermCerdaDelete     = "cerda:delete"
	PermPadrilloWrite   = "padrillo:write"
	PermPadrilloDelete  = "padrillo:delete"
	PermVentaRead       = "venta:read"
	PermVentaWrite      = "venta:write"
	PermStatsRead       = "stats:read"
	PermStatsVentasRead = "stats:ventas:read"
	PermRecursoWrite    = "recurso:write"
)

var rolePermissions = map[string]map[string]bool{
	models.RolAdmin: {
		PermUsersManage:     true,
		PermUsersManageAll:  true,
		PermGranjaRead:      true,
		PermGranjaWrite:     true,
		PermGranjaDelete:    true,
		PermCerdaWrite:      true,
		PermCerdaDelete:     true,
		PermPadrilloWrite:   true,
		PermPadrilloDelete:  true,
		PermVentaRead:       true,
		PermVentaWrite:      true,
		PermStatsRead:       true,
		PermStatsVentasRead: true,
		PermRecursoWrite:    true,
	},
	models.RolPropietario: {
		PermUsersManage:     true,
		PermGranjaRead:      true,
		PermGranjaWrite:     true,
		PermGranjaDelete:    true,
		PermCerdaWrite:      true,
		PermCerdaDelete:     true,
		PermPadrilloWrite:   true,
		PermPadrilloDelete:  true,
		PermVentaRead:       true,
		PermVentaWrite:      true,
		PermStatsRead:       true,
		PermStatsVentasRead: true,
		PermRecursoWrite:    true,
	},
	models.RolEmpleado: {
		PermGranjaRead:     true,
		PermCerdaWrite:     true,
		PermPadrilloWrite:  true,
		PermStatsRead:      true,
		PermRecursoWrite:   true,
	},
}

// HasPermission indica si un rol tiene un permiso global (sin scope de granja).
func HasPermission(rol, perm string) bool {
	if rol == models.RolAdmin {
		return true
	}
	perms, ok := rolePermissions[rol]
	if !ok {
		return false
	}
	return perms[perm]
}

// PermisosEfectivos lista los permisos del rol (admin recibe todos).
func PermisosEfectivos(rol string) []string {
	all := []string{
		PermUsersManage,
		PermUsersManageAll,
		PermGranjaRead,
		PermGranjaWrite,
		PermGranjaDelete,
		PermCerdaWrite,
		PermCerdaDelete,
		PermPadrilloWrite,
		PermPadrilloDelete,
		PermVentaRead,
		PermVentaWrite,
		PermStatsRead,
		PermStatsVentasRead,
		PermRecursoWrite,
	}
	if rol == models.RolAdmin {
		return all
	}
	perms, ok := rolePermissions[rol]
	if !ok {
		return nil
	}
	var result []string
	for _, p := range all {
		if perms[p] {
			result = append(result, p)
		}
	}
	return result
}
