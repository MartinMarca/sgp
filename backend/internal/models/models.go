package models

// Este archivo facilita la importación de todos los modelos desde un solo lugar
// y define constantes útiles para los enums

// Estados de Cerda
const (
	EstadoCerdaDisponible = "disponible"
	EstadoCerdaServicio   = "servicio"
	EstadoCerdaGestacion  = "gestacion"
	EstadoCerdaCria       = "cria"
)

// Estados de Lote
const (
	EstadoLoteActivo  = "activo"
	EstadoLoteCerrado = "cerrado"
	EstadoLoteVendido = "vendido"
)

// Tipos de Monta
const (
	TipoMontaNatural      = "natural"
	TipoMontaInseminacion = "inseminacion"
)

// Motivos de Baja
const (
	MotivoBajaMuerte = "muerte"
	MotivoBajaVenta  = "venta"
)

// Roles de Usuario (sistema)
const (
	RolAdmin       = "admin"
	RolUsuario     = "usuario"
	RolVeterinario = "veterinario"
)

// Causas de Muerte de Lechones
const (
	CausaMuerteAplastamiento = "aplastamiento"
	CausaMuerteEnfermedad    = "enfermedad"
	CausaMuerteInanicion     = "inanicion"
	CausaMuerteOtro          = "otro"
)

// Roles en Granja (usuario-granja)
const (
	RolGranjaPropietario   = "propietario"
	RolGranjaAdministrador = "administrador"
	RolGranjaOperador      = "operador"
)

// AllModels retorna un slice con todos los modelos para migraciones
func AllModels() []interface{} {
	return []interface{}{
		&Usuario{},
		&Granja{},
		&UsuarioGranja{},
		&Corral{},
		&Lote{},
		&Padrillo{},
		&Cerda{},
		&Servicio{},
		&Parto{},
		&Destete{},
		&MuerteLechon{},
		&Venta{},
	}
}
