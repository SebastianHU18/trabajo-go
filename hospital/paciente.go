package hospital

// Paciente representa a un paciente del hospital. Persona va embutida por
// valor: el paciente reutiliza ID(), Nombre() y Edad() sin repetir los campos.
//
// Todos los campos son privados. El hospital (mismo paquete) mueve al paciente
// a una cama escribiendo estado y habitacionActual; desde fuera del paquete
// solo se puede leer por los metodos de abajo.
type Paciente struct {
	Persona
	nivel            NivelNarcolepsia
	estado           EstadoPaciente
	ubicacionActual  string
	medicoAsignado   *Medico
	habitacionActual *Habitacion // la cama que ocupa; nil si no tiene
}

// NuevoPaciente construye un paciente despierto. Devuelve puntero porque el
// hospital, los medicos y los episodios comparten al mismo paciente.
func NuevoPaciente(id, nombre string, edad int, nivel NivelNarcolepsia) *Paciente {
	return &Paciente{
		Persona: NuevaPersona(id, nombre, edad),
		nivel:   nivel,
		estado:  Despierto,
	}
}

// SufrirAtaqueDeSueno deja al paciente dormido en el pasillo, en la ubicacion
// donde le dio el ataque. Se usa *Paciente porque el metodo modifica: con
// receptor por valor mutaria una copia y el estado real nunca cambiaria.
func (p *Paciente) SufrirAtaqueDeSueno(ubicacion string) {
	p.estado = DormidoEnPasillo
	p.ubicacionActual = ubicacion
}

// Despertar devuelve al paciente al estado Despierto y falla con ErrYaDespierto
// si ya lo estaba. NO libera la habitacion: eso lo hace el hospital, que es el
// unico que conoce las camas.
func (p *Paciente) Despertar() error {
	if p.estado == Despierto {
		return ErrYaDespierto
	}
	p.estado = Despierto
	return nil
}

// AsignarMedico registra al medico que atiende al paciente.
func (p *Paciente) AsignarMedico(m *Medico) {
	p.medicoAsignado = m
}

// Estado devuelve el estado actual del paciente.
func (p Paciente) Estado() EstadoPaciente { return p.estado }

// Nivel devuelve el nivel de narcolepsia del paciente.
func (p Paciente) Nivel() NivelNarcolepsia { return p.nivel }

// Ubicacion devuelve el lugar donde le dio el ultimo ataque.
func (p Paciente) Ubicacion() string { return p.ubicacionActual }

// MedicoAsignado devuelve el medico que lo atiende, o nil si no tiene.
func (p Paciente) MedicoAsignado() *Medico { return p.medicoAsignado }

// HabitacionActual devuelve la habitacion que ocupa, o nil si no tiene cama.
func (p Paciente) HabitacionActual() *Habitacion { return p.habitacionActual }
