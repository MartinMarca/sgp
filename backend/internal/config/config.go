package config

import "os"

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
}

// Load carga la configuración desde las variables de entorno
func Load() *Config {
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
		Env:      getEnv("ENV", "development"),
		LogLevel: getEnv("LOG_LEVEL", "debug"),
	}
}

// getEnv obtiene una variable de entorno o devuelve un valor por defecto
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
