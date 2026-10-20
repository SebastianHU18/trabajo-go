package hospital

import "errors"

// Camillero puede asistir a un paciente y trasladarlo a una cama asignada,
// pero no tiene metodos de diagnostico.
type Camillero struct {
	Persona
	turno string
}

// NuevoCamillero crea personal de traslado con su turno.
func NuevoCamillero(id, nombre string, edad int, turno string) *Camillero {
	return &Camillero{Persona: Persona{id: id, nombre: nombre, edad: edad}, turno: turno}
}

// Turno devuelve el horario de trabajo del camillero.
func (c *Camillero) Turno() string { return c.turno }

// Atender registra la intervencion del camillero sin diagnosticar al paciente.
// El Hospital ya ha colocado al paciente en una cama, o lo mantiene en pasillo.
func (c *Camillero) Atender(p *Paciente, ubicacion string) (RegistroEpisodio, error) {
	if c == nil {
		return RegistroEpisodio{}, errors.New("camillero obligatorio")
	}
	return nuevoEpisodio(p, c, ubicacion)
}
