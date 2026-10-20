// Command web inicia una demostracion interactiva del hospital.
// Ejecutar desde la raiz del repositorio con: go run ./cmd/web
package main

import (
	"log"
	"net/http"
	"os"

	"costenos-narcolepsia/api"
	"costenos-narcolepsia/hospital"
)

func main() {
	h := hospital.NuevoHospital("Hospital de los Costenos con Narcolepsia")
	for _, m := range []*hospital.Medico{
		hospital.NuevoMedico("D-01", "Karen Ospina", 38, "Neurologia"),
		hospital.NuevoMedico("D-02", "Mauricio Perez", 45, "Medicina interna"),
	} {
		if err := h.ContratarMedico(m); err != nil {
			log.Fatal(err)
		}
	}
	if err := h.ContratarCamillero(hospital.NuevoCamillero("C-01", "Wilson Castro", 30, "manana")); err != nil {
		log.Fatal(err)
	}
	for _, numero := range []int{101, 102, 103} {
		if err := h.RegistrarHabitacion(hospital.NuevaHabitacion(numero, 1)); err != nil {
			log.Fatal(err)
		}
	}
	for _, p := range []*hospital.Paciente{
		hospital.NuevoPaciente("P-001", "Goma jugador", 29, hospital.Severo),
		hospital.NuevoPaciente("P-002", "Urrego galleta", 41, hospital.Severo),
		hospital.NuevoPaciente("P-003", "Chanda chambeador", 33, hospital.Moderado),
		hospital.NuevoPaciente("P-004", "Portu goleador", 52, hospital.Leve),
		hospital.NuevoPaciente("P-005", "Manuela paisa", 26, hospital.Severo),
	} {
		if err := h.AdmitirPaciente(p); err != nil {
			log.Fatal(err)
		}
	}
	puerto := os.Getenv("HOSPITAL_PORT")
	if puerto == "" {
		puerto = "8080"
	}
	log.Printf("Interfaz disponible en http://localhost:%s", puerto)
	log.Fatal(http.ListenAndServe("127.0.0.1:"+puerto, api.NuevoServidor(h)))
}
