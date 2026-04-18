package main

import (
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/martin/sgp/internal/config"
	"github.com/martin/sgp/internal/database"
	"github.com/martin/sgp/internal/repositories"
	"github.com/martin/sgp/internal/routes"
	"github.com/martin/sgp/internal/services"
)

func main() {
	// Cargar variables de entorno
	if err := godotenv.Load(); err != nil {
		log.Println("No se encontró archivo .env, usando variables del sistema")
	}

	// Cargar configuración
	cfg := config.Load()

	// Conectar a la base de datos
	db, err := database.Connect(cfg)
	if err != nil {
		log.Fatalf("Error al conectar con la base de datos: %v", err)
	}

	// Ejecutar migraciones automáticas (opcional)
	if os.Getenv("AUTO_MIGRATE") == "true" {
		if err := database.AutoMigrate(db); err != nil {
			log.Fatalf("Error en migraciones: %v", err)
		}
	}

	// Parsear expiración JWT
	jwtExpiration, err := time.ParseDuration(cfg.JWTExpiration)
	if err != nil {
		jwtExpiration = 24 * time.Hour // Default: 24 horas
	}

	// Inicializar capas
	repos := repositories.NewRepositoryContainer(db)
	svc := services.NewServiceContainer(db, repos, services.ServiceConfig{
		JWTSecret:     cfg.JWTSecret,
		JWTExpiration: jwtExpiration,
	})

	// Configurar rutas
	router := routes.SetupRoutes(cfg, svc)

	// Iniciar servidor
	addr := cfg.ServerHost + ":" + cfg.ServerPort
	log.Printf("Servidor iniciado en http://%s", addr)
	if err := router.Run(addr); err != nil {
		log.Fatalf("Error al iniciar el servidor: %v", err)
	}
}
