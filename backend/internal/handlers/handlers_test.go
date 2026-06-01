package handlers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/martin/sgp/internal/config"
	"github.com/martin/sgp/internal/routes"
	"github.com/martin/sgp/internal/services"
	"github.com/martin/sgp/internal/testutil"
	"github.com/martin/sgp/internal/utils"
)

// setupRouter crea un router con toda la infraestructura de test
func setupRouter(t *testing.T) (*testutil.TestEnv, http.Handler) {
	t.Helper()
	env := testutil.SetupTestDB(t)

	cfg := &config.Config{
		JWTSecret:         "test-secret-key",
		JWTExpiration:     "1h",
		CORSOrigin:        "*",
		RegisterEnabled:   true,
	}

	router := routes.SetupRoutes(cfg, env.Services, env.Repos)
	return env, router
}

// getAuthToken registra y logea un usuario, retorna el token
func getAuthToken(t *testing.T, router http.Handler) string {
	t.Helper()

	// Registrar
	body := `{"username":"testadmin","email":"admin@test.com","password":"123456","nombre_completo":"Admin Test"}`
	req := httptest.NewRequest("POST", "/api/auth/register", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("Register: esperaba 201, obtuvo %d: %s", w.Code, w.Body.String())
	}

	// Login
	loginBody := `{"username":"testadmin","password":"123456"}`
	req = httptest.NewRequest("POST", "/api/auth/login", bytes.NewBufferString(loginBody))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Login: esperaba 200, obtuvo %d: %s", w.Code, w.Body.String())
	}

	var resp utils.Response
	json.Unmarshal(w.Body.Bytes(), &resp)
	data := resp.Data.(map[string]interface{})
	return data["token"].(string)
}

// doRequest helper para hacer requests autenticados
func doRequest(t *testing.T, router http.Handler, method, path, token string, body interface{}) *httptest.ResponseRecorder {
	t.Helper()
	var reqBody *bytes.Buffer
	if body != nil {
		b, _ := json.Marshal(body)
		reqBody = bytes.NewBuffer(b)
	} else {
		reqBody = bytes.NewBuffer(nil)
	}

	req := httptest.NewRequest(method, path, reqBody)
	req.Header.Set("Content-Type", "application/json")
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	return w
}

// parseResponse parsea un Response estándar
func parseResponse(t *testing.T, w *httptest.ResponseRecorder) utils.Response {
	t.Helper()
	var resp utils.Response
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("Error parseando response: %v\nBody: %s", err, w.Body.String())
	}
	return resp
}

// ============================================================
// Tests E2E - Health Check
// ============================================================

func TestHealthCheck(t *testing.T) {
	_, router := setupRouter(t)

	w := doRequest(t, router, "GET", "/api/health", "", nil)
	if w.Code != http.StatusOK {
		t.Fatalf("Health: esperaba 200, obtuvo %d", w.Code)
	}
}

// ============================================================
// Tests E2E - Auth
// ============================================================

func TestE2E_RegistroYLogin(t *testing.T) {
	_, router := setupRouter(t)

	// Registro
	w := doRequest(t, router, "POST", "/api/auth/register", "", map[string]string{
		"username": "e2euser", "email": "e2e@test.com",
		"password": "123456", "nombre_completo": "E2E User",
	})
	if w.Code != http.StatusCreated {
		t.Fatalf("Register: esperaba 201, obtuvo %d: %s", w.Code, w.Body.String())
	}

	// Login
	w = doRequest(t, router, "POST", "/api/auth/login", "", map[string]string{
		"username": "e2euser", "password": "123456",
	})
	if w.Code != http.StatusOK {
		t.Fatalf("Login: esperaba 200, obtuvo %d", w.Code)
	}

	resp := parseResponse(t, w)
	if !resp.Success {
		t.Fatal("Login debería ser exitoso")
	}
}

func TestE2E_AccesoSinToken(t *testing.T) {
	_, router := setupRouter(t)

	w := doRequest(t, router, "GET", "/api/granjas", "", nil)
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("Sin token: esperaba 401, obtuvo %d", w.Code)
	}
}

// ============================================================
// Tests E2E - CRUD Granjas
// ============================================================

