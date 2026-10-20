package hospital

import (
	"errors"
	"testing"
)

// Estos tests viven en package hospital (no en hospital_test) a proposito: al
// estar dentro del paquete pueden armar structs con campos privados a mano, sin
// depender de AdmitirPaciente ni del constructor de otro integrante.

// TestEstadosSeImprimenComoTexto verifica que String() devuelva texto y no 0, 1, 2.
func TestEstadosSeImprimenComoTexto(t *testing.T) {
	casos := []struct {
		got  string
		want string
	}{
		{Despierto.String(), "Despierto"},
		{DormidoEnPasillo.String(), "Dormido en el pasillo"},
		{DormidoEnCama.String(), "Dormido en cama"},
		{Leve.String(), "Leve"},
		{Moderado.String(), "Moderado"},
		{Severo.String(), "Severo"},
	}
	for _, c := range casos {
		if c.got != c.want {
			t.Errorf("esperaba %q, obtuve %q", c.want, c.got)
		}
	}
}

// TestPacienteNuevoEstaDespierto verifica el estado inicial y los datos que
// vienen de la Persona embutida.
func TestPacienteNuevoEstaDespierto(t *testing.T) {
	p := NuevoPaciente("P-001", "Goma jugador", 29, Severo)

	if p.Estado() != Despierto {
		t.Fatalf("un paciente nuevo debe estar Despierto, esta %v", p.Estado())
	}
	if p.Nivel() != Severo {
		t.Fatalf("esperaba nivel Severo, obtuve %v", p.Nivel())
	}
	// ID, Nombre y Edad no estan escritos en Paciente: vienen de Persona.
	if p.ID() != "P-001" || p.Nombre() != "Goma jugador" || p.Edad() != 29 {
		t.Fatalf("la Persona embutida no entrego sus datos: %s %s %d", p.ID(), p.Nombre(), p.Edad())
	}
}

// TestAtaqueDejaAlPacienteEnElPasillo verifica que el ataque cambie el estado de
// verdad (receptor puntero) y guarde la ubicacion.
func TestAtaqueDejaAlPacienteEnElPasillo(t *testing.T) {
	p := NuevoPaciente("P-001", "Goma jugador", 29, Severo)

	p.SufrirAtaqueDeSueno("cafeteria")

	if p.Estado() != DormidoEnPasillo {
		t.Fatalf("esperaba DormidoEnPasillo, obtuve %v", p.Estado())
	}
	if p.Ubicacion() != "cafeteria" {
		t.Fatalf("esperaba ubicacion cafeteria, obtuve %q", p.Ubicacion())
	}
}

// TestDespertarDosVecesDaError verifica que despertar a quien ya esta despierto
// devuelva ErrYaDespierto y no un panic.
func TestDespertarDosVecesDaError(t *testing.T) {
	p := NuevoPaciente("P-001", "Goma jugador", 29, Severo)
	p.SufrirAtaqueDeSueno("cafeteria")

	if err := p.Despertar(); err != nil {
		t.Fatalf("la primera vez no debio fallar: %v", err)
	}
	if p.Estado() != Despierto {
		t.Fatalf("tras despertar debia quedar Despierto, quedo %v", p.Estado())
	}
	if err := p.Despertar(); !errors.Is(err, ErrYaDespierto) {
		t.Fatalf("esperaba ErrYaDespierto, obtuve %v", err)
	}
}

// TestConsultaPasilloConHospitalArmadoAMano arma el hospital sin AdmitirPaciente
// y comprueba que solo salgan los dormidos en el pasillo.
func TestConsultaPasilloConHospitalArmadoAMano(t *testing.T) {
	despierto := NuevoPaciente("P-1", "Ana", 30, Leve)
	enPasillo := NuevoPaciente("P-2", "Beto", 40, Severo)
	enCama := NuevoPaciente("P-3", "Caro", 25, Moderado)
	enPasillo.SufrirAtaqueDeSueno("recepcion")
	enCama.SufrirAtaqueDeSueno("escaleras")
	enCama.estado = DormidoEnCama

	h := &Hospital{pacientes: []*Paciente{despierto, enPasillo, enCama}}

	got := h.PacientesEnPasillo()
	if len(got) != 1 || got[0] != enPasillo {
		t.Fatalf("esperaba solo a P-2 en el pasillo, obtuve %v", got)
	}
}

// TestReporteSeverosIncluyeSeverosSinEpisodios verifica que un paciente Severo
// sin ataques aparezca con 0 y que uno que no es Severo no aparezca.
func TestReporteSeverosIncluyeSeverosSinEpisodios(t *testing.T) {
	conEpisodio := NuevoPaciente("P-1", "Ana", 30, Severo)
	sinEpisodio := NuevoPaciente("P-2", "Beto", 40, Severo)
	moderado := NuevoPaciente("P-3", "Caro", 25, Moderado)

	h := &Hospital{pacientes: []*Paciente{conEpisodio, sinEpisodio, moderado}}
	h.historial = []RegistroEpisodio{{paciente: conEpisodio}, {paciente: moderado}}

	reporte := h.ReporteSeveros()

	if len(reporte) != 2 {
		t.Fatalf("esperaba 2 pacientes Severo, obtuve %d", len(reporte))
	}
	if reporte[conEpisodio] != 1 {
		t.Fatalf("esperaba 1 episodio para P-1, obtuve %d", reporte[conEpisodio])
	}
	if n, esta := reporte[sinEpisodio]; !esta || n != 0 {
		t.Fatalf("P-2 debia aparecer con 0 episodios, obtuve %d (aparece: %v)", n, esta)
	}
}

// TestSimularDiaMantieneCamasConsistentes corre la simulacion concurrente y
// verifica que ninguna cama quede ocupada de mas o de menos. Pensado para
// correrse con -race: si dos goroutines asignaran la misma habitacion, o el
// estado se tocara sin candado, el detector lo reportaria.
func TestSimularDiaMantieneCamasConsistentes(t *testing.T) {
	h := NuevoHospital("simulacion")
	_ = h.ContratarMedico(NuevoMedico("D-01", "Karen Ospina", 38, "Neurologia"))
	_ = h.ContratarCamillero(NuevoCamillero("C-01", "Wilson Castro", 30, "manana"))
	for _, num := range []int{101, 102} {
		_ = h.RegistrarHabitacion(NuevaHabitacion(num, 1))
	}
	for _, id := range []string{"P-1", "P-2", "P-3", "P-4", "P-5", "P-6"} {
		_ = h.AdmitirPaciente(NuevoPaciente(id, "Paciente "+id, 30, Severo))
	}

	h.SimularDia(5)

	enCama := 0
	for _, p := range h.pacientes {
		if p.Estado() == DormidoEnCama {
			enCama++
			if p.HabitacionActual() == nil {
				t.Fatalf("%s esta DormidoEnCama pero no tiene habitacion", p.ID())
			}
		}
	}
	ocupadas := 0
	for _, hab := range h.habitaciones {
		if !hab.EstaDisponible() {
			ocupadas++
		}
	}
	if enCama != ocupadas {
		t.Fatalf("%d pacientes en cama pero %d habitaciones ocupadas", enCama, ocupadas)
	}
	if enCama > len(h.habitaciones) {
		t.Fatalf("hay %d pacientes en cama y solo %d habitaciones", enCama, len(h.habitaciones))
	}
}
