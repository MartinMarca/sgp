package services_test

import (
	"testing"

	"github.com/martin/sgp/internal/models"
	"github.com/martin/sgp/internal/services"
	"github.com/martin/sgp/internal/testutil"
)

// ============================================================
// Tests del ciclo de maternidad completo
// ============================================================

func TestCicloMaternidadCompleto(t *testing.T) {
	env := testutil.SetupTestDB(t)
	svc := env.Services

	granjaID := testutil.SeedGranja(t, svc)
	corralID := testutil.SeedCorral(t, svc, granjaID)
	padrilloID := testutil.SeedPadrillo(t, svc, granjaID)
	cerdaID := testutil.SeedCerda(t, svc, granjaID)

	// 1. Verificar estado inicial: disponible
	cerda, err := svc.Cerda.ObtenerPorID(cerdaID)
	if err != nil {
		t.Fatalf("Error obteniendo cerda: %v", err)
	}
	if cerda.Estado != models.EstadoCerdaDisponible {
		t.Fatalf("Estado inicial: esperaba 'disponible', obtuvo '%s'", cerda.Estado)
	}

	// 2. Registrar servicio: disponible -> servicio
	servicio, err := svc.Servicio.Crear(services.CrearServicioInput{
		CerdaID:       cerdaID,
		FechaServicio: "2026-01-15",
		TipoMonta:     "natural",
		PadrilloID:    &padrilloID,
	})
	if err != nil {
		t.Fatalf("Error creando servicio: %v", err)
	}
	if servicio.ID == 0 {
		t.Fatal("Servicio debería tener un ID")
	}

	cerda, _ = svc.Cerda.ObtenerPorID(cerdaID)
	if cerda.Estado != models.EstadoCerdaServicio {
		t.Fatalf("Después de servicio: esperaba 'servicio', obtuvo '%s'", cerda.Estado)
	}

	// 3. Confirmar preñez: servicio -> gestación
	_, err = svc.Servicio.ConfirmarPrenez(servicio.ID, services.ConfirmarPrenezInput{
		FechaConfirmacion: "2026-01-25",
	})
	if err != nil {
		t.Fatalf("Error confirmando preñez: %v", err)
	}

	cerda, _ = svc.Cerda.ObtenerPorID(cerdaID)
	if cerda.Estado != models.EstadoCerdaGestacion {
		t.Fatalf("Después de confirmar: esperaba 'gestacion', obtuvo '%s'", cerda.Estado)
	}

	// 4. Registrar parto: gestación -> cría
	parto, err := svc.Parto.Crear(services.CrearPartoInput{
		CerdaID:                cerdaID,
		FechaParto:             "2026-05-09",
		LechonesNacidosVivos:   10,
		LechonesNacidosTotales: 11,
		LechonesHembras:        6,
		LechonesMachos:         4,
	})
	if err != nil {
		t.Fatalf("Error creando parto: %v", err)
	}
	if parto.FechaEstimada.IsZero() {
		t.Fatal("FechaEstimada del parto no debería ser cero")
	}

	cerda, _ = svc.Cerda.ObtenerPorID(cerdaID)
	if cerda.Estado != models.EstadoCerdaCria {
		t.Fatalf("Después de parto: esperaba 'cria', obtuvo '%s'", cerda.Estado)
	}

	// 5. Registrar destete con nuevo lote: cría -> disponible
	destete, err := svc.Destete.Crear(services.CrearDesteteInput{
		CerdaID:                    cerdaID,
		FechaDestete:               "2026-06-08",
		CantidadLechonesDestetados: 9,
		NuevoLote: &services.NuevoLoteInput{
			CorralID: corralID,
			Nombre:   "Lote Test",
		},
	})
	if err != nil {
		t.Fatalf("Error creando destete: %v", err)
	}
	if destete.LoteID == 0 {
		t.Fatal("Destete debería tener un LoteID")
	}

	cerda, _ = svc.Cerda.ObtenerPorID(cerdaID)
	if cerda.Estado != models.EstadoCerdaDisponible {
		t.Fatalf("Después de destete: esperaba 'disponible', obtuvo '%s'", cerda.Estado)
	}

	// 6. Verificar lote creado
	lote, err := svc.Lote.ObtenerPorID(destete.LoteID)
	if err != nil {
		t.Fatalf("Error obteniendo lote: %v", err)
	}
	if lote.CantidadLechones != 9 {
		t.Fatalf("Lote: esperaba 9 lechones, obtuvo %d", lote.CantidadLechones)
	}
	if lote.Nombre != "Lote Test" {
		t.Fatalf("Lote nombre: esperaba 'Lote Test', obtuvo '%s'", lote.Nombre)
	}
}

