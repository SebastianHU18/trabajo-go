package hospital

// EstadoPaciente indica donde esta un paciente respecto a su ataque de sueno.
// Es un tipo con constantes y no una cadena suelta: el compilador atrapa los
// errores de escritura que con "Dormido_En_Pasillo" pasarian desapercibidos.
type EstadoPaciente int

const (
	// Despierto es el estado inicial del paciente.
	Despierto EstadoPaciente = iota
	// DormidoEnPasillo es el de quien se durmio y aun no tiene cama.
	DormidoEnPasillo
	// DormidoEnCama es el de quien ya tiene habitacion asignada.
	DormidoEnCama
)

// String devuelve el estado como texto legible. Sin este metodo los estados se
// imprimirian como 0, 1 y 2.
func (e EstadoPaciente) String() string {
	switch e {
	case Despierto:
		return "Despierto"
	case DormidoEnPasillo:
		return "Dormido en el pasillo"
	case DormidoEnCama:
		return "Dormido en cama"
	default:
		return "Estado desconocido"
	}
}

// NivelNarcolepsia clasifica la gravedad de la narcolepsia de un paciente.
type NivelNarcolepsia int

const (
	// Leve es el nivel menos grave.
	Leve NivelNarcolepsia = iota
	// Moderado es el nivel intermedio.
	Moderado
	// Severo es el nivel mas grave; es el que cruza el reporte 5.4.
	Severo
)

// String devuelve el nivel como texto legible.
func (n NivelNarcolepsia) String() string {
	switch n {
	case Leve:
		return "Leve"
	case Moderado:
		return "Moderado"
	case Severo:
		return "Severo"
	default:
		return "Nivel desconocido"
	}
}
