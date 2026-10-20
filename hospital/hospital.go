package hospital

import (
	"strings"
	"sync"
)

// Hospital es el nucleo del sistema: coordina pacientes, personal, habitaciones
// y el historial de episodios, y es el unico que conoce las camas.
//
// Lleva un sync.Mutex para el bono de concurrencia, asi que TODOS sus metodos
// usan receptor puntero: copiar el struct copiaria el mutex y go vet lo reporta.
type Hospital struct {
	nombre       string
	medicos      []*Medico
	pacientes    []*Paciente
	habitaciones []*Habitacion
	// personal guarda al staff como []Atendedor (no como tipos concretos) para
	// poder despachar por polimorfismo sin preguntar quien es cada uno.
	personal  []Atendedor
	historial []RegistroEpisodio
	turno     int        // indice round-robin del siguiente atendedor
	mu        sync.Mutex // protege todo el estado del hospital
}

// NuevoHospital construye un hospital vacio con su nombre.
func NuevoHospital(nombre string) *Hospital {
	return &Hospital{nombre: nombre}
}

// Nombre devuelve el nombre del hospital.
func (h *Hospital) Nombre() string { return h.nombre }

// AdmitirPaciente registra un paciente nuevo y rechaza los id repetidos.
func (h *Hospital) AdmitirPaciente(p *Paciente) error {
	h.mu.Lock()
	defer h.mu.Unlock()

	for _, existente := range h.pacientes {
		if strings.EqualFold(existente.ID(), p.ID()) {
			return ErrPacienteRepetido
		}
	}
	h.pacientes = append(h.pacientes, p)
	return nil
}

// ContratarMedico suma un medico al hospital y lo registra como personal que
// puede atender ataques de sueno.
func (h *Hospital) ContratarMedico(m *Medico) error {
	h.mu.Lock()
	defer h.mu.Unlock()

	h.medicos = append(h.medicos, m)
	h.personal = append(h.personal, m)
	return nil
}

// ContratarCamillero suma un camillero al personal que puede atender ataques.
func (h *Hospital) ContratarCamillero(c *Camillero) error {
	h.mu.Lock()
	defer h.mu.Unlock()

	h.personal = append(h.personal, c)
	return nil
}

// ContratarPersonal agrega cualquier Atendedor al staff sin mirar su tipo
// concreto. Permite enganchar en vivo un tipo nuevo (por ejemplo una enfermera)
// que cumpla la interfaz, sin tocar el resto del hospital.
func (h *Hospital) ContratarPersonal(a Atendedor) error {
	h.mu.Lock()
	defer h.mu.Unlock()

	h.personal = append(h.personal, a)
	return nil
}

// RegistrarHabitacion suma una habitacion al hospital.
func (h *Hospital) RegistrarHabitacion(hab *Habitacion) error {
	h.mu.Lock()
	defer h.mu.Unlock()

	h.habitaciones = append(h.habitaciones, hab)
	return nil
}

// AsignarHabitacion busca la primera habitacion libre y la asigna al paciente.
// Si no hay ninguna devuelve ErrSinHabitacion y el paciente se queda en el
// pasillo. Consulta 5.3 del enunciado.
func (h *Hospital) AsignarHabitacion(p *Paciente) (*Habitacion, error) {
	h.mu.Lock()
	defer h.mu.Unlock()

	return h.asignarHabitacionSinCandado(p)
}

// asignarHabitacionSinCandado hace el trabajo real de buscar y ocupar una cama,
// SIN tomar el candado. Existe aparte porque sync.Mutex no es reentrante:
// AsignarHabitacion toma el candado y delega aqui, y DespacharAtencion, que ya
// lo tiene tomado, la reusa sin volver a bloquear.
func (h *Hospital) asignarHabitacionSinCandado(p *Paciente) (*Habitacion, error) {
	for _, hab := range h.habitaciones {
		if hab.EstaDisponible() {
			if err := hab.Ocupar(p); err != nil {
				continue
			}
			// El hospital es el unico que mueve al paciente a una cama: cambia
			// su estado y recuerda la habitacion para poder liberarla despues.
			p.estado = DormidoEnCama
			p.habitacionActual = hab
			return hab, nil
		}
	}
	return nil, ErrSinHabitacion
}

// DespacharAtencion marca el ataque del paciente, envia al primer atendedor del
// turno sin conocer su tipo concreto (aqui se cobra el polimorfismo) y le busca
// cama. Devuelve el episodio y, si no hubo cama, ErrSinHabitacion; en ese caso
// el paciente queda en el pasillo pero igual queda atendido y registrado.
func (h *Hospital) DespacharAtencion(p *Paciente, ubicacion string) (RegistroEpisodio, error) {
	h.mu.Lock()
	defer h.mu.Unlock()

	// Solo se atiende el ataque de quien esta despierto. Si el paciente ya
	// estaba dormido (en pasillo o en cama), despacharlo de nuevo le asignaria
	// una segunda habitacion sin liberar la primera; por eso se rechaza. Para
	// re-atenderlo primero hay que despertarlo (eso libera su cama).
	if p.Estado() != Despierto {
		return RegistroEpisodio{}, ErrPacienteNoDespierto
	}

	atendedor := h.siguienteAtendedor()
	if atendedor == nil {
		return RegistroEpisodio{}, ErrSinPersonal
	}

	p.SufrirAtaqueDeSueno(ubicacion)
	// Se asigna la cama ANTES de atender para que el atendedor registre en el
	// episodio la habitacion que quedo (o ninguna). errCama se guarda y se
	// devuelve al final sin cortar el flujo: el paciente se atiende igual.
	_, errCama := h.asignarHabitacionSinCandado(p)

	reg, err := atendedor.Atender(p, ubicacion)
	if err != nil {
		return reg, err
	}
	h.historial = append(h.historial, reg)
	return reg, errCama
}

// DespertarPaciente despierta al paciente y, si estaba en una cama, la libera.
// La liberacion vive aqui y no en Paciente porque solo el hospital conoce las
// habitaciones.
func (h *Hospital) DespertarPaciente(p *Paciente) error {
	h.mu.Lock()
	defer h.mu.Unlock()

	if err := p.Despertar(); err != nil {
		return err
	}
	if p.habitacionActual != nil {
		_ = p.habitacionActual.Liberar(p)
		p.habitacionActual = nil
	}
	return nil
}

// RegistrarEpisodio agrega un episodio al historial del hospital.
func (h *Hospital) RegistrarEpisodio(reg RegistroEpisodio) error {
	h.mu.Lock()
	defer h.mu.Unlock()

	h.historial = append(h.historial, reg)
	return nil
}

// Historial devuelve una copia del historial de episodios del hospital, para
// que quien lo lea no pueda modificar el original.
func (h *Hospital) Historial() []RegistroEpisodio {
	h.mu.Lock()
	defer h.mu.Unlock()

	copia := make([]RegistroEpisodio, len(h.historial))
	copy(copia, h.historial)
	return copia
}

// siguienteAtendedor devuelve al proximo atendedor en rotacion. Es privado y
// NO toma el candado: solo se llama desde metodos que ya lo tienen tomado.
func (h *Hospital) siguienteAtendedor() Atendedor {
	if len(h.personal) == 0 {
		return nil
	}
	a := h.personal[h.turno%len(h.personal)]
	h.turno++
	return a
}
