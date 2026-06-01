package authz_test

import (
	"testing"

	"github.com/martin/sgp/internal/authz"
	"github.com/martin/sgp/internal/models"
)

func TestCanViewUser_PropietarioVeEmpleado(t *testing.T) {
	svc := authz.NewService(nil)
	propID := uint(1)
	empPropID := uint(1)
	actor := authz.Actor{ID: 1, Rol: models.RolPropietario}
	target := &models.Usuario{ID: 2, Rol: models.RolEmpleado, PropietarioID: &empPropID}

	if err := svc.CanViewUser(actor, target); err != nil {
		t.Fatalf("propietario debería ver empleado propio: %v", err)
	}
	_ = propID
}

func TestCanViewUser_EmpleadoNoVeOtro(t *testing.T) {
	svc := authz.NewService(nil)
	actor := authz.Actor{ID: 2, Rol: models.RolEmpleado, PropietarioID: uintPtr(1)}
	target := &models.Usuario{ID: 3, Rol: models.RolEmpleado, PropietarioID: uintPtr(1)}

	if err := svc.CanViewUser(actor, target); err == nil {
		t.Fatal("empleado no debería ver a otro empleado")
	}
}

func TestValidateCreateRole_PropietarioNoCreaAdmin(t *testing.T) {
	svc := authz.NewService(nil)
	actor := authz.Actor{ID: 1, Rol: models.RolPropietario}

	if err := svc.ValidateCreateRole(actor, models.RolAdmin); err == nil {
		t.Fatal("propietario no debería crear admin")
	}
	if err := svc.ValidateCreateRole(actor, models.RolEmpleado); err != nil {
		t.Fatalf("propietario debería crear empleado: %v", err)
	}
}

func uintPtr(v uint) *uint {
	return &v
}