// ============================================================
// Tests de validaciones de negocio - Cerda
// ============================================================

func TestCrearCerda_CaravanaDuplicada(t *testing.T) {
	env := testutil.SetupTestDB(t)
	svc := env.Services
	granjaID := testutil.SeedGranja(t, svc)

	_, err := svc.Cerda.Crear(services.CrearCerdaInput{
		GranjaID: granjaID, NumeroCaravana: "C-DUP", Estado: "disponible",
	})
	if err != nil {
		t.Fatalf("Primera cerda no debería fallar: %v", err)
	}

	_, err = svc.Cerda.Crear(services.CrearCerdaInput{
		GranjaID: granjaID, NumeroCaravana: "C-DUP", Estado: "disponible",
	})
	if err == nil {
		t.Fatal("Segunda cerda con misma caravana debería fallar")
	}
	if err != services.ErrCaravanaDuplicada {
		t.Fatalf("Esperaba ErrCaravanaDuplicada, obtuvo: %v", err)
	}
}

func TestCrearCerda_CualquierEstadoInicial(t *testing.T) {
	env := testutil.SetupTestDB(t)
	svc := env.Services
	granjaID := testutil.SeedGranja(t, svc)

	estados := []string{"disponible", "servicio", "gestacion", "cria"}
	for i, estado := range estados {
		cerda, err := svc.Cerda.Crear(services.CrearCerdaInput{
			GranjaID:       granjaID,
			NumeroCaravana: "EST-" + string(rune('A'+i)),
			Estado:         estado,
		})
		if err != nil {
			t.Fatalf("Crear cerda en estado '%s' no debería fallar: %v", estado, err)
		}
		if cerda.Estado != estado {
			t.Fatalf("Cerda: esperaba estado '%s', obtuvo '%s'", estado, cerda.Estado)
		}
	}
}

func TestBajaCerda_EnServicioFalla(t *testing.T) {
	env := testutil.SetupTestDB(t)
	svc := env.Services
	granjaID := testutil.SeedGranja(t, svc)
	padrilloID := testutil.SeedPadrillo(t, svc, granjaID)
	cerdaID := testutil.SeedCerda(t, svc, granjaID)

	// Pasar a estado servicio
	svc.Servicio.Crear(services.CrearServicioInput{
		CerdaID: cerdaID, FechaServicio: "2026-01-15",
		TipoMonta: "natural", PadrilloID: &padrilloID,
	})

	// Intentar dar de baja
	_, err := svc.Cerda.DarDeBaja(cerdaID, services.BajaCerdaInput{MotivoBaja: "muerte"})
	if err == nil {
		t.Fatal("Baja de cerda en servicio debería fallar")
	}
	if err != services.ErrCerdaTieneServicioActivo {
		t.Fatalf("Esperaba ErrCerdaTieneServicioActivo, obtuvo: %v", err)
	}
}

// ============================================================
// Tests de validaciones de negocio - Servicio
// ============================================================

func TestCrearServicio_CerdaNoDisponible(t *testing.T) {
	env := testutil.SetupTestDB(t)
	svc := env.Services
	granjaID := testutil.SeedGranja(t, svc)

	// Crear cerda en estado gestación
	cerda, _ := svc.Cerda.Crear(services.CrearCerdaInput{
		GranjaID: granjaID, NumeroCaravana: "C-GEST", Estado: "gestacion",
	})

	_, err := svc.Servicio.Crear(services.CrearServicioInput{
		CerdaID: cerda.ID, FechaServicio: "2026-01-15",
		TipoMonta: "inseminacion", NumeroPajuela: strPtr("PAJ-001"),
	})
	if err == nil {
		t.Fatal("Servicio en cerda no disponible debería fallar")
	}
}

func TestCrearServicio_MontaNaturalSinPadrillo(t *testing.T) {
	env := testutil.SetupTestDB(t)
	svc := env.Services
	granjaID := testutil.SeedGranja(t, svc)
	cerdaID := testutil.SeedCerda(t, svc, granjaID)

	_, err := svc.Servicio.Crear(services.CrearServicioInput{
		CerdaID: cerdaID, FechaServicio: "2026-01-15",
		TipoMonta: "natural", PadrilloID: nil,
	})
	if err == nil {
		t.Fatal("Monta natural sin padrillo debería fallar")
	}
}

