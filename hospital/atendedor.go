package hospital

// Atendedor abstrae al personal que atiende emergencias de sueno.
// El hospital puede usar medicos y camilleros sin consultar su tipo concreto.
type Atendedor interface {
	ID() string
	Nombre() string
	Atender(p *Paciente, ubicacion string) (RegistroEpisodio, error)
}

// Las aserciones estaticas detectan cambios incompatibles de firmas.
var (
	_ Atendedor = (*Medico)(nil)
	_ Atendedor = (*Camillero)(nil)
)
