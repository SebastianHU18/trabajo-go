// Command costenos-narcolepsia arma y corre el escenario de demostracion del
// enunciado (seccion 6) e imprime las cuatro consultas (seccion 5). Aqui solo
// va el escenario y la impresion: toda la logica de negocio vive en el paquete
// hospital.
package main

import (
	"fmt"
	"sort"

	"costenos-narcolepsia/hospital"
)

func main() {
	// --- Montaje: 1 hospital, 2 medicos, 1 camillero, 3 habitaciones, 5 pacientes.
	h := hospital.NuevoHospital("Hospital de los Costenos con Narcolepsia")

	karen := hospital.NuevoMedico("D-01", "Karen Ospina", 38, "Neurologia")
	mauricio := hospital.NuevoMedico("D-02", "Mauricio Perez", 45, "Medicina interna")
	wilson := hospital.NuevoCamillero("C-01", "Wilson Castro", 30, "manana")

	_ = h.ContratarMedico(karen)
	_ = h.ContratarMedico(mauricio)
	_ = h.ContratarCamillero(wilson)

	for _, hab := range []*hospital.Habitacion{
		hospital.NuevaHabitacion(101, 1),
		hospital.NuevaHabitacion(102, 1),
		hospital.NuevaHabitacion(103, 1),
	} {
		_ = h.RegistrarHabitacion(hab)
	}

	yeimy := hospital.NuevoPaciente("P-001", "Yeimy Padilla", 29, hospital.Severo)
	wilfrido := hospital.NuevoPaciente("P-002", "Wilfrido Berrio", 41, hospital.Severo)
	karla := hospital.NuevoPaciente("P-003", "Karla Fontalvo", 33, hospital.Moderado)
	dagoberto := hospital.NuevoPaciente("P-004", "Dagoberto Movil", 52, hospital.Leve)
	esther := hospital.NuevoPaciente("P-005", "Esthercita Polo", 26, hospital.Severo)

	for _, p := range []*hospital.Paciente{yeimy, wilfrido, karla, dagoberto, esther} {
		_ = h.AdmitirPaciente(p)
	}

	fmt.Printf("%s\n", h.Nombre())
	fmt.Println("Personal: 2 medicos y 1 camillero | Habitaciones: 3 | Pacientes admitidos: 5")

	// --- [5.3] Ataques de sueno y asignacion de camas.
	fmt.Println("\n=== [5.3] Ataques de sueno y asignacion de camas ===")
	type ataque struct {
		p         *hospital.Paciente
		ubicacion string
	}
	for _, a := range []ataque{
		{yeimy, "cafeteria"},
		{wilfrido, "sala de espera"},
		{karla, "fila de radiologia"},
		{dagoberto, "pasillo del segundo piso"},
	} {
		reg, err := h.DespacharAtencion(a.p, a.ubicacion)
		if err != nil {
			fmt.Printf("  %s %-16s ataque en %-24q -> ERROR: %v (se queda en el pasillo)\n",
				a.p.ID(), a.p.Nombre(), a.ubicacion, err)
			continue
		}
		fmt.Printf("  %s\n", reg.Resumen())
	}

	// --- Un paciente despierta, libera su cama y el del pasillo la recibe.
	fmt.Println("\n--- Un paciente despierta y cede su habitacion ---")
	if err := h.DespertarPaciente(yeimy); err != nil {
		fmt.Printf("  no se pudo despertar a %s: %v\n", yeimy.ID(), err)
	} else {
		fmt.Printf("  %s %s despierta y libera su habitacion\n", yeimy.ID(), yeimy.Nombre())
	}
	if hab, err := h.AsignarHabitacion(dagoberto); err != nil {
		fmt.Printf("  %s sigue sin cama: %v\n", dagoberto.ID(), err)
	} else {
		fmt.Printf("  %s %s pasa del pasillo a la habitacion %d\n",
			dagoberto.ID(), dagoberto.Nombre(), hab.Numero())
	}

	// --- [5.1] Pacientes dormidos en el pasillo.
	fmt.Println("\n=== [5.1] Pacientes dormidos en el pasillo ===")
	enPasillo := h.PacientesEnPasillo()
	if len(enPasillo) == 0 {
		fmt.Println("  (ninguno: todos quedaron ubicados)")
	}
	for _, p := range enPasillo {
		fmt.Printf("  %s %s (%s)\n", p.ID(), p.Nombre(), p.Ubicacion())
	}

	// --- [5.2] Episodios atendidos por un medico.
	fmt.Printf("\n=== [5.2] Episodios atendidos por Dra. %s ===\n", karen.Nombre())
	episodios := karen.MisEpisodios()
	if len(episodios) == 0 {
		fmt.Println("  (sin episodios)")
	}
	for _, e := range episodios {
		fmt.Printf("  %s\n", e.Resumen())
	}

	// --- [5.4] Reporte de severidad.
	fmt.Println("\n=== [5.4] Reporte de severidad (pacientes Severo) ===")
	imprimirReporteSeveros(h.ReporteSeveros())
}

// imprimirReporteSeveros imprime el reporte ordenado por id de paciente. El
// reporte es un map, cuyo recorrido en Go es aleatorio; se ordena aqui, en la
// capa de impresion, para que la salida sea estable entre corridas.
func imprimirReporteSeveros(reporte map[*hospital.Paciente]int) {
	pacientes := make([]*hospital.Paciente, 0, len(reporte))
	for p := range reporte {
		pacientes = append(pacientes, p)
	}
	sort.Slice(pacientes, func(i, j int) bool {
		return pacientes[i].ID() < pacientes[j].ID()
	})
	for _, p := range pacientes {
		fmt.Printf("  %s %-16s %d episodios hoy\n", p.ID(), p.Nombre(), reporte[p])
	}
}