func TestConfirmarPrenez_DobleConfirmacion(t *testing.T) {
	env := testutil.SetupTestDB(t)
	svc := env.Services
	granjaID := testutil.SeedGranja(t, svc)
	padrilloID := testutil.SeedPadrillo(t, svc, granjaID)
	cerdaID := testutil.SeedCerda(t, svc, granjaID)

	servicio, _ := svc.Servicio.Crear(services.CrearServicioInput{
		CerdaID: cerdaID, FechaServicio: "2026-01-15",
		TipoMonta: "natural", PadrilloID: &padrilloID,
	})

	// Primera confirmación
	_, err := svc.Servicio.ConfirmarPrenez(servicio.ID, services.ConfirmarPrenezInput{})
	if err != nil {
		t.Fatalf("Primera confirmación no debería fallar: %v", err)
	}

	// Segunda confirmación debería fallar
	_, err = svc.Servicio.ConfirmarPrenez(servicio.ID, services.ConfirmarPrenezInput{})
	if err == nil {
		t.Fatal("Doble confirmación debería fallar")
	}
}

func TestCancelarPrenez_DevuelveDisponible(t *testing.T) {
	env := testutil.SetupTestDB(t)
	svc := env.Services
	granjaID := testutil.SeedGranja(t, svc)
	padrilloID := testutil.SeedPadrillo(t, svc, granjaID)
	cerdaID := testutil.SeedCerda(t, svc, granjaID)

	servicio, _ := svc.Servicio.Crear(services.CrearServicioInput{
		CerdaID: cerdaID, FechaServicio: "2026-01-15",
		TipoMonta: "natural", PadrilloID: &padrilloID,
	})

	svc.Servicio.ConfirmarPrenez(servicio.ID, services.ConfirmarPrenezInput{})

	// Cancelar
	_, err := svc.Servicio.CancelarPrenez(servicio.ID, services.CancelarPrenezInput{
		Motivo: "Aborto espontáneo",
	})
	if err != nil {
		t.Fatalf("Cancelar preñez no debería fallar: %v", err)
	}

	cerda, _ := svc.Cerda.ObtenerPorID(cerdaID)
	if cerda.Estado != models.EstadoCerdaDisponible {
		t.Fatalf("Después de cancelar: esperaba 'disponible', obtuvo '%s'", cerda.Estado)
	}
}

// ============================================================
// Tests de validaciones de negocio - Parto
// ============================================================

func TestCrearParto_CerdaNoEnGestacion(t *testing.T) {
	env := testutil.SetupTestDB(t)
	svc := env.Services
	granjaID := testutil.SeedGranja(t, svc)
	cerdaID := testutil.SeedCerda(t, svc, granjaID) // Estado disponible

	_, err := svc.Parto.Crear(services.CrearPartoInput{
		CerdaID: cerdaID, FechaParto: "2026-05-09",
		LechonesNacidosVivos: 5, LechonesNacidosTotales: 5,
		LechonesHembras: 3, LechonesMachos: 2,
	})
	if err == nil {
		t.Fatal("Parto en cerda disponible debería fallar")
	}
}

func TestCrearParto_HembrasMachosDistintoVivos(t *testing.T) {
	env := testutil.SetupTestDB(t)
	svc := env.Services
	granjaID := testutil.SeedGranja(t, svc)
	padrilloID := testutil.SeedPadrillo(t, svc, granjaID)
	cerdaID := testutil.SeedCerda(t, svc, granjaID)

	// Llevar a gestación
	srv, _ := svc.Servicio.Crear(services.CrearServicioInput{
		CerdaID: cerdaID, FechaServicio: "2026-01-15",
		TipoMonta: "natural", PadrilloID: &padrilloID,
	})
	svc.Servicio.ConfirmarPrenez(srv.ID, services.ConfirmarPrenezInput{})

	// Intentar parto con h+m != vivos
	_, err := svc.Parto.Crear(services.CrearPartoInput{
		CerdaID: cerdaID, FechaParto: "2026-05-09",
		LechonesNacidosVivos: 8, LechonesNacidosTotales: 10,
		LechonesHembras: 3, LechonesMachos: 3, // 3+3=6 != 8
	})
	if err == nil {
		t.Fatal("Parto con h+m != vivos debería fallar")
	}
	if err != services.ErrLechonesInvalidos {
		t.Fatalf("Esperaba ErrLechonesInvalidos, obtuvo: %v", err)
	}
}

// ============================================================
// Tests de validaciones de negocio - Destete
// ============================================================

