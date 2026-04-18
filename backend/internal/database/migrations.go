package database

import (
	"log"

	"github.com/martin/sgp/internal/models"
	"gorm.io/gorm"
)

// AutoMigrate ejecuta las migraciones automáticas de GORM
// Nota: Esto es útil para desarrollo, en producción se recomienda usar migraciones SQL manuales
func AutoMigrate(db *gorm.DB) error {
	log.Println("Ejecutando migraciones automáticas...")

	// Usar models.AllModels() para obtener todos los modelos
	err := db.AutoMigrate(models.AllModels()...)

	if err != nil {
		return err
	}

	log.Println("✅ Migraciones ejecutadas correctamente")
	return nil
}
