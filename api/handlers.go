package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"costenos-narcolepsia/hospital"
)

type servidor struct{ hospital *hospital.Hospital }

func escribirJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Cache-Control", "no-store")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}

func metodo(w http.ResponseWriter, r *http.Request, esperado string) bool {
	if r.Method == esperado {
		return true
	}
	w.Header().Set("Allow", esperado)
	escribirJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "metodo no permitido"})
	return false
}

func (s *servidor) salud(w http.ResponseWriter, r *http.Request) {
	if !metodo(w, r, http.MethodGet) {
		return
	}
	escribirJSON(w, http.StatusOK, map[string]string{"estado": "ok"})
}

func (s *servidor) pacientes(w http.ResponseWriter, r *http.Request) {
	if !metodo(w, r, http.MethodGet) {
		return
	}
	escribirJSON(w, http.StatusOK, s.hospital.VistaPacientes())
}

func (s *servidor) habitaciones(w http.ResponseWriter, r *http.Request) {
	if !metodo(w, r, http.MethodGet) {
		return
	}
	escribirJSON(w, http.StatusOK, s.hospital.VistaHabitaciones())
}

func (s *servidor) medicos(w http.ResponseWriter, r *http.Request) {
	if !metodo(w, r, http.MethodGet) {
		return
	}
	escribirJSON(w, http.StatusOK, s.hospital.VistaMedicos())
}

func (s *servidor) severos(w http.ResponseWriter, r *http.Request) {
	if !metodo(w, r, http.MethodGet) {
		return
	}
	escribirJSON(w, http.StatusOK, s.hospital.VistaSeveros())
}

// episodioJSON es la vista de un episodio del historial general. Se arma aqui,
// con los metodos de lectura publicos, para no agregar JSON al paquete hospital.
type episodioJSON struct {
	ID         string `json:"id"`
	Fecha      string `json:"fecha"`
	PacienteID string `json:"pacienteId"`
	Atendedor  string `json:"atendedor"`
	Ubicacion  string `json:"ubicacion"`
	Habitacion *int   `json:"habitacion"`
	Resumen    string `json:"resumen"`
}

// episodios devuelve el historial de todo el hospital, incluidos los episodios
// atendidos por camilleros, que no aparecen en el historial de ningun medico.
func (s *servidor) episodios(w http.ResponseWriter, r *http.Request) {
	if !metodo(w, r, http.MethodGet) {
		return
	}
	historial := s.hospital.Historial()
	vista := make([]episodioJSON, 0, len(historial))
	for _, e := range historial {
		item := episodioJSON{ID: e.ID(), Fecha: e.Fecha().Format("2006-01-02 15:04:05"), Ubicacion: e.Ubicacion(), Resumen: e.Resumen()}
		if e.Paciente() != nil {
			item.PacienteID = e.Paciente().ID()
		}
		if e.AtendidoPor() != nil {
			item.Atendedor = e.AtendidoPor().Nombre()
		}
		if e.Habitacion() != nil {
			numero := e.Habitacion().Numero()
			item.Habitacion = &numero
		}
		vista = append(vista, item)
	}
	escribirJSON(w, http.StatusOK, vista)
}

func (s *servidor) pacienteIndividual(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/api/pacientes/pasillo" {
		if !metodo(w, r, http.MethodGet) {
			return
		}
		escribirJSON(w, http.StatusOK, s.hospital.VistaPacientesEnPasillo())
		return
	}
	trozos := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(trozos) != 4 || trozos[0] != "api" || trozos[1] != "pacientes" || trozos[2] == "" {
		escribirJSON(w, http.StatusNotFound, map[string]string{"error": "ruta no encontrada"})
		return
	}
	if !metodo(w, r, http.MethodPost) {
		return
	}
	var err error
	switch trozos[3] {
	case "ataque":
		r.Body = http.MaxBytesReader(w, r.Body, 4096)
		var cuerpo struct {
			Ubicacion string `json:"ubicacion"`
		}
		if e := json.NewDecoder(r.Body).Decode(&cuerpo); e != nil || strings.TrimSpace(cuerpo.Ubicacion) == "" {
			escribirJSON(w, http.StatusBadRequest, map[string]string{"error": "se requiere ubicacion"})
			return
		}
		var reg hospital.RegistroEpisodio
		reg, err = s.hospital.AtacarPorID(trozos[2], cuerpo.Ubicacion)
		if err == nil || errors.Is(err, hospital.ErrSinHabitacion) {
			respuesta := map[string]any{"episodio": reg.Resumen(), "camaAsignada": err == nil}
			if err != nil {
				respuesta["aviso"] = err.Error()
			}
			escribirJSON(w, http.StatusOK, respuesta)
			return
		}
	case "despertar":
		err = s.hospital.DespertarPorID(trozos[2])
		if err == nil {
			escribirJSON(w, http.StatusOK, map[string]string{"resultado": "paciente despierto"})
			return
		}
	default:
		escribirJSON(w, http.StatusNotFound, map[string]string{"error": "ruta no encontrada"})
		return
	}
	status := http.StatusConflict
	if errors.Is(err, hospital.ErrPacienteNoAdmitido) {
		status = http.StatusNotFound
	}
	escribirJSON(w, status, map[string]string{"error": err.Error()})
}

func (s *servidor) episodiosMedico(w http.ResponseWriter, r *http.Request) {
	if !metodo(w, r, http.MethodGet) {
		return
	}
	trozos := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(trozos) != 4 || trozos[0] != "api" || trozos[1] != "medicos" || trozos[3] != "episodios" {
		escribirJSON(w, http.StatusNotFound, map[string]string{"error": "ruta no encontrada"})
		return
	}
	episodios, ok := s.hospital.VistaEpisodiosDeMedico(trozos[2])
	if !ok {
		escribirJSON(w, http.StatusNotFound, map[string]string{"error": "medico no encontrado"})
		return
	}
	escribirJSON(w, http.StatusOK, episodios)
}

func (s *servidor) pagina(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	if !metodo(w, r, http.MethodGet) {
		return
	}
	http.ServeFile(w, r, "web/index.html")
}
