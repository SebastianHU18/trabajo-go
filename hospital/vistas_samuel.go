package hospital

import "strings"

// VistaPaciente es una proyeccion de lectura para la API, no una entidad de dominio.
type VistaPaciente struct {
	ID        string `json:"id"`
	Nombre    string `json:"nombre"`
	Nivel     string `json:"nivel"`
	Estado    string `json:"estado"`
	Ubicacion string `json:"ubicacion"`
}

// VistaHabitacion es una proyeccion de lectura del inventario de camas.
type VistaHabitacion struct {
	Numero    int    `json:"numero"`
	Capacidad int    `json:"capacidad"`
	Ocupantes int    `json:"ocupantes"`
	Estado    string `json:"estado"`
}

// VistaMedico es una proyeccion de lectura del personal medico.
type VistaMedico struct {
	ID           string `json:"id"`
	Nombre       string `json:"nombre"`
	Especialidad string `json:"especialidad"`
}

// VistaEpisodio contiene los datos de un episodio que necesita la interfaz.
type VistaEpisodio struct {
	ID         string `json:"id"`
	Fecha      string `json:"fecha"`
	PacienteID string `json:"pacienteId"`
	Atendedor  string `json:"atendedor"`
	Ubicacion  string `json:"ubicacion"`
	Habitacion *int   `json:"habitacion"`
	Resumen    string `json:"resumen"`
}

// VistaPacientes copia el estado visible de los pacientes bajo el candado.
func (h *Hospital) VistaPacientes() []VistaPaciente {
	h.mu.Lock()
	defer h.mu.Unlock()
	resultado := make([]VistaPaciente, 0, len(h.pacientes))
	for _, p := range h.pacientes {
		resultado = append(resultado, vistaPaciente(p))
	}
	return resultado
}

// VistaPacientesEnPasillo devuelve solo los pacientes dormidos en pasillo.
func (h *Hospital) VistaPacientesEnPasillo() []VistaPaciente {
	h.mu.Lock()
	defer h.mu.Unlock()
	resultado := make([]VistaPaciente, 0)
	for _, p := range h.pacientes {
		if p.estado == DormidoEnPasillo {
			resultado = append(resultado, vistaPaciente(p))
		}
	}
	return resultado
}

func vistaPaciente(p *Paciente) VistaPaciente {
	return VistaPaciente{ID: p.ID(), Nombre: p.Nombre(), Nivel: p.nivel.String(), Estado: p.estado.String(), Ubicacion: p.ubicacionActual}
}

// VistaHabitaciones copia el estado de las camas sin exponer punteros privados.
func (h *Hospital) VistaHabitaciones() []VistaHabitacion {
	h.mu.Lock()
	defer h.mu.Unlock()
	resultado := make([]VistaHabitacion, 0, len(h.habitaciones))
	for _, hab := range h.habitaciones {
		resultado = append(resultado, VistaHabitacion{Numero: hab.numero, Capacidad: hab.capacidad, Ocupantes: len(hab.ocupantes), Estado: hab.estado.String()})
	}
	return resultado
}

// VistaMedicos copia el listado del personal medico.
func (h *Hospital) VistaMedicos() []VistaMedico {
	h.mu.Lock()
	defer h.mu.Unlock()
	resultado := make([]VistaMedico, 0, len(h.medicos))
	for _, m := range h.medicos {
		resultado = append(resultado, VistaMedico{ID: m.ID(), Nombre: m.Nombre(), Especialidad: m.Especialidad()})
	}
	return resultado
}

// VistaEpisodiosDeMedico retorna registros visibles de un medico por su ID.
func (h *Hospital) VistaEpisodiosDeMedico(id string) ([]VistaEpisodio, bool) {
	h.mu.Lock()
	defer h.mu.Unlock()
	for _, m := range h.medicos {
		if strings.EqualFold(m.ID(), id) {
			registros := m.MisEpisodios()
			resultado := make([]VistaEpisodio, 0, len(registros))
			for _, r := range registros {
				resultado = append(resultado, vistaEpisodio(r))
			}
			return resultado, true
		}
	}
	return nil, false
}

func vistaEpisodio(r RegistroEpisodio) VistaEpisodio {
	vista := VistaEpisodio{ID: r.ID(), Fecha: r.Fecha().Format("2006-01-02 15:04:05"), Ubicacion: r.Ubicacion(), Resumen: r.Resumen()}
	if r.Paciente() != nil {
		vista.PacienteID = r.Paciente().ID()
	}
	if r.AtendidoPor() != nil {
		vista.Atendedor = r.AtendidoPor().Nombre()
	}
	if r.Habitacion() != nil {
		numero := r.Habitacion().Numero()
		vista.Habitacion = &numero
	}
	return vista
}

// VistaSeveros resume la cantidad de episodios por paciente grave.
// Utiliza la misma fuente de datos que el reporte obligatorio del grupo.
func (h *Hospital) VistaSeveros() map[string]int {
	reporte := h.ReporteSeveros()
	vista := make(map[string]int, len(reporte))
	for p, cantidad := range reporte {
		vista[p.ID()] = cantidad
	}
	return vista
}

// AtacarPorID realiza una emergencia identificando al paciente sin exponer listas.
func (h *Hospital) AtacarPorID(id, ubicacion string) (RegistroEpisodio, error) {
	p := h.buscarPacientePorID(id)
	if p == nil {
		return RegistroEpisodio{}, ErrPacienteNoAdmitido
	}
	return h.DespacharAtencion(p, ubicacion)
}

// DespertarPorID despierta a un paciente y libera su habitacion.
func (h *Hospital) DespertarPorID(id string) error {
	p := h.buscarPacientePorID(id)
	if p == nil {
		return ErrPacienteNoAdmitido
	}
	return h.DespertarPaciente(p)
}

func (h *Hospital) buscarPacientePorID(id string) *Paciente {
	h.mu.Lock()
	defer h.mu.Unlock()
	for _, p := range h.pacientes {
		if strings.EqualFold(p.ID(), id) {
			return p
		}
	}
	return nil
}
