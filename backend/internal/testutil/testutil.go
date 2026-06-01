package testutil

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"testing"
	"time"

	"github.com/martin/sgp/internal/authz"
	"github.com/martin/sgp/internal/models"
	"github.com/martin/sgp/internal/repositories"
	"github.com/martin/sgp/internal/services"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// TestEnv contiene toda la infraestructura de test
type TestEnv struct {
	DB       *gorm.DB
	Repos    *repositories.RepositoryContainer
	Services *services.ServiceContainer
}

// SetupTestDB crea una base de datos de test y retorna la conexión
func SetupTestDB(t *testing.T) *TestEnv {
	t.Helper()

	dbHost := getEnv("TEST_DB_HOST", "localhost")
	dbPort := getEnv("TEST_DB_PORT", "3306")
	dbUser := getEnv("TEST_DB_USER", "root")
	dbPass := getEnv("TEST_DB_PASSWORD", "")
	dbName := fmt.Sprintf("sgp_test_%d", rand.Intn(999999))

	// Conectar sin DB para crearla
	rootDSN := fmt.Sprintf("%s:%s@tcp(%s:%s)/?charset=utf8mb4&parseTime=True&loc=Local",
		dbUser, dbPass, dbHost, dbPort)

	rootDB, err := gorm.Open(mysql.Open(rootDSN), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		t.Skipf("No se pudo conectar a MySQL: %v (skipping test)", err)
	}

	rootDB.Exec(fmt.Sprintf("DROP DATABASE IF EXISTS %s", dbName))
	rootDB.Exec(fmt.Sprintf("CREATE DATABASE %s CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci", dbName))

	// Conectar a la DB de test
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		dbUser, dbPass, dbHost, dbPort, dbName)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger:                                   logger.Default.LogMode(logger.Silent),
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		t.Fatalf("Error conectando a DB de test: %v", err)
	}

	// Migrar todas las tablas
	if err := db.AutoMigrate(models.AllModels()...); err != nil {
		t.Fatalf("Error en migraciones de test: %v", err)
	}

	repos := repositories.NewRepositoryContainer(db)
	svc := services.NewServiceContainer(db, repos, services.ServiceConfig{
		JWTSecret:       "test-secret-key",
		JWTExpiration:   1 * time.Hour,
		RegisterEnabled: true,
	})

	// Registrar cleanup
	t.Cleanup(func() {
		sqlDB, _ := db.DB()
		if sqlDB != nil {
			sqlDB.Close()
		}
		// Limpiar: conectar como root y dropear
		cleanDB, err := gorm.Open(mysql.Open(rootDSN), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Silent),
		})
		if err == nil {
			cleanDB.Exec(fmt.Sprintf("DROP DATABASE IF EXISTS %s", dbName))
			cdb, _ := cleanDB.DB()
			if cdb != nil {
				cdb.Close()
			}
		}
	})

	return &TestEnv{DB: db, Repos: repos, Services: svc}
}

// SeedAdmin crea un usuario admin de test
func SeedAdmin(t *testing.T, env *TestEnv) uint {
	t.Helper()
	nombre := fmt.Sprintf("admin_%d", time.Now().UnixNano())
	user := &models.Usuario{
		Username:     nombre,
		Email:        fmt.Sprintf("%s@test.com", nombre),
		PasswordHash: "hash",
		Rol:          models.RolAdmin,
		Activo:       true,
	}
	if err := env.Repos.Usuario.Create(user); err != nil {
		t.Fatalf("Error creando admin seed: %v", err)
	}
	return user.ID
}

// SeedPropietario crea un usuario propietario de test
func SeedPropietario(t *testing.T, env *TestEnv) uint {
	t.Helper()
	nombre := fmt.Sprintf("prop_%d", time.Now().UnixNano())
	user := &models.Usuario{
		Username:     nombre,
		Email:        fmt.Sprintf("%s@test.com", nombre),
		PasswordHash: "hash",
		Rol:          models.RolPropietario,
		Activo:       true,
	}
	if err := env.Repos.Usuario.Create(user); err != nil {
		t.Fatalf("Error creando propietario seed: %v", err)
	}
	return user.ID
}

// SeedEmpleado crea un empleado vinculado a un propietario
func SeedEmpleado(t *testing.T, env *TestEnv, propietarioID uint) uint {
	t.Helper()
	nombre := fmt.Sprintf("emp_%d", time.Now().UnixNano())
	user := &models.Usuario{
		Username:      nombre,
		Email:         fmt.Sprintf("%s@test.com", nombre),
		PasswordHash:  "hash",
		Rol:           models.RolEmpleado,
		PropietarioID: &propietarioID,
		Activo:        true,
	}
	if err := env.Repos.Usuario.Create(user); err != nil {
		t.Fatalf("Error creando empleado seed: %v", err)
	}
	return user.ID
}

// SeedGranja crea una granja de test (nuevo propietario) y retorna su ID
func SeedGranja(t *testing.T, env *TestEnv) uint {
	t.Helper()
	return SeedGranjaFor(t, env, SeedPropietario(t, env))
}

// SeedGranjaFor crea una granja para un propietario existente
func SeedGranjaFor(t *testing.T, env *TestEnv, propietarioID uint) uint {
	t.Helper()
	granja, err := env.Services.Granja.Crear(services.CrearGranjaInput{
		Nombre:      "Granja Test",
		Descripcion: "Granja para tests",
		Ubicacion:   "Test City",
	}, propietarioID)
	if err != nil {
		t.Fatalf("Error creando granja seed: %v", err)
	}
	return granja.ID
}

// SeedActor carga un Actor desde la BD por ID de usuario
func SeedActor(t *testing.T, env *TestEnv, userID uint) authz.Actor {
	t.Helper()
	user, err := env.Repos.Usuario.FindByID(userID)
	if err != nil {
		t.Fatalf("Error cargando usuario seed %d: %v", userID, err)
	}
	return authz.ActorFromUsuario(user)
}

// SeedCorral crea un corral de test
func SeedCorral(t *testing.T, svc *services.ServiceContainer, granjaID uint) uint {
	t.Helper()
	corral, err := svc.Corral.Crear(services.CrearCorralInput{
		GranjaID: granjaID,
		Nombre:   "Corral Test",
	})
	if err != nil {
		t.Fatalf("Error creando corral seed: %v", err)
	}
	return corral.ID
}

// SeedPadrillo crea un padrillo de test
func SeedPadrillo(t *testing.T, svc *services.ServiceContainer, granjaID uint) uint {
	t.Helper()
	padrillo, err := svc.Padrillo.Crear(services.CrearPadrilloInput{
		GranjaID:       granjaID,
		NumeroCaravana: fmt.Sprintf("P-%d", time.Now().UnixNano()),
		Nombre:         "Padrillo Test",
	})
	if err != nil {
		t.Fatalf("Error creando padrillo seed: %v", err)
	}
	return padrillo.ID
}

// SeedCerda crea una cerda de test en estado disponible
func SeedCerda(t *testing.T, svc *services.ServiceContainer, granjaID uint) uint {
	t.Helper()
	cerda, err := svc.Cerda.Crear(services.CrearCerdaInput{
		GranjaID:       granjaID,
		NumeroCaravana: fmt.Sprintf("C-%d", time.Now().UnixNano()),
		Estado:         "disponible",
	})
	if err != nil {
		t.Fatalf("Error creando cerda seed: %v", err)
	}
	return cerda.ID
}

func getEnv(key, defaultVal string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return defaultVal
}

func init() {
	log.SetOutput(os.Stderr)
}
