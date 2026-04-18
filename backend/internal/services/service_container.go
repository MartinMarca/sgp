package services

import (
	"time"

	"github.com/martin/sgp/internal/repositories"
	"gorm.io/gorm"
)

// ServiceContainer contiene todos los services de la aplicación
type ServiceContainer struct {
	Auth           *AuthService
	Granja         *GranjaService
	Corral         *CorralService
	Lote           *LoteService
	Padrillo       *PadrilloService
	Cerda          *CerdaService
	Servicio       *ServicioService
	Parto          *PartoService
	Destete        *DesteteService
	MuerteLechon   *MuerteLechonService
	Venta          *VentaService
	Calendario     *CalendarioService
	Estadisticas   *EstadisticasService
}

// ServiceConfig configuración para inicializar los services
type ServiceConfig struct {
	JWTSecret     string
	JWTExpiration time.Duration
}

// NewServiceContainer crea una nueva instancia con todos los services
func NewServiceContainer(db *gorm.DB, repos *repositories.RepositoryContainer, cfg ServiceConfig) *ServiceContainer {
	return &ServiceContainer{
		Auth:         NewAuthService(repos, cfg.JWTSecret, cfg.JWTExpiration),
		Granja:       NewGranjaService(repos),
		Corral:       NewCorralService(repos),
		Lote:         NewLoteService(repos),
		Padrillo:     NewPadrilloService(repos),
		Cerda:        NewCerdaService(repos),
		Servicio:     NewServicioService(db, repos),
		Parto:        NewPartoService(db, repos),
		Destete:      NewDesteteService(db, repos),
		MuerteLechon: NewMuerteLechonService(db, repos),
		Venta:        NewVentaService(db, repos),
		Calendario:   NewCalendarioService(repos),
		Estadisticas: NewEstadisticasService(repos),
	}
}
