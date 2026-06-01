package config

import (
	"os"
	"strconv"
	"strings"
)

// Config contiene la configuración de la aplicación
type Config struct {
	// Database
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string

	// Server
	ServerHost string
	ServerPort string

	// JWT
	JWTSecret     string
	JWTExpiration string

	// CORS
	CORSOrigin string

	// Environment
	Env      string
	LogLevel string

	// Auth
	RegisterEnabled     bool
	BootstrapAdminEmail string
}

// Load carga la configuración desde las variables de entorno
func Load() *Config {
	env := getEnv("ENV", "development")
	registerEnabled := true
	if v := os.Getenv("REGISTER_ENABLED"); v != "" {
		registerEnabled, _ = strconv.ParseBool(v)
	} else if env == "production" {
		registerEnabled = false
	}

	return &Config{
		// Database
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "3306"),
		DBUser:     getEnv("DB_USER", "root"),
		DBPassword: getEnv("DB_PASSWORD", ""),
		DBName:     getEnv("DB_NAME", "granja_porcina"),

		// Server
		ServerHost: getEnv("SERVER_HOST", "localhost"),
		ServerPort: getEnv("SERVER_PORT", "8080"),

		// JWT
		JWTSecret:     getEnv("JWT_SECRET", "tu_clave_secreta_super_segura"),
		JWTExpiration: getEnv("JWT_EXPIRATION", "24h"),

		// CORS
		CORSOrigin: getEnv("CORS_ORIGIN", "http://localhost:8080"),

		// Environment
		Env:      env,
		LogLevel: getEnv("LOG_LEVEL", "debug"),

		RegisterEnabled:     registerEnabled,
		BootstrapAdminEmail: strings.TrimSpace(getEnv("BOOTSTRAP_ADMIN_EMAIL", "")),
	}
}

// getEnv obtiene una variable de entorno o devuelve un valor por defecto
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
