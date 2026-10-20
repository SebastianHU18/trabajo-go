package api

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"testing"

	"costenos-narcolepsia/hospital"
)

func TestEndpointsPacientesAtaquesHabitaciones(t *testing.T) {
	h := hospital.NuevoHospital("prueba")
	_ = h.ContratarMedico(hospital.NuevoMedico("D-1", "Karen", 35, "Neurologia"))
	_ = h.RegistrarHabitacion(hospital.NuevaHabitacion(101, 1))
	_ = h.AdmitirPaciente(hospital.NuevoPaciente("P-1", "Ana", 29, hospital.Severo))
	_ = h.AdmitirPaciente(hospital.NuevoPaciente("P-2", "Beto", 27, hospital.Leve))
	s := NuevoServidor(h)
	check := func(metodo, ruta string, cuerpo []byte, esperado int) []byte {
		t.Helper()
		r := httptest.NewRequest(metodo, ruta, bytes.NewReader(cuerpo))
		w := httptest.NewRecorder()
		s.ServeHTTP(w, r)
		if w.Code != esperado {
			t.Fatalf("%s %s: esperaba %d, obtuvo %d: %s", metodo, ruta, esperado, w.Code, w.Body.String())
		}
		return w.Body.Bytes()
	}
	check("GET", "/health", nil, 200)
	check("GET", "/api/pacientes", nil, 200)
	check("GET", "/api/habitaciones", nil, 200)
	check("GET", "/api/medicos", nil, 200)
	check("GET", "/api/reportes/severos", nil, 200)
	check("GET", "/api/pacientes/pasillo", nil, 200)
	check("POST", "/api/pacientes/P-1/ataque", []byte(`{"ubicacion":"cafeteria"}`), 200)
	respuesta := check("POST", "/api/pacientes/P-2/ataque", []byte(`{"ubicacion":"pasillo"}`), 200)
	var aviso map[string]any
	if err := json.Unmarshal(respuesta, &aviso); err != nil || aviso["camaAsignada"] != false {
		t.Fatal("el cuarto cupo inexistente debe quedar en pasillo")
	}
	check("GET", "/api/medicos/D-1/episodios", nil, 200)
	check("POST", "/api/pacientes/P-1/despertar", nil, 200)
	check("GET", "/api/pacientes/pasillo", nil, 200)
	check("GET", "/api/medicos/inexistente/episodios", nil, 404)
	check("POST", "/api/pacientes/P-3/despertar", nil, 404)
	check("PUT", "/api/pacientes", nil, 405)
	check("POST", "/api/pacientes/P-1/ataque", []byte(`{}`), 400)
	check("GET", "/api/inexistente", nil, 404)
}
