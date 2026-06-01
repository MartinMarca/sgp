package database

import (
	"fmt"
	"log"
	"strings"

	"gorm.io/gorm"
)

// migrateRolesLegacy prepara BD existentes antes de AutoMigrate (roles + propietario_id).
func migrateRolesLegacy(db *gorm.DB) error {
	if !tableExists(db, "usuarios") {
		return nil
	}

	log.Println("Verificando migración de roles (legacy)...")

	if err := migrateUsuariosLegacy(db); err != nil {
		return fmt.Errorf("migración usuarios: %w", err)
	}
	if err := migrateGranjasPropietario(db); err != nil {
		return fmt.Errorf("migración granjas.propietario_id: %w", err)
	}

	return nil
}

func migrateUsuariosLegacy(db *gorm.DB) error {
	if !columnExists(db, "usuarios", "propietario_id") {
		if err := db.Exec(`ALTER TABLE usuarios ADD COLUMN propietario_id BIGINT UNSIGNED NULL`).Error; err != nil {
			return err
		}
	} else {
		// Alinear tipo con GORM (bigint unsigned) si quedó como INT
		if err := db.Exec(`ALTER TABLE usuarios MODIFY COLUMN propietario_id BIGINT UNSIGNED NULL`).Error; err != nil {
			return err
		}
	}

	rolType := columnType(db, "usuarios", "rol")
	if rolType == "" {
		return nil
	}
	if strings.Contains(rolType, "usuario") || strings.Contains(rolType, "veterinario") {
		log.Println("Actualizando enum de usuarios.rol...")
		if err := db.Exec(`ALTER TABLE usuarios MODIFY COLUMN rol ENUM('admin','usuario','veterinario','propietario','empleado') DEFAULT 'usuario'`).Error; err != nil {
			return err
		}
		if err := db.Exec(`UPDATE usuarios SET rol = 'propietario' WHERE rol IN ('usuario', 'veterinario')`).Error; err != nil {
			return err
		}
		if err := db.Exec(`ALTER TABLE usuarios MODIFY COLUMN rol ENUM('admin','propietario','empleado') DEFAULT 'propietario'`).Error; err != nil {
			return err
		}
	}

	return nil
}

func migrateGranjasPropietario(db *gorm.DB) error {
	if !tableExists(db, "granjas") {
		return nil
	}

	if !columnExists(db, "granjas", "propietario_id") {
		if err := db.Exec(`ALTER TABLE granjas ADD COLUMN propietario_id BIGINT UNSIGNED NULL`).Error; err != nil {
			return err
		}
	} else {
		// Permite backfill aunque GORM haya dejado NOT NULL con valor 0
		if err := db.Exec(`ALTER TABLE granjas MODIFY COLUMN propietario_id BIGINT UNSIGNED NULL`).Error; err != nil {
			return err
		}
	}

	if tableExists(db, "usuario_granja") {
		if err := db.Exec(`
			UPDATE granjas g
			INNER JOIN usuario_granja ug ON ug.granja_id = g.id AND ug.rol = 'propietario'
			SET g.propietario_id = ug.usuario_id
			WHERE g.propietario_id IS NULL OR g.propietario_id = 0
		`).Error; err != nil {
			return err
		}
	}

	if err := db.Exec(`
		UPDATE granjas g
		SET g.propietario_id = (SELECT MIN(u.id) FROM usuarios u WHERE u.rol IN ('propietario', 'admin'))
		WHERE (g.propietario_id IS NULL OR g.propietario_id = 0)
		  AND EXISTS (SELECT 1 FROM usuarios u WHERE u.rol IN ('propietario', 'admin'))
	`).Error; err != nil {
		return err
	}

	if err := db.Exec(`
		UPDATE granjas g
		SET g.propietario_id = (SELECT MIN(u.id) FROM usuarios u)
		WHERE (g.propietario_id IS NULL OR g.propietario_id = 0)
		  AND EXISTS (SELECT 1 FROM usuarios u)
	`).Error; err != nil {
		return err
	}

	var invalid int64
	if err := db.Raw(`SELECT COUNT(*) FROM granjas WHERE propietario_id IS NULL OR propietario_id = 0`).Scan(&invalid).Error; err != nil {
		return err
	}
	if invalid > 0 {
		log.Printf("Advertencia: %d granja(s) sin propietario_id válido (sin usuarios en BD)", invalid)
		return nil
	}

	log.Println("propietario_id en granjas actualizado correctamente")
	return db.Exec(`ALTER TABLE granjas MODIFY COLUMN propietario_id BIGINT UNSIGNED NOT NULL`).Error
}

func tableExists(db *gorm.DB, table string) bool {
	var count int64
	db.Raw(`
		SELECT COUNT(*) FROM information_schema.tables
		WHERE table_schema = DATABASE() AND table_name = ?`, table).Scan(&count)
	return count > 0
}

func columnExists(db *gorm.DB, table, column string) bool {
	var count int64
	db.Raw(`
		SELECT COUNT(*) FROM information_schema.columns
		WHERE table_schema = DATABASE() AND table_name = ? AND column_name = ?`, table, column).Scan(&count)
	return count > 0
}

func columnType(db *gorm.DB, table, column string) string {
	var colType string
	db.Raw(`
		SELECT COLUMN_TYPE FROM information_schema.columns
		WHERE table_schema = DATABASE() AND table_name = ? AND column_name = ?`, table, column).Scan(&colType)
	return colType
}