func TestE2E_CRUDGranjas(t *testing.T) {
	_, router := setupRouter(t)
	token := getAuthToken(t, router)

	// Crear granja
	w := doRequest(t, router, "POST", "/api/granjas", token, services.CrearGranjaInput{
		Nombre: "Granja E2E", Descripcion: "Desc", Ubicacion: "Lugar",
	})
	if w.Code != http.StatusCreated {
		t.Fatalf("Crear granja: esperaba 201, obtuvo %d: %s", w.Code, w.Body.String())
	}

	// Listar granjas
	w = doRequest(t, router, "GET", "/api/granjas", token, nil)
	if w.Code != http.StatusOK {
		t.Fatalf("Listar granjas: esperaba 200, obtuvo %d", w.Code)
	}

	// Obtener granja por ID
	w = doRequest(t, router, "GET", "/api/granjas/1", token, nil)
	if w.Code != http.StatusOK {
		t.Fatalf("Obtener granja: esperaba 200, obtuvo %d", w.Code)
	}

	// Actualizar granja
	w = doRequest(t, router, "PUT", "/api/granjas/1", token, services.ActualizarGranjaInput{
		Nombre: "Granja Actualizada",
	})
	if w.Code != http.StatusOK {
		t.Fatalf("Actualizar granja: esperaba 200, obtuvo %d: %s", w.Code, w.Body.String())
	}
}

// ============================================================
// Tests E2E - Ciclo de maternidad HTTP
// ============================================================

func TestE2E_CicloMaternidad(t *testing.T) {
	_, router := setupRouter(t)
	token := getAuthToken(t, router)

	// 1. Crear granja
	w := doRequest(t, router, "POST", "/api/granjas", token, map[string]string{
		"nombre": "Granja Ciclo",
	})
	if w.Code != http.StatusCreated {
		t.Fatalf("Crear granja: %d %s", w.Code, w.Body.String())
	}

	// 2. Crear corral
	w = doRequest(t, router, "POST", "/api/granjas/1/corrales", token, map[string]string{
		"nombre": "Corral 1",
	})
	if w.Code != http.StatusCreated {
		t.Fatalf("Crear corral: %d %s", w.Code, w.Body.String())
	}

	// 3. Crear padrillo
	w = doRequest(t, router, "POST", "/api/granjas/1/padrillos", token, map[string]interface{}{
		"numero_caravana": "P-E2E", "nombre": "Padrillo E2E",
	})
	if w.Code != http.StatusCreated {
		t.Fatalf("Crear padrillo: %d %s", w.Code, w.Body.String())
	}

	// 4. Crear cerda
	w = doRequest(t, router, "POST", "/api/granjas/1/cerdas", token, map[string]interface{}{
		"numero_caravana": "C-E2E", "estado": "disponible",
	})
	if w.Code != http.StatusCreated {
		t.Fatalf("Crear cerda: %d %s", w.Code, w.Body.String())
	}

	// 5. Registrar servicio
	padrilloID := 1
	w = doRequest(t, router, "POST", "/api/servicios", token, map[string]interface{}{
		"cerda_id": 1, "fecha_servicio": "2026-01-15",
		"tipo_monta": "natural", "padrillo_id": padrilloID,
	})
	if w.Code != http.StatusCreated {
		t.Fatalf("Crear servicio: %d %s", w.Code, w.Body.String())
	}

	// 6. Confirmar preñez
	w = doRequest(t, router, "POST", "/api/servicios/1/confirmar-prenez", token, map[string]string{
		"fecha_confirmacion": "2026-01-25",
	})
	if w.Code != http.StatusOK {
		t.Fatalf("Confirmar preñez: %d %s", w.Code, w.Body.String())
	}

	// 7. Registrar parto
	w = doRequest(t, router, "POST", "/api/partos", token, map[string]interface{}{
		"cerda_id": 1, "fecha_parto": "2026-05-09",
		"lechones_nacidos_vivos": 10, "lechones_nacidos_totales": 11,
		"lechones_hembras": 6, "lechones_machos": 4,
	})
	if w.Code != http.StatusCreated {
		t.Fatalf("Crear parto: %d %s", w.Code, w.Body.String())
	}

	// 8. Registrar destete
	w = doRequest(t, router, "POST", "/api/destetes", token, map[string]interface{}{
		"cerda_id": 1, "fecha_destete": "2026-06-08",
		"cantidad_lechones_destetados": 9,
		"nuevo_lote": map[string]interface{}{
			"corral_id": 1, "nombre": "Lote E2E",
		},
	})
	if w.Code != http.StatusCreated {
		t.Fatalf("Crear destete: %d %s", w.Code, w.Body.String())
	}

	// 9. Verificar cerda está disponible de nuevo
	w = doRequest(t, router, "GET", "/api/cerdas/1", token, nil)
	if w.Code != http.StatusOK {
		t.Fatalf("Obtener cerda: %d", w.Code)
	}
	resp := parseResponse(t, w)
	data := resp.Data.(map[string]interface{})
	if data["estado"] != "disponible" {
		t.Fatalf("Cerda debería estar disponible, está en '%v'", data["estado"])
	}

	// 10. Verificar lote
	w = doRequest(t, router, "GET", "/api/lotes/1", token, nil)
	if w.Code != http.StatusOK {
		t.Fatalf("Obtener lote: %d", w.Code)
	}
	resp = parseResponse(t, w)
	loteData := resp.Data.(map[string]interface{})
	if int(loteData["cantidad_lechones"].(float64)) != 9 {
		t.Fatalf("Lote debería tener 9 lechones, tiene %v", loteData["cantidad_lechones"])
	}

	// 11. Calendario
	w = doRequest(t, router, "GET", "/api/calendario?granja_id=1&dias=365", token, nil)
	if w.Code != http.StatusOK {
		t.Fatalf("Calendario: %d %s", w.Code, w.Body.String())
	}
}

