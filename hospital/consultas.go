package hospital

// Las consultas devuelven datos, no imprimen: la impresion ocurre en main.
// Viven aparte de hospital.go, como el consultas.go de goexample.
//
// Toman el candado porque pueden llamarse mientras otras goroutines despachan
// ataques (bono de concurrencia o API). Ninguna se llama desde un metodo que ya
// tenga el candado, porque sync.Mutex no es reentrante.

// PacientesEnPasillo devuelve los pacientes dormidos en un pasillo, para poder
// enviar a un camillero. Consulta 5.1 del enunciado.
func (h *Hospital) PacientesEnPasillo() []*Paciente {
	h.mu.Lock()
	defer h.mu.Unlock()

	enPasillo := []*Paciente{}
	for _, p := range h.pacientes {
		if p.Estado() == DormidoEnPasillo {
			enPasillo = append(enPasillo, p)
		}
	}
	return enPasillo
}

// ReporteSeveros cruza los pacientes de nivel Severo con la cantidad de
// episodios que tuvo cada uno. Consulta 5.4 del enunciado.
//
// Devuelve map[*Paciente]int y no un slice de pares porque la clave natural es
// la identidad del paciente: la consulta por paciente es directa, sin recorrer.
// La clave es *Paciente y no Paciente para no copiar el struct y que la clave
// siga siendo el mismo paciente. Los Severo sin episodios aparecen con 0.
func (h *Hospital) ReporteSeveros() map[*Paciente]int {
	h.mu.Lock()
	defer h.mu.Unlock()

	reporte := make(map[*Paciente]int)
	for _, p := range h.pacientes {
		if p.Nivel() == Severo {
			reporte[p] = 0
		}
	}
	for _, reg := range h.historial {
		if _, esSevero := reporte[reg.paciente]; esSevero {
			reporte[reg.paciente]++
		}
	}
	return reporte
}
