package hospital

import (
	"errors"
	"fmt"
	"strings"
	"sync/atomic"
	"time"
)

var secuenciaEpisodios atomic.Uint64

// RegistroEpisodio guarda la atencion de un ataque, incluso sin habitacion.
// Sus campos son privados: las capas de presentacion usan metodos de lectura.
type RegistroEpisodio struct {
	id          string
	fecha       time.Time
	paciente    *Paciente
	atendidoPor Atendedor
	ubicacion   string
	habitacion  *Habitacion
}

// NuevoRegistroEpisodio valida los datos y construye un registro inmutable por convencion.
func NuevoRegistroEpisodio(id string, fecha time.Time, paciente *Paciente, atendidoPor Atendedor, ubicacion string, habitacion *Habitacion) (RegistroEpisodio, error) {
	if strings.TrimSpace(id) == "" || fecha.IsZero() || paciente == nil || atendidoPor == nil || strings.TrimSpace(ubicacion) == "" {
		return RegistroEpisodio{}, errors.New("datos incompletos del episodio")
	}
	return RegistroEpisodio{id: id, fecha: fecha, paciente: paciente, atendidoPor: atendidoPor, ubicacion: ubicacion, habitacion: habitacion}, nil
}

func nuevoEpisodio(p *Paciente, a Atendedor, ubicacion string) (RegistroEpisodio, error) {
	if p == nil || a == nil {
		return RegistroEpisodio{}, errors.New("paciente y atendedor son obligatorios")
	}
	if strings.TrimSpace(ubicacion) == "" {
		return RegistroEpisodio{}, errors.New("la ubicacion del ataque es obligatoria")
	}
	return NuevoRegistroEpisodio(fmt.Sprintf("E-%06d", secuenciaEpisodios.Add(1)), time.Now(), p, a, ubicacion, p.habitacionActual)
}

// ID identifica el episodio sin exponer sus campos.
func (e RegistroEpisodio) ID() string { return e.id }

// Fecha retorna el instante del ataque.
func (e RegistroEpisodio) Fecha() time.Time { return e.fecha }

// Paciente retorna el paciente relacionado.
func (e RegistroEpisodio) Paciente() *Paciente { return e.paciente }

// AtendidoPor retorna el personal que atendio el evento.
func (e RegistroEpisodio) AtendidoPor() Atendedor { return e.atendidoPor }

// Ubicacion indica donde ocurrio el ataque.
func (e RegistroEpisodio) Ubicacion() string { return e.ubicacion }

// Habitacion retorna la cama asignada, o nil si ninguna estaba disponible.
func (e RegistroEpisodio) Habitacion() *Habitacion { return e.habitacion }

// Resumen representa el evento con fecha, paciente, personal y cama.
func (e RegistroEpisodio) Resumen() string {
	paciente, atendedor, habitacion := "(sin paciente)", "(sin personal)", "sin habitacion"
	if e.paciente != nil {
		paciente = e.paciente.ID()
	}
	if e.atendidoPor != nil {
		atendedor = e.atendidoPor.Nombre()
	}
	if e.habitacion != nil {
		habitacion = fmt.Sprintf("habitacion %d", e.habitacion.Numero())
	}
	return fmt.Sprintf("%s %s %s (%s) -> %s [atiende: %s]", e.fecha.Format("2006-01-02 15:04"), e.id, paciente, e.ubicacion, habitacion, atendedor)
}