func TestCrearDestete_DestetadosExcedenVivos(t *testing.T) {
	env := testutil.SetupTestDB(t)
	svc := env.Services
	granjaID := testutil.SeedGranja(t, svc)
	corralID := testutil.SeedCorral(t, svc, granjaID)
	padrilloID := testutil.SeedPadrillo(t, svc, granjaID)
	cerdaID := testutil.SeedCerda(t, svc, granjaID)

	// Ciclo completo hasta cría
	srv, _ := svc.Servicio.Crear(services.CrearServicioInput{
		CerdaID: cerdaID, FechaServicio: "2026-01-15",
		TipoMonta: "natural", PadrilloID: &padrilloID,
	})
	svc.Servicio.ConfirmarPrenez(srv.ID, services.ConfirmarPrenezInput{})
	svc.Parto.Crear(services.CrearPartoInput{
		CerdaID: cerdaID, FechaParto: "2026-05-09",
		LechonesNacidosVivos: 8, LechonesNacidosTotales: 10,
		LechonesHembras: 4, LechonesMachos: 4,
	})

	// Intentar destetar más de los nacidos vivos
	_, err := svc.Destete.Crear(services.CrearDesteteInput{
		CerdaID: cerdaID, FechaDestete: "2026-06-08",
		CantidadLechonesDestetados: 15, // Más que 8 nacidos vivos
		NuevoLote: &services.NuevoLoteInput{
			CorralID: corralID, Nombre: "Lote Exceso",
		},
	})
	if err == nil {
		t.Fatal("Destete con más destetados que vivos debería fallar")
	}
}

func TestCrearDestete_SinLoteFalla(t *testing.T) {
	env := testutil.SetupTestDB(t)
	svc := env.Services
	granjaID := testutil.SeedGranja(t, svc)
	padrilloID := testutil.SeedPadrillo(t, svc, granjaID)
	cerdaID := testutil.SeedCerda(t, svc, granjaID)

	// Ciclo hasta cría
	srv, _ := svc.Servicio.Crear(services.CrearServicioInput{
		CerdaID: cerdaID, FechaServicio: "2026-01-15",
		TipoMonta: "natural", PadrilloID: &padrilloID,
	})
	svc.Servicio.ConfirmarPrenez(srv.ID, services.ConfirmarPrenezInput{})
	svc.Parto.Crear(services.CrearPartoInput{
		CerdaID: cerdaID, FechaParto: "2026-05-09",
		LechonesNacidosVivos: 8, LechonesNacidosTotales: 8,
		LechonesHembras: 4, LechonesMachos: 4,
	})

	// Destete sin lote
	_, err := svc.Destete.Crear(services.CrearDesteteInput{
		CerdaID: cerdaID, FechaDestete: "2026-06-08",
		CantidadLechonesDestetados: 8,
		LoteID: nil, NuevoLote: nil,
	})
	if err == nil {
		t.Fatal("Destete sin lote debería fallar")
	}
	if err != services.ErrLoteRequerido {
		t.Fatalf("Esperaba ErrLoteRequerido, obtuvo: %v", err)
	}
}

func TestCrearDestete_ConLoteExistente(t *testing.T) {
	env := testutil.SetupTestDB(t)
	svc := env.Services
	granjaID := testutil.SeedGranja(t, svc)
	corralID := testutil.SeedCorral(t, svc, granjaID)
	padrilloID := testutil.SeedPadrillo(t, svc, granjaID)

	// Crear lote anticipado
	lote, err := svc.Lote.Crear(services.CrearLoteInput{
		CorralID: corralID, Nombre: "Lote Previo", Fecha: "2026-01-01",
	})
	if err != nil {
		t.Fatalf("Error creando lote: %v", err)
	}
	if lote.CantidadLechones != 0 {
		t.Fatalf("Lote nuevo debería tener 0 lechones, tiene %d", lote.CantidadLechones)
	}

	// Ciclo completo de una cerda hasta destete
	cerdaID := testutil.SeedCerda(t, svc, granjaID)
	srv, _ := svc.Servicio.Crear(services.CrearServicioInput{
		CerdaID: cerdaID, FechaServicio: "2026-01-15",
		TipoMonta: "natural", PadrilloID: &padrilloID,
	})
	svc.Servicio.ConfirmarPrenez(srv.ID, services.ConfirmarPrenezInput{})
	svc.Parto.Crear(services.CrearPartoInput{
		CerdaID: cerdaID, FechaParto: "2026-05-09",
		LechonesNacidosVivos: 10, LechonesNacidosTotales: 10,
		LechonesHembras: 5, LechonesMachos: 5,
	})

	// Destete asignando a lote existente
	loteID := lote.ID
	_, err = svc.Destete.Crear(services.CrearDesteteInput{
		CerdaID: cerdaID, FechaDestete: "2026-06-08",
		CantidadLechonesDestetados: 9,
		LoteID: &loteID,
	})
	if err != nil {
		t.Fatalf("Destete con lote existente no debería fallar: %v", err)
	}

	// Verificar que el lote tiene los lechones
	loteActualizado, _ := svc.Lote.ObtenerPorID(loteID)
	if loteActualizado.CantidadLechones != 9 {
		t.Fatalf("Lote debería tener 9 lechones, tiene %d", loteActualizado.CantidadLechones)
	}
}

