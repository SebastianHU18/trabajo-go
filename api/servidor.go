// Package api expone las operaciones del hospital mediante HTTP y JSON.
// No introduce dependencias de HTTP dentro del modelo de dominio.
package api

import (
	"net/http"

	"costenos-narcolepsia/hospital"
)

// NuevoServidor construye las rutas HTTP usando solo la biblioteca estandar.
func NuevoServidor(h *hospital.Hospital) http.Handler {
	s := &servidor{hospital: h}
	mux := http.NewServeMux()
	mux.HandleFunc("/health", s.salud)
	mux.HandleFunc("/api/pacientes", s.pacientes)
	mux.HandleFunc("/api/pacientes/", s.pacienteIndividual)
	mux.HandleFunc("/api/medicos", s.medicos)
	mux.HandleFunc("/api/medicos/", s.episodiosMedico)
	mux.HandleFunc("/api/habitaciones", s.habitaciones)
	mux.HandleFunc("/api/reportes/severos", s.severos)
	mux.HandleFunc("/", s.pagina)
	return mux
}
