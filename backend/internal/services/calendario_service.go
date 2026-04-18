package services

import (
	"time"

	"github.com/martin/sgp/internal/models"
	"github.com/martin/sgp/internal/repositories"
)

// CalendarioService maneja la lógica de eventos futuros (calendario)
type CalendarioService struct {
	repos *repositories.RepositoryContainer
}

// NewCalendarioService crea una nueva instancia del servicio
func NewCalendarioService(repos *repositories.RepositoryContainer) *CalendarioService {
	return &CalendarioService{repos: repos}
}

// --- DTOs ---

// EventoCalendario representa un evento futuro en el calendario
type EventoCalendario struct {
	ID              uint      `json:"id"`
	Tipo            string    `json:"tipo"`             // "parto_estimado", "destete_estimado", "confirmacion_pendiente"
	FechaEstimada   time.Time `json:"fecha_estimada"`
	CerdaID         uint      `json:"cerda_id"`
	CerdaCaravana   string    `json:"cerda_caravana"`
	GranjaID        uint      `json:"granja_id,omitempty"`
	GranjaNombre    string    `json:"granja_nombre,omitempty"`
	Descripcion     string    `json:"descripcion"`
	DiasRestantes   int       `json:"dias_restantes"`
}

// --- Métodos ---

// GetEventosFuturos obtiene todos los eventos futuros combinados
func (s *CalendarioService) GetEventosFuturos(granjaID *uint, diasAntes int) ([]EventoCalendario, error) {
	var eventos []EventoCalendario

	// 1. Partos futuros (cerdas en gestación con fecha estimada de parto)
	partosFuturos, err := s.getPartosEstimados(granjaID, diasAntes)
	if err != nil {
		return nil, err
	}
	eventos = append(eventos, partosFuturos...)

	// 2. Destetes futuros (cerdas en cría con fecha estimada de destete)
	destetesFuturos, err := s.getDestetesEstimados(granjaID, diasAntes)
	if err != nil {
		return nil, err
	}
	eventos = append(eventos, destetesFuturos...)

	// 3. Servicios pendientes de confirmación
	pendientes, err := s.getConfirmacionesPendientes(granjaID)
	if err != nil {
		return nil, err
	}
	eventos = append(eventos, pendientes...)

	return eventos, nil
}

// getPartosEstimados obtiene los partos estimados próximos
// Busca cerdas en gestación y calcula la fecha estimada basándose en fecha_servicio + 114 días
func (s *CalendarioService) getPartosEstimados(granjaID *uint, diasAntes int) ([]EventoCalendario, error) {
	var eventos []EventoCalendario

	// Obtener cerdas en gestación
	estado := models.EstadoCerdaGestacion
	cerdas, err := s.repos.Cerda.FindByEstado(estado, granjaID)
	if err != nil {
		return nil, err
	}

	ahora := time.Now()
	limiteInferior := ahora.AddDate(0, 0, -diasAntes)
	limiteSuperior := ahora.AddDate(0, 0, diasAntes)

	for _, cerda := range cerdas {
		// Obtener el servicio con preñez confirmada
		servicio, err := s.repos.Servicio.GetServicioConPrenezConfirmada(cerda.ID)
		if err != nil {
			continue
		}

		fechaEstimadaParto := servicio.FechaServicio.AddDate(0, 0, DiasGestacion)

		// Filtrar por rango de fechas
		if fechaEstimadaParto.Before(limiteInferior) || fechaEstimadaParto.After(limiteSuperior) {
			continue
		}

		diasRestantes := int(fechaEstimadaParto.Sub(ahora).Hours() / 24)

		eventos = append(eventos, EventoCalendario{
			ID:            cerda.ID,
			Tipo:          "parto_estimado",
			FechaEstimada: fechaEstimadaParto,
			CerdaID:       cerda.ID,
			CerdaCaravana: cerda.NumeroCaravana,
			GranjaID:      cerda.GranjaID,
			Descripcion:   "Parto estimado",
			DiasRestantes: diasRestantes,
		})
	}

	return eventos, nil
}

// getDestetesEstimados obtiene los destetes estimados próximos
// Busca cerdas en cría y calcula la fecha estimada basándose en fecha_parto + 30 días
func (s *CalendarioService) getDestetesEstimados(granjaID *uint, diasAntes int) ([]EventoCalendario, error) {
	var eventos []EventoCalendario

	estado := models.EstadoCerdaCria
	cerdas, err := s.repos.Cerda.FindByEstado(estado, granjaID)
	if err != nil {
		return nil, err
	}

	ahora := time.Now()
	limiteInferior := ahora.AddDate(0, 0, -diasAntes)
	limiteSuperior := ahora.AddDate(0, 0, diasAntes)

	for _, cerda := range cerdas {
		// Obtener el último parto de esta cerda
		parto, err := s.repos.Cerda.GetUltimoParto(cerda.ID)
		if err != nil {
			continue
		}

		fechaEstimadaDestete := parto.FechaParto.AddDate(0, 0, DiasCria)

		if fechaEstimadaDestete.Before(limiteInferior) || fechaEstimadaDestete.After(limiteSuperior) {
			continue
		}

		diasRestantes := int(fechaEstimadaDestete.Sub(ahora).Hours() / 24)

		eventos = append(eventos, EventoCalendario{
			ID:            cerda.ID,
			Tipo:          "destete_estimado",
			FechaEstimada: fechaEstimadaDestete,
			CerdaID:       cerda.ID,
			CerdaCaravana: cerda.NumeroCaravana,
			GranjaID:      cerda.GranjaID,
			Descripcion:   "Destete estimado",
			DiasRestantes: diasRestantes,
		})
	}

	return eventos, nil
}

// getConfirmacionesPendientes obtiene servicios pendientes de confirmar preñez
func (s *CalendarioService) getConfirmacionesPendientes(granjaID *uint) ([]EventoCalendario, error) {
	var eventos []EventoCalendario

	servicios, err := s.repos.Servicio.GetServiciosPendientesConfirmacion(granjaID)
	if err != nil {
		return nil, err
	}

	ahora := time.Now()

	for _, servicio := range servicios {
		diasDesdeServicio := int(ahora.Sub(servicio.FechaServicio).Hours() / 24)

		eventos = append(eventos, EventoCalendario{
			ID:            servicio.ID,
			Tipo:          "confirmacion_pendiente",
			FechaEstimada: servicio.FechaServicio,
			CerdaID:       servicio.CerdaID,
			CerdaCaravana: servicio.Cerda.NumeroCaravana,
			Descripcion:   "Confirmación de preñez pendiente",
			DiasRestantes: -diasDesdeServicio, // Negativo porque ya pasó
		})
	}

	return eventos, nil
}
