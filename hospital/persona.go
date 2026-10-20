package hospital

// Persona agrupa los datos de identidad que comparten pacientes, medicos y
// camilleros. Se embute en cada uno en vez de repetir id, nombre y edad: asi
// reutilizan ID(), Nombre() y Edad() sin escribirlos otra vez.
type Persona struct {
	id     string
	nombre string
	edad   int
}

// NuevaPersona construye una Persona. Devuelve el valor y no un puntero porque
// quien la usa la embute por valor: cada entidad tiene su propia copia y nunca
// hay una Persona nula que haga estallar Nombre().
func NuevaPersona(id, nombre string, edad int) Persona {
	return Persona{id: id, nombre: nombre, edad: edad}
}

// ID devuelve el identificador de la persona.
// Receptor por valor: solo lee, y asi el metodo se promueve tanto a T como a *T,
// que es lo que permite a *Medico y *Camillero cumplir la interfaz Atendedor.
func (p Persona) ID() string { return p.id }

// Nombre devuelve el nombre de la persona.
func (p Persona) Nombre() string { return p.nombre }

// Edad devuelve la edad de la persona.
func (p Persona) Edad() int { return p.edad }
