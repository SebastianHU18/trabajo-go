package hospital

import (
	"errors"
	"testing"
)

// Estos tests viven en package hospital (no en hospital_test) a proposito: al
// estar dentro del paquete pueden armar structs con campos privados a mano,
// sin depender de que otro integrante termine su constructor.

// TestAsignarHabitacionConCupo cubre el caso feliz de la consulta 5.3: hay una
// habitacion libre, se asigna y el paciente queda en cama.
func TestAsignarHabitacionConCupo(t *testing.T) {
	h := NuevoHospital("prueba")
	_ = h.RegistrarHabitacion(NuevaHabitacion(101, 1))

	p := NuevoPaciente("P-1", "Ana", 30, Leve)
	_ = h.AdmitirPaciente(p)
	p.SufrirAtaqueDeSueno("cafeteria")

	hab, err := h.AsignarHabitacion(p)
	if err != nil {
		t.Fatalf("no debio fallar la asignacion: %v", err)
	}
	if hab == nil || hab.Numero() != 101 {
		t.Fatalf("esperaba la habitacion 101, obtuve %v", hab)
	}
	if p.Estado() != DormidoEnCama {
		t.Fatalf("el paciente debio quedar DormidoEnCama, quedo %v", p.Estado())
	}
	if hab.EstaDisponible() {
		t.Fatalf("la habitacion debio quedar ocupada")
	}
}

// TestAsignarHabitacionSinCupo cubre el caso de falla de la consulta 5.3: no
// queda cama, se devuelve ErrSinHabitacion y el paciente sigue en el pasillo.
func TestAsignarHabitacionSinCupo(t *testing.T) {
	h := NuevoHospital("prueba")
	_ = h.RegistrarHabitacion(NuevaHabitacion(101, 1))

	primero := NuevoPaciente("P-1", "Ana", 30, Leve)
	_ = h.AdmitirPaciente(primero)
	primero.SufrirAtaqueDeSueno("cafeteria")
	if _, err := h.AsignarHabitacion(primero); err != nil {
		t.Fatalf("la primera asignacion no debio fallar: %v", err)
	}

	segundo := NuevoPaciente("P-2", "Beto", 40, Severo)
	_ = h.AdmitirPaciente(segundo)
	segundo.SufrirAtaqueDeSueno("pasillo del segundo piso")

	hab, err := h.AsignarHabitacion(segundo)
	if !errors.Is(err, ErrSinHabitacion) {
		t.Fatalf("esperaba ErrSinHabitacion, obtuve %v", err)
	}
	if hab != nil {
		t.Fatalf("no debio devolver habitacion, devolvio %v", hab)
	}
	if segundo.Estado() != DormidoEnPasillo {
		t.Fatalf("el paciente sin cama debio quedar en el pasillo, quedo %v", segundo.Estado())
	}
}

// TestPacientesEnPasillo cubre la consulta 5.1: solo devuelve los que estan
// dormidos en un pasillo, no los despiertos ni los que estan en cama.
func TestPacientesEnPasillo(t *testing.T) {
	h := NuevoHospital("prueba")
	_ = h.RegistrarHabitacion(NuevaHabitacion(101, 1))

	despierto := NuevoPaciente("P-1", "Ana", 30, Leve)
	enPasillo := NuevoPaciente("P-2", "Beto", 40, Severo)
	enCama := NuevoPaciente("P-3", "Caro", 25, Moderado)
	for _, p := range []*Paciente{despierto, enPasillo, enCama} {
		_ = h.AdmitirPaciente(p)
	}

	enPasillo.SufrirAtaqueDeSueno("pasillo del segundo piso")
	enCama.SufrirAtaqueDeSueno("sala de espera")
	if _, err := h.AsignarHabitacion(enCama); err != nil {
		t.Fatalf("no debio fallar la asignacion: %v", err)
	}

	got := h.PacientesEnPasillo()
	if len(got) != 1 {
		t.Fatalf("esperaba 1 paciente en el pasillo, obtuve %d", len(got))
	}
	if got[0].ID() != "P-2" {
		t.Fatalf("esperaba a P-2 en el pasillo, obtuve %s", got[0].ID())
	}
}

// TestReporteSeveros cubre la consulta 5.4: cruza los pacientes Severo con su
// numero de episodios. Se arma el historial a mano con structs de campos
// privados, sin depender de time.Now ni del flujo completo.
func TestReporteSeveros(t *testing.T) {
	h := NuevoHospital("prueba")
	severo1 := NuevoPaciente("P-1", "Ana", 30, Severo)
	severo2 := NuevoPaciente("P-2", "Beto", 40, Severo)
	leve := NuevoPaciente("P-3", "Caro", 25, Leve)
	for _, p := range []*Paciente{severo1, severo2, leve} {
		_ = h.AdmitirPaciente(p)
	}

	h.historial = []RegistroEpisodio{
		{paciente: severo1},
		{paciente: severo1},
		{paciente: severo2},
		{paciente: leve},
	}

	reporte := h.ReporteSeveros()
	if reporte[severo1] != 2 {
		t.Fatalf("esperaba 2 episodios para P-1, obtuve %d", reporte[severo1])
	}
	if reporte[severo2] != 1 {
		t.Fatalf("esperaba 1 episodio para P-2, obtuve %d", reporte[severo2])
	}
	if _, aparece := reporte[leve]; aparece {
		t.Fatalf("un paciente de nivel Leve no debe aparecer en el reporte de severos")
	}
}