// ============================================================
// Tests de Lote y Corral
// ============================================================

func TestCerrarLote(t *testing.T) {
	env := testutil.SetupTestDB(t)
	svc := env.Services
	granjaID := testutil.SeedGranja(t, svc)
	corralID := testutil.SeedCorral(t, svc, granjaID)

	lote, _ := svc.Lote.Crear(services.CrearLoteInput{
		CorralID: corralID, Nombre: "Lote a cerrar",
	})

	// Cerrar lote
	cerrado, err := svc.Lote.Cerrar(lote.ID, services.CerrarLoteInput{
		MotivoCierre: "Vendidos", Estado: "vendido",
	})
	if err != nil {
		t.Fatalf("Cerrar lote no debería fallar: %v", err)
	}
	if cerrado.Estado != models.EstadoLoteVendido {
		t.Fatalf("Estado: esperaba 'vendido', obtuvo '%s'", cerrado.Estado)
	}

	// Intentar cerrar de nuevo
	_, err = svc.Lote.Cerrar(lote.ID, services.CerrarLoteInput{
		MotivoCierre: "Otro", Estado: "cerrado",
	})
	if err == nil {
		t.Fatal("Cerrar lote ya cerrado debería fallar")
	}
}

func TestEliminarCorral_ConLotesActivosFalla(t *testing.T) {
	env := testutil.SetupTestDB(t)
	svc := env.Services
	granjaID := testutil.SeedGranja(t, svc)
	corralID := testutil.SeedCorral(t, svc, granjaID)

	// Crear lote activo en el corral
	svc.Lote.Crear(services.CrearLoteInput{
		CorralID: corralID, Nombre: "Lote activo",
	})

	err := svc.Corral.Eliminar(corralID)
	if err == nil {
		t.Fatal("Eliminar corral con lotes activos debería fallar")
	}
}

// ============================================================
// Tests de Auth
// ============================================================

func TestAuth_RegistroYLogin(t *testing.T) {
	env := testutil.SetupTestDB(t)
	svc := env.Services

	// Registrar
	usuario, err := svc.Auth.Registrar(services.RegistroInput{
		Username:       "testuser",
		Email:          "test@test.com",
		Password:       "password123",
		NombreCompleto: "Test User",
	})
	if err != nil {
		t.Fatalf("Registro no debería fallar: %v", err)
	}
	if usuario.Username != "testuser" {
		t.Fatalf("Username: esperaba 'testuser', obtuvo '%s'", usuario.Username)
	}
	if usuario.PasswordHash != "" {
		t.Fatal("PasswordHash no debería estar en la respuesta")
	}

	// Login
	response, err := svc.Auth.Login(services.LoginInput{
		Username: "testuser", Password: "password123",
	})
	if err != nil {
		t.Fatalf("Login no debería fallar: %v", err)
	}
	if response.Token == "" {
		t.Fatal("Token no debería estar vacío")
	}

	// Login con password incorrecto
	_, err = svc.Auth.Login(services.LoginInput{
		Username: "testuser", Password: "wrongpass",
	})
	if err == nil {
		t.Fatal("Login con password incorrecto debería fallar")
	}
}

func TestAuth_RegistroDuplicado(t *testing.T) {
	env := testutil.SetupTestDB(t)
	svc := env.Services

	svc.Auth.Registrar(services.RegistroInput{
		Username: "dup", Email: "dup@test.com", Password: "123456",
	})

	_, err := svc.Auth.Registrar(services.RegistroInput{
		Username: "dup", Email: "other@test.com", Password: "123456",
	})
	if err == nil {
		t.Fatal("Registro con username duplicado debería fallar")
	}
}

// ============================================================
// Helpers
// ============================================================

func strPtr(s string) *string {
	return &s
}
