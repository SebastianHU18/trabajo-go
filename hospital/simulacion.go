package hospital

import (
	"errors"
	"math/rand"
	"sync"
	"time"
)

// ubicacionesSimulacion son los lugares donde la simulacion deja caer dormidos
// a los pacientes.
var ubicacionesSimulacion = []string{
	"cafeteria",
	"sala de espera",
	"fila de radiologia",
	"pasillo del segundo piso",
	"recepcion",
	"escaleras",
}

// SimularDia lanza una goroutine por paciente. Cada una provoca ataques en
// intervalos aleatorios y el hospital los atiende. Espera a que todas terminen.
// Los errores esperables (sin habitacion, sin personal) no detienen la
// simulacion: un paciente sin cama simplemente se queda en el pasillo.
func (h *Hospital) SimularDia(ataquesPorPaciente int) {
	// Se copia la lista bajo el candado para no recorrer h.pacientes mientras
	// otra goroutine admite un paciente.
	h.mu.Lock()
	pacientes := make([]*Paciente, len(h.pacientes))
	copy(pacientes, h.pacientes)
	h.mu.Unlock()

	var wg sync.WaitGroup
	for _, p := range pacientes {
		wg.Add(1) // Add antes de lanzar, nunca dentro de la goroutine

		// p entra como parametro: con go 1.21 la variable del ciclo es una sola
		// y todas las goroutines verian al ultimo paciente.
		go func(p *Paciente) {
			defer wg.Done()

			for i := 0; i < ataquesPorPaciente; i++ {
				time.Sleep(time.Duration(rand.Intn(80)) * time.Millisecond)
				_ = h.ManejarAtaque(p, ubicacionAleatoria())
			}
		}(p)
	}
	wg.Wait()
}

// ManejarAtaque provoca un ataque de sueno a un paciente y lo hace atender.
//
// Si el paciente seguia dormido de un ataque anterior, primero lo despierta: eso
// libera su cama y evita que ocupe dos habitaciones a la vez. Despertar y
// despachar son operaciones del hospital que toman el candado completas, asi que
// buscar y ocupar una cama ocurren en la misma seccion critica y dos goroutines
// no pueden quedarse con la misma habitacion.
func (h *Hospital) ManejarAtaque(p *Paciente, ubicacion string) error {
	if err := h.DespertarPaciente(p); err != nil && !errors.Is(err, ErrYaDespierto) {
		return err
	}
	_, err := h.DespacharAtencion(p, ubicacion)
	return err
}

// ubicacionAleatoria elige un lugar de la simulacion al azar. rand.Intn del
// paquete es seguro para usar desde varias goroutines.
func ubicacionAleatoria() string {
	return ubicacionesSimulacion[rand.Intn(len(ubicacionesSimulacion))]
}
