package hospital

import (
	"errors"
	"strings"
	"sync"
)

// Medico reutiliza la identidad de Persona y registra sus atenciones.
type Medico struct {
	Persona
	especialidad string
	pacientes    []*Paciente
	episodios    []RegistroEpisodio
	mu           sync.Mutex
}

// NuevoMedico construye un medico con identidad y especialidad.
func NuevoMedico(id, nombre string, edad int, especialidad string) *Medico {
	return &Medico{Persona: Persona{id: id, nombre: nombre, edad: edad}, especialidad: especialidad}
}

// Especialidad devuelve el area de trabajo del medico.
func (m *Medico) Especialidad() string { return m.especialidad }

// Pacientes devuelve una copia de los pacientes diagnosticados.
func (m *Medico) Pacientes() []*Paciente {
	m.mu.Lock()
	defer m.mu.Unlock()
	return append([]*Paciente(nil), m.pacientes...)
}

// DiagnosticarPaciente agrega al paciente bajo cuidado sin duplicarlo.
func (m *Medico) DiagnosticarPaciente(p *Paciente) error {
	if m == nil || p == nil {
		return errors.New("medico y paciente son obligatorios")
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	for _, actual := range m.pacientes {
		if strings.EqualFold(actual.ID(), p.ID()) {
			return nil
		}
	}
	p.AsignarMedico(m)
	m.pacientes = append(m.pacientes, p)
	return nil
}

// AtenderEmergenciaSueno atiende un episodio usando la ubicacion del paciente.
func (m *Medico) AtenderEmergenciaSueno(p *Paciente) error {
	if p == nil {
		return errors.New("paciente obligatorio")
	}
	_, err := m.Atender(p, p.Ubicacion())
	return err
}

// Atender satisface Atendedor y guarda el episodio en el historial propio.
// La asignacion de cama corresponde a Hospital y ocurre antes de esta llamada.
func (m *Medico) Atender(p *Paciente, ubicacion string) (RegistroEpisodio, error) {
	if m == nil {
		return RegistroEpisodio{}, errors.New("medico obligatorio")
	}
	registro, err := nuevoEpisodio(p, m, ubicacion)
	if err != nil {
		return RegistroEpisodio{}, err
	}
	m.mu.Lock()
	m.episodios = append(m.episodios, registro)
	m.mu.Unlock()
	return registro, nil
}

// MisEpisodios devuelve una copia de los episodios atendidos por el medico.
func (m *Medico) MisEpisodios() []RegistroEpisodio {
	m.mu.Lock()
	defer m.mu.Unlock()
	return append([]RegistroEpisodio(nil), m.episodios...)
}
