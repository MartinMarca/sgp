package services

import "errors"

// Errores de negocio comunes
var (
	// Generales
	ErrNotFound     = errors.New("registro no encontrado")
	ErrDuplicateKey = errors.New("registro duplicado")
	ErrForbidden    = errors.New("operación no permitida")

	// Cerdas
	ErrCerdaNoDisponible     = errors.New("la cerda no está en estado disponible")
	ErrCerdaNoEnServicio     = errors.New("la cerda no está en estado servicio")
	ErrCerdaNoEnGestacion    = errors.New("la cerda no está en estado gestación")
	ErrCerdaNoEnCria         = errors.New("la cerda no está en estado cría")
	ErrCerdaNoActiva         = errors.New("la cerda no está activa")
	ErrCerdaTieneServicioActivo = errors.New("la cerda tiene un servicio activo, cancele la preñez primero")
	ErrCaravanaDuplicada     = errors.New("ya existe una cerda/padrillo con ese número de caravana en la granja")

	// Servicios
	ErrServicioRequierePadrillo = errors.New("la monta natural requiere un padrillo")
	ErrServicioRequierePajuela  = errors.New("la inseminación requiere un número de pajuela")

	// Preñez
	ErrPrenezYaConfirmada   = errors.New("la preñez ya fue confirmada")
	ErrPrenezYaCancelada    = errors.New("la preñez ya fue cancelada")
	ErrPrenezNoConfirmada   = errors.New("la preñez no está confirmada")
	ErrMotivoRequerido      = errors.New("el motivo de cancelación es requerido")

	// Partos
	ErrNoHayServicioConfirmado = errors.New("no hay servicio con preñez confirmada para esta cerda")
	ErrLechonesInvalidos       = errors.New("la suma de hembras + machos debe ser igual a los nacidos vivos")
	ErrTotalMenorQueVivos      = errors.New("los lechones totales no pueden ser menor a los nacidos vivos")

	// Destetes
	ErrNoHayPartoSinDestete   = errors.New("no hay parto sin destete para esta cerda")
	ErrDestetadosExcedenVivos = errors.New("la cantidad de destetados no puede superar los nacidos vivos del parto")
	ErrLoteRequerido          = errors.New("se debe asignar un lote para los lechones destetados")

	// Lotes
	ErrLoteNoActivo    = errors.New("el lote no está en estado activo")
	ErrCorralRequerido = errors.New("el lote debe estar asignado a un corral")

	// Corrales
	ErrCorralTieneLotesActivos = errors.New("el corral tiene lotes activos, reasígnelos antes de dar de baja")

	// Granjas
	ErrGranjaTieneDatosActivos = errors.New("la granja tiene cerdas, padrillos o corrales activos")

	// Ventas
	ErrVentaExcedeLote    = errors.New("la cantidad a vender supera los lechones disponibles en el lote")
	ErrTipoAnimalInvalido = errors.New("el tipo de animal no es válido")

	// Muertes de lechones
	ErrMuerteRequierePartoOLote  = errors.New("la muerte debe estar asociada a un parto (lactancia) o a un lote (engorde), no ambos")
	ErrMuertesExcedenVivos       = errors.New("las muertes registradas superan los lechones nacidos vivos del parto")
	ErrMuertesExcedenLote        = errors.New("las muertes registradas superan los lechones disponibles en el lote")
	ErrCausaMuerteInvalida       = errors.New("la causa de muerte no es válida")
)