// ============================================================
// Tests E2E - Validaciones HTTP
// ============================================================

func TestE2E_ValidacionPartoHMDistintoVivos(t *testing.T) {
	_, router := setupRouter(t)
	token := getAuthToken(t, router)

	// Setup: granja + cerda en gestación
	doRequest(t, router, "POST", "/api/granjas", token, map[string]string{"nombre": "G"})
	doRequest(t, router, "POST", "/api/granjas/1/padrillos", token, map[string]interface{}{
		"numero_caravana": "P1", "nombre": "Padrillo",
	})
	doRequest(t, router, "POST", "/api/granjas/1/cerdas", token, map[string]interface{}{
		"numero_caravana": "C1", "estado": "disponible",
	})
	doRequest(t, router, "POST", "/api/servicios", token, map[string]interface{}{
		"cerda_id": 1, "fecha_servicio": "2026-01-15",
		"tipo_monta": "natural", "padrillo_id": 1,
	})
	doRequest(t, router, "POST", "/api/servicios/1/confirmar-prenez", token, map[string]string{})

	// Parto con h+m != vivos
	w := doRequest(t, router, "POST", "/api/partos", token, map[string]interface{}{
		"cerda_id": 1, "fecha_parto": "2026-05-09",
		"lechones_nacidos_vivos": 8, "lechones_nacidos_totales": 10,
		"lechones_hembras": 3, "lechones_machos": 3,
	})
	if w.Code != http.StatusUnprocessableEntity {
		t.Fatalf("Parto inválido: esperaba 422, obtuvo %d: %s", w.Code, w.Body.String())
	}
}

// ============================================================
// Tests E2E - Usuarios (roles)
// ============================================================

func TestE2E_UsuariosRoles(t *testing.T) {
	_, router := setupRouter(t)
	tokenProp := getAuthToken(t, router)

	// Propietario crea empleado
	w := doRequest(t, router, "POST", "/api/usuarios", tokenProp, map[string]interface{}{
		"username": "empleado1",
		"email":    "empleado1@test.com",
		"password": "123456",
		"rol":      "empleado",
	})
	if w.Code != http.StatusCreated {
		t.Fatalf("Crear empleado: esperaba 201, obtuvo %d: %s", w.Code, w.Body.String())
	}

	// Propietario no puede crear admin
	w = doRequest(t, router, "POST", "/api/usuarios", tokenProp, map[string]interface{}{
		"username": "hacker",
		"email":    "hacker@test.com",
		"password": "123456",
		"rol":      "admin",
	})
	if w.Code != http.StatusForbidden {
		t.Fatalf("Crear admin como propietario: esperaba 403, obtuvo %d", w.Code)
	}

	// Propietario lista empleados
	w = doRequest(t, router, "GET", "/api/usuarios", tokenProp, nil)
	if w.Code != http.StatusOK {
		t.Fatalf("Listar usuarios propietario: esperaba 200, obtuvo %d", w.Code)
	}
	resp := parseResponse(t, w)
	list := resp.Data.([]interface{})
	if len(list) != 1 {
		t.Fatalf("Propietario debería ver 1 empleado, obtuvo %d", len(list))
	}

	// Login empleado
	w = doRequest(t, router, "POST", "/api/auth/login", "", map[string]string{
		"username": "empleado1", "password": "123456",
	})
	if w.Code != http.StatusOK {
		t.Fatalf("Login empleado: %d %s", w.Code, w.Body.String())
	}
	resp = parseResponse(t, w)
	tokenEmp := resp.Data.(map[string]interface{})["token"].(string)

	// Empleado no puede listar usuarios
	w = doRequest(t, router, "GET", "/api/usuarios", tokenEmp, nil)
	if w.Code != http.StatusForbidden {
		t.Fatalf("Empleado listar usuarios: esperaba 403, obtuvo %d", w.Code)
	}

	// Empleado puede ver su propio perfil
	w = doRequest(t, router, "GET", "/api/usuarios/2", tokenEmp, nil)
	if w.Code != http.StatusOK {
		t.Fatalf("Empleado ver self: esperaba 200, obtuvo %d: %s", w.Code, w.Body.String())
	}
}

