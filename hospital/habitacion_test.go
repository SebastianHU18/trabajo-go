package hospital

import (
	"errors"
	"strings"
	"testing"
	"time"
)

func TestHabitacionCapacidadYLiberacion(t *testing.T) {
	h := NuevaHabitacion(101, 2)
	a := NuevoPaciente("P-1", "Ana", 25, Leve)
	b := NuevoPaciente("P-2", "Beto", 26, Moderado)
	c := NuevoPaciente("P-3", "Caro", 27, Severo)
	if !h.EstaDisponible() || h.Estado() != Disponible {
		t.Fatal("la habitacion nueva debe estar disponible")
	}
	if err := h.Ocupar(a); err != nil {
		t.Fatal(err)
	}
	if h.Estado() != Ocupada || !h.EstaDisponible() {
		t.Fatal("una cama ocupada no elimina el otro cupo")
	}
	if err := h.Ocupar(b); err != nil {
		t.Fatal(err)
	}
	if h.EstaDisponible() {
		t.Fatal("habitacion llena figura disponible")
	}
	if err := h.Ocupar(c); !errors.Is(err, ErrHabitacionOcupada) {
		t.Fatalf("esperaba ErrHabitacionOcupada: %v", err)
	}
	if err := h.Ocupar(a); !errors.Is(err, ErrHabitacionOcupada) {
		t.Fatal("paciente duplicado debe fallar")
	}
	lista := h.Ocupantes()
	lista[0] = c
	if h.Ocupantes()[0] != a {
		t.Fatal("Ocupantes expone el slice interno")
	}
	if err := h.Liberar(c); !errors.Is(err, ErrPacienteNoEnHabitacion) {
		t.Fatal("liberar ausente debe fallar")
	}
	if err := h.Liberar(a); err != nil {
		t.Fatal(err)
	}
	if h.Estado() != Ocupada || !h.EstaDisponible() {
		t.Fatal("queda un ocupante con cupo libre")
	}
	if err := h.Liberar(b); err != nil {
		t.Fatal(err)
	}
	if h.Estado() != Disponible || len(h.Ocupantes()) != 0 {
		t.Fatal("no se libero completamente")
	}
}

func TestHabitacionCapacidadInvalida(t *testing.T) {
	h := NuevaHabitacion(99, 0)
	if h.EstaDisponible() {
		t.Fatal("no puede haber cupos con capacidad 0")
	}
	if !errors.Is(h.Ocupar(NuevoPaciente("P-1", "Ana", 23, Leve)), ErrHabitacionOcupada) {
		t.Fatal("debe rechazar capacidad invalida")
	}
}

func TestPersonalPolimorficoYEpisodios(t *testing.T) {
	medico := NuevoMedico("D-1", "Karen", 35, "Neurologia")
	camillero := NuevoCamillero("C-1", "Wilson", 29, "manana")
	personal := []Atendedor{medico, camillero}
	p := NuevoPaciente("P-1", "Ana", 29, Severo)
	p.SufrirAtaqueDeSueno("cafeteria")
	for _, a := range personal {
		e, err := a.Atender(p, "cafeteria")
		if err != nil {
			t.Fatalf("no atendio %T: %v", a, err)
		}
		if e.Paciente() != p || e.AtendidoPor() != a || e.Habitacion() != nil {
			t.Fatal("registro incorrecto")
		}
		if e.Fecha().IsZero() || !strings.Contains(e.Resumen(), "sin habitacion") {
			t.Fatal("resumen invalido")
		}
	}
	if len(medico.MisEpisodios()) != 1 {
		t.Fatal("medico no guardo su episodio")
	}
	if err := medico.DiagnosticarPaciente(p); err != nil {
		t.Fatal(err)
	}
	if err := medico.DiagnosticarPaciente(p); err != nil {
		t.Fatal(err)
	}
	if len(medico.Pacientes()) != 1 {
		t.Fatal("paciente duplicado")
	}
}

func TestNuevoRegistroEpisodioValida(t *testing.T) {
	p := NuevoPaciente("P-1", "Ana", 29, Severo)
	d := NuevoMedico("D-1", "Karen", 35, "Neurologia")
	if _, err := NuevoRegistroEpisodio("", time.Now(), p, d, "cafeteria", nil); err == nil {
		t.Fatal("id vacio debe fallar")
	}
	reg, err := NuevoRegistroEpisodio("E-1", time.Date(2026, 10, 19, 12, 0, 0, 0, time.UTC), p, d, "cafeteria", nil)
	if err != nil || reg.ID() != "E-1" {
		t.Fatalf("registro invalido: %v", err)
	}
}
