package hospital

import (
	"errors"
	"fmt"
)

// EstadoHabitacion identifica el estado de ocupacion de una habitacion.
type EstadoHabitacion int

const (
	Disponible EstadoHabitacion = iota
	Ocupada
)

// String representa el estado de habitacion de forma legible.
func (e EstadoHabitacion) String() string {
	switch e {
	case Disponible:
		return "Disponible"
	case Ocupada:
		return "Ocupada"
	default:
		return "Desconocido"
	}
}

// Habitacion administra la capacidad y los pacientes que ocupan una cama.
type Habitacion struct {
	numero    int
	capacidad int
	estado    EstadoHabitacion
	ocupantes []*Paciente
}

// NuevaHabitacion crea una habitacion sin pacientes.
// La capacidad no valida se rechaza al intentar ocupar, sin inventar cupos.
func NuevaHabitacion(numero, capacidad int) *Habitacion {
	return &Habitacion{numero: numero, capacidad: capacidad, estado: Disponible}
}

// Numero devuelve el numero identificador de la habitacion.
func (h *Habitacion) Numero() int { return h.numero }

// Capacidad devuelve el numero total de camas.
func (h *Habitacion) Capacidad() int { return h.capacidad }

// Estado indica si hay al menos un ocupante en la habitacion.
func (h *Habitacion) Estado() EstadoHabitacion { return h.estado }

// Ocupantes entrega una copia de la lista sin exponer el slice interno.
func (h *Habitacion) Ocupantes() []*Paciente {
	return append([]*Paciente(nil), h.ocupantes...)
}

// EstaDisponible informa si queda al menos una cama libre.
func (h *Habitacion) EstaDisponible() bool {
	return h != nil && h.capacidad > 0 && len(h.ocupantes) < h.capacidad
}

// Ocupar asigna una cama al paciente, rechazando duplicados y falta de cupo.
// El hospital gestiona los cambios de estado del paciente.
func (h *Habitacion) Ocupar(p *Paciente) error {
	if h == nil || p == nil {
		return errors.New("habitacion y paciente son obligatorios")
	}
	for _, existente := range h.ocupantes {
		if existente == p || existente.ID() == p.ID() {
			return fmt.Errorf("paciente ya ocupa la habitacion %d: %w", h.numero, ErrHabitacionOcupada)
		}
	}
	if !h.EstaDisponible() {
		return fmt.Errorf("habitacion %d sin cupo: %w", h.numero, ErrHabitacionOcupada)
	}
	h.ocupantes = append(h.ocupantes, p)
	h.estado = Ocupada
	return nil
}

// Liberar desocupa la cama del paciente indicado.
func (h *Habitacion) Liberar(p *Paciente) error {
	if h == nil || p == nil {
		return ErrPacienteNoEnHabitacion
	}
	for i, actual := range h.ocupantes {
		if actual == p {
			h.ocupantes = append(h.ocupantes[:i], h.ocupantes[i+1:]...)
			if len(h.ocupantes) == 0 {
				h.estado = Disponible
			}
			return nil
		}
	}
	return ErrPacienteNoEnHabitacion
}
