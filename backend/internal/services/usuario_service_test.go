package services_test

import (
	"errors"
	"testing"

	"github.com/martin/sgp/internal/authz"
	"github.com/martin/sgp/internal/models"
	"github.com/martin/sgp/internal/services"
	"github.com/martin/sgp/internal/testutil"
)

func TestUsuarioService_PropietarioCreaEmpleado(t *testing.T) {
	env := testutil.SetupTestDB(t)
	propID := testutil.SeedPropietario(t, env)
	actor := testutil.SeedActor(t, env, propID)

	emp, err := env.Services.Usuario.Crear(actor, services.CrearUsuarioInput{
		Username: "empleado_svc",
		Email:    "empleado_svc@test.com",
		Password: "123456",
		Rol:      models.RolEmpleado,
	})
	if err != nil {
		t.Fatalf("Crear empleado: %v", err)
	}
	if emp.PropietarioID == nil || *emp.PropietarioID != propID {
		t.Fatalf("propietario_id esperaba %d, obtuvo %v", propID, emp.PropietarioID)
	}
}

func TestUsuarioService_DuplicateUsername(t *testing.T) {
	env := testutil.SetupTestDB(t)
	propID := testutil.SeedPropietario(t, env)
	actor := testutil.SeedActor(t, env, propID)

	input := services.CrearUsuarioInput{
		Username: "dup_user",
		Email:    "a@test.com",
		Password: "123456",
		Rol:      models.RolEmpleado,
	}
	if _, err := env.Services.Usuario.Crear(actor, input); err != nil {
		t.Fatalf("primer create: %v", err)
	}

	input.Email = "b@test.com"
	_, err := env.Services.Usuario.Crear(actor, input)
	if !errors.Is(err, services.ErrDuplicateKey) {
		t.Fatalf("esperaba ErrDuplicateKey, obtuvo %v", err)
	}
}

func TestUsuarioService_DesactivarUltimoAdmin(t *testing.T) {
	env := testutil.SetupTestDB(t)
	adminID := testutil.SeedAdmin(t, env)

	err := env.Services.Usuario.Desactivar(adminID)
	if !errors.Is(err, services.ErrUltimoAdmin) {
		t.Fatalf("esperaba ErrUltimoAdmin, obtuvo %v", err)
	}
}

func TestUsuarioService_PropietarioListaSoloEmpleados(t *testing.T) {
	env := testutil.SetupTestDB(t)
	propID := testutil.SeedPropietario(t, env)
	actor := testutil.SeedActor(t, env, propID)
	testutil.SeedEmpleado(t, env, propID)
	testutil.SeedPropietario(t, env)

	list, err := env.Services.Usuario.Listar(actor, nil)
	if err != nil {
		t.Fatalf("Listar: %v", err)
	}
	if len(list) != 1 {
		t.Fatalf("propietario debería ver 1 empleado, obtuvo %d", len(list))
	}
	if list[0].Rol != models.RolEmpleado {
		t.Fatalf("esperaba rol empleado, obtuvo %s", list[0].Rol)
	}
}

func TestUsuarioService_EmpleadoNoListaUsuarios(t *testing.T) {
	env := testutil.SetupTestDB(t)
	propID := testutil.SeedPropietario(t, env)
	empID := testutil.SeedEmpleado(t, env, propID)
	actor := testutil.SeedActor(t, env, empID)

	_, err := env.Services.Usuario.Listar(actor, nil)
	if !errors.Is(err, services.ErrForbidden) {
		t.Fatalf("esperaba ErrForbidden, obtuvo %v", err)
	}
}

func TestUsuarioService_AdminCreaEmpleadoRequierePropietario(t *testing.T) {
	env := testutil.SetupTestDB(t)
	adminID := testutil.SeedAdmin(t, env)
	actor := authz.Actor{ID: adminID, Rol: models.RolAdmin}

	_, err := env.Services.Usuario.Crear(actor, services.CrearUsuarioInput{
		Username: "emp_sin_prop",
		Email:    "emp_sin_prop@test.com",
		Password: "123456",
		Rol:      models.RolEmpleado,
	})
	if !errors.Is(err, services.ErrPropietarioRequerido) {
		t.Fatalf("esperaba ErrPropietarioRequerido, obtuvo %v", err)
	}
}
