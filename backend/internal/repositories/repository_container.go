package repositories

import "gorm.io/gorm"

// RepositoryContainer contiene todos los repositories de la aplicación
type RepositoryContainer struct {
	Usuario       *UsuarioRepository
	Granja        *GranjaRepository
	Corral        *CorralRepository
	Lote          *LoteRepository
	Padrillo      *PadrilloRepository
	Cerda         *CerdaRepository
	Servicio      *ServicioRepository
	Parto         *PartoRepository
	Destete       *DesteteRepository
	MuerteLechon  *MuerteLechonRepository
	Venta         *VentaRepository
}

// NewRepositoryContainer crea una nueva instancia de RepositoryContainer
func NewRepositoryContainer(db *gorm.DB) *RepositoryContainer {
	return &RepositoryContainer{
		Usuario:       NewUsuarioRepository(db),
		Granja:        NewGranjaRepository(db),
		Corral:        NewCorralRepository(db),
		Lote:          NewLoteRepository(db),
		Padrillo:      NewPadrilloRepository(db),
		Cerda:         NewCerdaRepository(db),
		Servicio:      NewServicioRepository(db),
		Parto:         NewPartoRepository(db),
		Destete:       NewDesteteRepository(db),
		MuerteLechon:  NewMuerteLechonRepository(db),
		Venta:         NewVentaRepository(db),
	}
}
