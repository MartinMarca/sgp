package database

import (
	"fmt"
	"log"

	"github.com/martin/sgp/internal/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Connect establece la conexión con la base de datos MySQL
func Connect(cfg *config.Config) (*gorm.DB, error) {
	// Construir DSN (Data Source Name)
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=UTC",
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBName,
	)

	// Configurar logger según el entorno
	var logLevel logger.LogLevel
	if cfg.Env == "production" {
		logLevel = logger.Silent
	} else {
		logLevel = logger.Info
	}

	// Conectar a la base de datos
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logLevel),
	})

	if err != nil {
		return nil, fmt.Errorf("error al conectar con MySQL: %w", err)
	}

	// Obtener la conexión SQL subyacente para configurar el pool
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("error al obtener conexión SQL: %w", err)
	}

	// Configurar pool de conexiones
	sqlDB.SetMaxOpenConns(25)
	sqlDB.SetMaxIdleConns(5)

	log.Println("✅ Conexión a MySQL establecida correctamente")
	return db, nil
}
