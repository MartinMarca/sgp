package bootstrap

import (
	"log"

	"github.com/martin/sgp/internal/config"
	"github.com/martin/sgp/internal/models"
	"github.com/martin/sgp/internal/repositories"
)

// Admin promueve un usuario a admin según configuración o crea el primero admin si no hay ninguno.
func Admin(cfg *config.Config, repos *repositories.RepositoryContainer) {
	if cfg.BootstrapAdminEmail != "" {
		user, err := repos.Usuario.FindByEmail(cfg.BootstrapAdminEmail)
		if err != nil {
			log.Printf("Bootstrap admin: usuario con email %s no encontrado", cfg.BootstrapAdminEmail)
			return
		}
		if user.Rol != models.RolAdmin {
			user.Rol = models.RolAdmin
			user.PropietarioID = nil
			if err := repos.Usuario.Update(user); err != nil {
				log.Printf("Bootstrap admin: error actualizando usuario: %v", err)
			} else {
				log.Printf("Bootstrap admin: %s promovido a admin", cfg.BootstrapAdminEmail)
			}
		}
		return
	}

	users, err := repos.Usuario.FindAll(nil)
	if err != nil || len(users) == 0 {
		return
	}
	for _, u := range users {
		if u.Rol == models.RolAdmin {
			return
		}
	}
	first := users[0]
	first.Rol = models.RolAdmin
	first.PropietarioID = nil
	if err := repos.Usuario.Update(&first); err != nil {
		log.Printf("Bootstrap admin: error promoviendo primer usuario: %v", err)
	} else {
		log.Printf("Bootstrap admin: usuario id=%d promovido a admin (no había admin)", first.ID)
	}
}