func TestE2E_AuthzRestricciones(t *testing.T) {
	_, router := setupRouter(t)
	tokenProp := getAuthToken(t, router)

	// Granja y cerda del propietario
	w := doRequest(t, router, "POST", "/api/granjas", tokenProp, map[string]string{"nombre": "Granja Authz"})
	if w.Code != http.StatusCreated {
		t.Fatalf("Crear granja: %d %s", w.Code, w.Body.String())
	}
	w = doRequest(t, router, "POST", "/api/granjas/1/cerdas", tokenProp, map[string]interface{}{
		"numero_caravana": "C-AUTHZ", "estado": "disponible",
	})
	if w.Code != http.StatusCreated {
		t.Fatalf("Crear cerda: %d %s", w.Code, w.Body.String())
	}

	// Segundo propietario con granja ajena
	w = doRequest(t, router, "POST", "/api/auth/register", "", map[string]string{
		"username": "otroprop", "email": "otro@test.com", "password": "123456",
	})
	if w.Code != http.StatusCreated {
		t.Fatalf("Register otro prop: %d %s", w.Code, w.Body.String())
	}
	w = doRequest(t, router, "POST", "/api/auth/login", "", map[string]string{
		"username": "otroprop", "password": "123456",
	})
	resp := parseResponse(t, w)
	tokenOtro := resp.Data.(map[string]interface{})["token"].(string)
	w = doRequest(t, router, "POST", "/api/granjas", tokenOtro, map[string]string{"nombre": "Granja Ajena"})
	if w.Code != http.StatusCreated {
		t.Fatalf("Crear granja ajena: %d %s", w.Code, w.Body.String())
	}

	// Propietario crea empleado
	w = doRequest(t, router, "POST", "/api/usuarios", tokenProp, map[string]interface{}{
		"username": "empauthz", "email": "empauthz@test.com", "password": "123456", "rol": "empleado",
	})
	if w.Code != http.StatusCreated {
		t.Fatalf("Crear empleado: %d %s", w.Code, w.Body.String())
	}
	w = doRequest(t, router, "POST", "/api/auth/login", "", map[string]string{
		"username": "empauthz", "password": "123456",
	})
	resp = parseResponse(t, w)
	tokenEmp := resp.Data.(map[string]interface{})["token"].(string)

	// Empleado no puede acceder a ventas (middleware)
	w = doRequest(t, router, "GET", "/api/ventas?granja_id=1", tokenEmp, nil)
	if w.Code != http.StatusForbidden {
		t.Fatalf("Empleado ventas: esperaba 403, obtuvo %d", w.Code)
	}

	// Empleado no puede dar de baja cerda
	w = doRequest(t, router, "POST", "/api/cerdas/1/baja", tokenEmp, map[string]string{
		"motivo_baja": "venta",
	})
	if w.Code != http.StatusForbidden {
		t.Fatalf("Empleado baja cerda: esperaba 403, obtuvo %d", w.Code)
	}

	// Propietario sí puede listar ventas (vacío)
	w = doRequest(t, router, "GET", "/api/ventas?granja_id=1", tokenProp, nil)
	if w.Code != http.StatusOK {
		t.Fatalf("Propietario ventas: esperaba 200, obtuvo %d: %s", w.Code, w.Body.String())
	}

	// Propietario no accede a granja ajena (id=2)
	w = doRequest(t, router, "GET", "/api/granjas/2", tokenProp, nil)
	if w.Code != http.StatusForbidden {
		t.Fatalf("Granja ajena: esperaba 403, obtuvo %d", w.Code)
	}
}
