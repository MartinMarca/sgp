package authz_test

import (
	"testing"

	"github.com/martin/sgp/internal/authz"
	"github.com/martin/sgp/internal/models"
)

func TestHasPermission_AdminTieneTodo(t *testing.T) {
	perms := authz.PermisosEfectivos(models.RolAdmin)
	if len(perms) != 14 {
		t.Fatalf("admin debería tener 14 permisos, obtuvo %d", len(perms))
	}
}

func TestHasPermission_EmpleadoSinVentas(t *testing.T) {
	if authz.HasPermission(models.RolEmpleado, authz.PermVentaRead) {
		t.Fatal("empleado no debería tener venta:read")
	}
	if authz.HasPermission(models.RolEmpleado, authz.PermCerdaDelete) {
		t.Fatal("empleado no debería tener cerda:delete")
	}
	if !authz.HasPermission(models.RolEmpleado, authz.PermGranjaRead) {
		t.Fatal("empleado debería tener granja:read")
	}
	if !authz.HasPermission(models.RolEmpleado, authz.PermRecursoWrite) {
		t.Fatal("empleado debería tener recurso:write")
	}
}

func TestHasPermission_PropietarioConVentas(t *testing.T) {
	if !authz.HasPermission(models.RolPropietario, authz.PermVentaWrite) {
		t.Fatal("propietario debería tener venta:write")
	}
	if authz.HasPermission(models.RolPropietario, authz.PermUsersManageAll) {
		t.Fatal("propietario no debería tener users:manage_all")
	}
}

func TestActor_OwnerID(t *testing.T) {
	propID := uint(5)
	emp := authz.Actor{Rol: models.RolEmpleado, PropietarioID: &propID}
	if emp.OwnerID() != 5 {
		t.Fatalf("OwnerID empleado: esperaba 5, obtuvo %d", emp.OwnerID())
	}
	prop := authz.Actor{ID: 3, Rol: models.RolPropietario}
	if prop.OwnerID() != 3 {
		t.Fatalf("OwnerID propietario: esperaba 3, obtuvo %d", prop.OwnerID())
	}
}
