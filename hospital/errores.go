package hospital

import "errors"

// Errores previsibles del hospital. Se declaran aqui, como variables del
// paquete, para que las pruebas los comparen con errors.Is en vez de comparar
// cadenas de texto frágiles. Este archivo se escribe en la fase 0 y despues no
// se vuelve a tocar: es el contrato de errores que comparten los tres bloques.
var (
	// ErrSinHabitacion se devuelve cuando no queda ninguna habitacion libre.
	ErrSinHabitacion = errors.New("no hay habitacion disponible")

	// ErrPacienteRepetido se devuelve al admitir un paciente ya admitido.
	ErrPacienteRepetido = errors.New("el paciente ya esta admitido")

	// ErrPacienteNoAdmitido se devuelve al operar sobre un paciente que el
	// hospital no conoce.
	ErrPacienteNoAdmitido = errors.New("el paciente no esta admitido en el hospital")

	// ErrYaDespierto se devuelve al despertar a quien ya estaba despierto.
	ErrYaDespierto = errors.New("el paciente ya esta despierto")

	// ErrHabitacionOcupada se devuelve al ocupar una habitacion sin cupo.
	ErrHabitacionOcupada = errors.New("la habitacion ya esta ocupada")

	// ErrPacienteNoEnHabitacion se devuelve al liberar a un paciente que no
	// ocupa esa habitacion.
	ErrPacienteNoEnHabitacion = errors.New("el paciente no ocupa esa habitacion")

	// ErrSinPersonal se devuelve al despachar una atencion sin personal.
	ErrSinPersonal = errors.New("no hay personal disponible para atender")

	// ErrPacienteNoDespierto se devuelve al despachar un ataque de un paciente
	// que no esta despierto: ya esta dormido y asignarle otra cama lo dejaria
	// ocupando dos habitaciones.
	ErrPacienteNoDespierto = errors.New("el paciente no esta despierto")
)
