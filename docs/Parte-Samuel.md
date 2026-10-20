# Parte de Samuel — Hospital de los Costeños con Narcolepsia

## Componentes

- `hospital/habitacion.go`: gestión de capacidad, ocupación y liberación de habitaciones.
- `hospital/episodio.go`: registro encapsulado de episodios, ID, instante, atendedor y habitación.
- `hospital/atendedor.go`: interfaz polimórfica consumida por el hospital.
- `hospital/medico.go`: composición de Persona; diagnóstico, historial personal y atención.
- `hospital/camillero.go`: composición de Persona; atención y traslado sin diagnóstico.
- `hospital/habitacion_test.go`: pruebas de habitaciones, personal, episodios y composición.
- `hospital/vistas_samuel.go`: proyecciones de lectura y operaciones por ID necesarias para la interfaz web.
- `api/`: adaptador HTTP sin lógica del dominio ni librerías externas.
- `web/index.html`: panel HTML/CSS/JavaScript.
- `cmd/web/main.go`: arranque aislado del bono gráfico.

## Ejecución

Desde la raíz del proyecto y con Go 1.21 o superior:

```bash
go run .
go test ./...
go vet ./...
go run ./cmd/web
```

Abrir `http://localhost:8080` para usar el panel web. El puerto escucha únicamente en 127.0.0.1. Si 8080 está ocupado, configura `HOSPITAL_PORT` antes de iniciar (por ejemplo, `$env:HOSPITAL_PORT=8189` en PowerShell).

## Rutas HTTP

| Método | Ruta | Función |
| --- | --- | --- |
| GET | `/health` | Comprobación del servidor |
| GET | `/api/pacientes` | Pacientes y estados |
| GET | `/api/pacientes/pasillo` | Consulta 5.1 |
| GET | `/api/medicos` | Médicos disponibles |
| GET | `/api/medicos/{id}/episodios` | Consulta 5.2 |
| GET | `/api/habitaciones` | Disponibilidad de camas |
| GET | `/api/reportes/severos` | Consulta 5.4 |
| POST | `/api/pacientes/{id}/ataque` | Ataque; cuerpo JSON `{"ubicacion":"cafeteria"}` |
| POST | `/api/pacientes/{id}/despertar` | Despierta y libera habitación |

## Reglas de integración

- `Hospital.DespacharAtencion` (Sebastián) asigna primero habitación, luego llama a `Atendedor.Atender`; el registro captura esa habitación (o `nil`).
- `Hospital.DespertarPaciente` libera las camas y cambia estado; `Habitacion.Liberar` por sí solo no toca estados de `Paciente`.
- Los campos del dominio permanecen privados. `vistas_samuel.go` devuelve proyecciones JSON de solo lectura, no punteros internos.
- No se modifica `go.mod`: la API utiliza exclusivamente `net/http` y `encoding/json` estándar, alternativa permitida para el bono.
- El estado y las habitaciones se leen bajo `Hospital.mu`, permitiendo a la API coexistir con los métodos concurrentes del hospital.
- La defensa exige explicar también las partes de Tomás y Sebastián y anotar en `AI_USAGE.md` las interacciones reales con IA; este archivo no sustituye ese reporte.

## Aspectos que debe revisar el equipo

- El `ReporteSeveros` de Tomás debe aplicar correctamente el filtro de fecha exigido por el enunciado. En la rama `sebastian` observada, sus tests crean registros sin fecha: deben alinear ese test con la consulta del día.
- El `Atender` del camillero registra su intervención sin diagnosticar; la asignación de la habitación corresponde exclusivamente al Hospital.
- Todos deben integrar sus ramas y verificar `go run .`, `go test ./...` y `go vet ./...` sobre un clon limpio.
