package models

import (
	"testing"

	"gorm.io/gorm"
)

// === Tests de validación de Parto ===

func TestParto_ValidateHembrasMachosIgualVivos(t *testing.T) {
	p := &Parto{
		LechonesNacidosVivos:   10,
		LechonesNacidosTotales: 12,
		LechonesHembras:        6,
		LechonesMachos:         4,
	}
	if err := p.validate(); err != nil {
		t.Errorf("esperaba nil, obtuvo error: %v", err)
	}
}

func TestParto_ValidateHembrasMachosDistintoVivos(t *testing.T) {
	p := &Parto{
		LechonesNacidosVivos:   10,
		LechonesNacidosTotales: 12,
		LechonesHembras:        3,
		LechonesMachos:         3, // 3+3=6 != 10
	}
	if err := p.validate(); err == nil {
		t.Error("esperaba error cuando h+m != vivos")
	}
}

func TestParto_ValidateTotalesMenorQueVivos(t *testing.T) {
	p := &Parto{
		LechonesNacidosVivos:   10,
		LechonesNacidosTotales: 8, // 8 < 10
		LechonesHembras:        6,
		LechonesMachos:         4,
	}
	if err := p.validate(); err == nil {
		t.Error("esperaba error cuando totales < vivos")
	}
}

func TestParto_ValidateTodoCero(t *testing.T) {
	p := &Parto{
		LechonesNacidosVivos:   0,
		LechonesNacidosTotales: 0,
		LechonesHembras:        0,
		LechonesMachos:         0,
	}
	if err := p.validate(); err != nil {
		t.Errorf("esperaba nil con todo cero, obtuvo: %v", err)
	}
}

func TestParto_ValidateTotalesIgualVivos(t *testing.T) {
	p := &Parto{
		LechonesNacidosVivos:   8,
		LechonesNacidosTotales: 8, // Sin muertos
		LechonesHembras:        5,
		LechonesMachos:         3,
	}
	if err := p.validate(); err != nil {
		t.Errorf("esperaba nil cuando totales == vivos, obtuvo: %v", err)
	}
}

// === Tests de validación de Servicio ===

func TestServicio_ValidateMontaNaturalConPadrillo(t *testing.T) {
	padrilloID := uint(1)
	s := &Servicio{
		TipoMonta:  TipoMontaNatural,
		PadrilloID: &padrilloID,
	}
	if err := s.validate(); err != nil {
		t.Errorf("esperaba nil, obtuvo: %v", err)
	}
}

func TestServicio_ValidateMontaNaturalSinPadrillo(t *testing.T) {
	s := &Servicio{
		TipoMonta:  TipoMontaNatural,
		PadrilloID: nil,
	}
	if err := s.validate(); err == nil {
		t.Error("esperaba error: monta natural sin padrillo")
	}
}

func TestServicio_ValidateInseminacionConPajuela(t *testing.T) {
	pajuela := "PAJ-001"
	s := &Servicio{
		TipoMonta:     TipoMontaInseminacion,
		NumeroPajuela: &pajuela,
	}
	if err := s.validate(); err != nil {
		t.Errorf("esperaba nil, obtuvo: %v", err)
	}
}

func TestServicio_ValidateInseminacionSinPajuela(t *testing.T) {
	s := &Servicio{
		TipoMonta:     TipoMontaInseminacion,
		NumeroPajuela: nil,
	}
	if err := s.validate(); err == nil {
		t.Error("esperaba error: inseminación sin pajuela")
	}
}

func TestServicio_ValidateInseminacionPajuelaVacia(t *testing.T) {
	pajuela := ""
	s := &Servicio{
		TipoMonta:     TipoMontaInseminacion,
		NumeroPajuela: &pajuela,
	}
	if err := s.validate(); err == nil {
		t.Error("esperaba error: inseminación con pajuela vacía")
	}
}

// === Tests de constantes ===

func TestConstantesEstadoCerda(t *testing.T) {
	estados := []string{EstadoCerdaDisponible, EstadoCerdaServicio, EstadoCerdaGestacion, EstadoCerdaCria}
	expected := []string{"disponible", "servicio", "gestacion", "cria"}
	for i, e := range estados {
		if e != expected[i] {
			t.Errorf("estado[%d]: esperaba %q, obtuvo %q", i, expected[i], e)
		}
	}
}

func TestConstantesEstadoLote(t *testing.T) {
	if EstadoLoteActivo != "activo" {
		t.Errorf("esperaba 'activo', obtuvo %q", EstadoLoteActivo)
	}
	if EstadoLoteCerrado != "cerrado" {
		t.Errorf("esperaba 'cerrado', obtuvo %q", EstadoLoteCerrado)
	}
}

func TestAllModelsRetornaModelos(t *testing.T) {
	all := AllModels()
	if len(all) != 11 {
		t.Errorf("esperaba 11 modelos, obtuvo %d", len(all))
	}
}

// === Test tabla names ===

func TestTableNames(t *testing.T) {
	tests := []struct {
		model    interface{ TableName() string }
		expected string
	}{
		{&Parto{}, "partos"},
		{&Servicio{}, "servicios"},
		{&Destete{}, "destetes"},
		{&Cerda{}, "cerdas"},
		{&Padrillo{}, "padrillos"},
		{&Granja{}, "granjas"},
		{&Corral{}, "corrales"},
		{&Lote{}, "lotes"},
		{&Usuario{}, "usuarios"},
		{&MuerteLechon{}, "muertes_lechones"},
	}

	for _, tt := range tests {
		if got := tt.model.TableName(); got != tt.expected {
			t.Errorf("TableName() = %q, esperaba %q", got, tt.expected)
		}
	}
}

// Hook test helper (simula la interfaz que GORM usa)
func TestPartoBeforeCreateHook(t *testing.T) {
	p := &Parto{
		LechonesNacidosVivos:   5,
		LechonesNacidosTotales: 5,
		LechonesHembras:        2,
		LechonesMachos:         3,
	}
	var tx *gorm.DB // nil es aceptable para la función validate interna
	if err := p.BeforeCreate(tx); err != nil {
		t.Errorf("BeforeCreate no debería fallar: %v", err)
	}

	p.LechonesMachos = 1 // 2+1=3 != 5
	if err := p.BeforeCreate(tx); err == nil {
		t.Error("BeforeCreate debería fallar con h+m != vivos")
	}
}
