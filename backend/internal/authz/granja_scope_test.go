package authz_test

import (
	"testing"

	"github.com/martin/sgp/internal/authz"
	"github.com/martin/sgp/internal/testutil"
)

func TestGranjasAccesibles_EmpleadoVeGranjasDelPropietario(t *testing.T) {
	env := testutil.SetupTestDB(t)
	propID := testutil.SeedPropietario(t, env)
	empID := testutil.SeedEmpleado(t, env, propID)
	granjaID := testutil.SeedGranjaFor(t, env, propID)

	svc := authz.NewService(env.Repos)
	actor := testutil.SeedActor(t, env, empID)

	granjas, err := svc.GranjasAccesibles(actor, nil)
	if err != nil {
		t.Fatalf("GranjasAccesibles: %v", err)
	}
	if len(granjas) != 1 || granjas[0].ID != granjaID {
		t.Fatalf("empleado debería ver granja %d, obtuvo %v", granjaID, granjas)
	}
}

func TestRequireGranja_PropietarioNoAccedeGranjaAjena(t *testing.T) {
	env := testutil.SetupTestDB(t)
	propA := testutil.SeedPropietario(t, env)
	propB := testutil.SeedPropietario(t, env)
	granjaAjena := testutil.SeedGranjaFor(t, env, propB)

	svc := authz.NewService(env.Repos)
	actor := testutil.SeedActor(t, env, propA)

	if err := svc.RequireGranja(actor, granjaAjena, authz.PermGranjaRead); err == nil {
		t.Fatal("propietario no debería acceder a granja ajena")
	}
}

func TestRequireGranjaQuery_EmpleadoSinGranjaID(t *testing.T) {
	env := testutil.SetupTestDB(t)
	propID := testutil.SeedPropietario(t, env)
	empID := testutil.SeedEmpleado(t, env, propID)
	testutil.SeedGranjaFor(t, env, propID)

	svc := authz.NewService(env.Repos)
	actor := testutil.SeedActor(t, env, empID)

	if err := svc.RequireGranjaQuery(actor, nil, authz.PermGranjaRead); err == nil {
		t.Fatal("empleado debería requerir granja_id en query")
	}
}

func TestRequireGranja_AdminAccedeCualquierGranja(t *testing.T) {
	env := testutil.SetupTestDB(t)
	adminID := testutil.SeedAdmin(t, env)
	propID := testutil.SeedPropietario(t, env)
	granjaID := testutil.SeedGranjaFor(t, env, propID)

	svc := authz.NewService(env.Repos)
	actor := testutil.SeedActor(t, env, adminID)

	if err := svc.RequireGranja(actor, granjaID, authz.PermGranjaRead); err != nil {
		t.Fatalf("admin debería acceder a cualquier granja: %v", err)
	}
}
