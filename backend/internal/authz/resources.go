package authz

import "github.com/martin/sgp/internal/models"

// RequireOnCerda verifica permiso sobre la granja de la cerda.
func (s *Service) RequireOnCerda(actor Actor, cerdaID uint, perm string) error {
	cerda, err := s.repos.Cerda.FindByID(cerdaID)
	if err != nil {
		return ErrNotFound
	}
	return s.RequireGranja(actor, cerda.GranjaID, perm)
}

// RequireOnPadrillo verifica permiso sobre la granja del padrillo.
func (s *Service) RequireOnPadrillo(actor Actor, padrilloID uint, perm string) error {
	padrillo, err := s.repos.Padrillo.FindByID(padrilloID)
	if err != nil {
		return ErrNotFound
	}
	return s.RequireGranja(actor, padrillo.GranjaID, perm)
}

// RequireOnCorral verifica permiso sobre la granja del corral.
func (s *Service) RequireOnCorral(actor Actor, corralID uint, perm string) error {
	corral, err := s.repos.Corral.FindByID(corralID)
	if err != nil {
		return ErrNotFound
	}
	return s.RequireGranja(actor, corral.GranjaID, perm)
}

// RequireOnLote verifica permiso sobre la granja del lote (vía corral).
func (s *Service) RequireOnLote(actor Actor, loteID uint, perm string) error {
	lote, err := s.repos.Lote.FindByID(loteID)
	if err != nil {
		return ErrNotFound
	}
	return s.RequireOnCorral(actor, lote.CorralID, perm)
}

// RequireOnVenta verifica permiso sobre la granja de la venta.
func (s *Service) RequireOnVenta(actor Actor, ventaID uint, perm string) error {
	venta, err := s.repos.Venta.FindByID(ventaID)
	if err != nil {
		return ErrNotFound
	}
	return s.RequireGranja(actor, venta.GranjaID, perm)
}

// RequireOnServicio verifica permiso sobre la granja del servicio (vía cerda).
func (s *Service) RequireOnServicio(actor Actor, servicioID uint, perm string) error {
	servicio, err := s.repos.Servicio.FindByID(servicioID)
	if err != nil {
		return ErrNotFound
	}
	return s.RequireOnCerda(actor, servicio.CerdaID, perm)
}

// RequireOnParto verifica permiso sobre la granja del parto (vía cerda).
func (s *Service) RequireOnParto(actor Actor, partoID uint, perm string) error {
	parto, err := s.repos.Parto.FindByID(partoID)
	if err != nil {
		return ErrNotFound
	}
	return s.RequireOnCerda(actor, parto.CerdaID, perm)
}

// RequireOnDestete verifica permiso sobre la granja del destete (vía cerda).
func (s *Service) RequireOnDestete(actor Actor, desteteID uint, perm string) error {
	destete, err := s.repos.Destete.FindByID(desteteID)
	if err != nil {
		return ErrNotFound
	}
	return s.RequireOnCerda(actor, destete.CerdaID, perm)
}

// RequireOnMuerte verifica permiso sobre la granja del registro de muerte.
func (s *Service) RequireOnMuerte(actor Actor, muerteID uint, perm string) error {
	muerte, err := s.repos.MuerteLechon.FindByID(muerteID)
	if err != nil {
		return ErrNotFound
	}
	return s.RequireGranja(actor, muerte.GranjaID, perm)
}

// RequireGranjaQuery valida acceso cuando hay granja_id en query (obligatorio para no-admin).
func (s *Service) RequireGranjaQuery(actor Actor, granjaID *uint, perm string) error {
	if actor.IsAdmin() {
		if granjaID != nil && *granjaID > 0 {
			return s.RequireGranja(actor, *granjaID, perm)
		}
		return nil
	}
	if granjaID == nil || *granjaID == 0 {
		return ErrForbidden
	}
	return s.RequireGranja(actor, *granjaID, perm)
}

// RequireAdminOnly restringe una acción solo a administradores.
func (s *Service) RequireAdminOnly(actor Actor) error {
	if !actor.IsAdmin() {
		return ErrForbidden
	}
	return nil
}

// GranjaIDFromCerda obtiene el ID de granja de una cerda.
func (s *Service) GranjaIDFromCerda(cerdaID uint) (uint, error) {
	cerda, err := s.repos.Cerda.FindByID(cerdaID)
	if err != nil {
		return 0, ErrNotFound
	}
	return cerda.GranjaID, nil
}

// GranjaIDFromCorral obtiene el ID de granja de un corral.
func (s *Service) GranjaIDFromCorral(corralID uint) (uint, error) {
	corral, err := s.repos.Corral.FindByID(corralID)
	if err != nil {
		return 0, ErrNotFound
	}
	return corral.GranjaID, nil
}

// ValidateGranjaIDForCreate valida acceso al crear recursos con granja_id explícito.
func (s *Service) ValidateGranjaIDForCreate(actor Actor, granjaID uint, perm string) error {
	return s.RequireGranja(actor, granjaID, perm)
}

// ValidateCerdaIDForCreate valida acceso al crear recursos referenciando una cerda.
func (s *Service) ValidateCerdaIDForCreate(actor Actor, cerdaID uint, perm string) error {
	return s.RequireOnCerda(actor, cerdaID, perm)
}

// Role constants re-export for middleware convenience.
const (
	RoleEmpleado = models.RolEmpleado
)
