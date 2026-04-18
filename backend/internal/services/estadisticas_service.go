package services

import (
	"github.com/martin/sgp/internal/models"
	"github.com/martin/sgp/internal/repositories"
)

// EstadisticasService maneja la lógica de estadísticas y métricas
type EstadisticasService struct {
	repos *repositories.RepositoryContainer
}

// NewEstadisticasService crea una nueva instancia del servicio
func NewEstadisticasService(repos *repositories.RepositoryContainer) *EstadisticasService {
	return &EstadisticasService{repos: repos}
}

// --- DTOs ---

// ResumenGranja estadísticas generales de una granja
type ResumenGranja struct {
	TotalCerdas          int64              `json:"total_cerdas"`
	CerdasPorEstado      map[string]int64   `json:"cerdas_por_estado"`
	TotalPadrillos       int64              `json:"total_padrillos"`
	TotalCorrales        int64              `json:"total_corrales"`
	TotalLotesActivos    int64              `json:"total_lotes_activos"`
	TotalLechones        int                `json:"total_lechones"`
}

// EstadisticasPeriodo estadísticas para un período de tiempo
type EstadisticasPeriodo struct {
	Partos    map[string]interface{} `json:"partos"`
	Destetes  map[string]interface{} `json:"destetes"`
	Servicios map[string]interface{} `json:"servicios"`
}

// --- Métodos ---

// GetResumenGranja obtiene un resumen estadístico de una granja
func (s *EstadisticasService) GetResumenGranja(granjaID uint) (*ResumenGranja, error) {
	resumen := &ResumenGranja{
		CerdasPorEstado: make(map[string]int64),
	}

	// Total cerdas activas
	activa := true
	cerdas, err := s.repos.Cerda.FindByGranjaID(granjaID, nil, &activa)
	if err != nil {
		return nil, err
	}
	resumen.TotalCerdas = int64(len(cerdas))

	// Cerdas por estado
	for _, estado := range []string{
		models.EstadoCerdaDisponible,
		models.EstadoCerdaServicio,
		models.EstadoCerdaGestacion,
		models.EstadoCerdaCria,
	} {
		cerdasEstado, err := s.repos.Cerda.FindByEstado(estado, &granjaID)
		if err != nil {
			continue
		}
		resumen.CerdasPorEstado[estado] = int64(len(cerdasEstado))
	}

	// Total padrillos activos
	padrillos, err := s.repos.Padrillo.FindByGranjaID(granjaID, &activa)
	if err == nil {
		resumen.TotalPadrillos = int64(len(padrillos))
	}

	// Total corrales activos
	corrales, err := s.repos.Corral.FindByGranjaID(granjaID, &activa)
	if err == nil {
		resumen.TotalCorrales = int64(len(corrales))
	}

	// Lotes activos y total de lechones
	estadoActivo := models.EstadoLoteActivo
	for _, corral := range corrales {
		lotes, err := s.repos.Lote.FindByCorralID(corral.ID, &estadoActivo)
		if err != nil {
			continue
		}
		resumen.TotalLotesActivos += int64(len(lotes))
		for _, lote := range lotes {
			resumen.TotalLechones += lote.CantidadLechones
		}
	}

	return resumen, nil
}

// GetEstadisticasPeriodo obtiene estadísticas de partos y destetes para un período
func (s *EstadisticasService) GetEstadisticasPeriodo(granjaID *uint, mes, anio int) (*EstadisticasPeriodo, error) {
	stats := &EstadisticasPeriodo{}

	partoStats, err := s.repos.Parto.GetEstadisticas(granjaID, mes, anio)
	if err != nil {
		return nil, err
	}
	stats.Partos = partoStats

	desteteStats, err := s.repos.Destete.GetEstadisticas(granjaID, mes, anio)
	if err != nil {
		return nil, err
	}
	stats.Destetes = desteteStats

	servicioStats, err := s.repos.Servicio.GetEstadisticas(granjaID, mes, anio)
	if err != nil {
		return nil, err
	}
	stats.Servicios = servicioStats

	return stats, nil
}
