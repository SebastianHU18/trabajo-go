# Hospital de los Costeños con Narcolepsia

Sistema central de gestión del hospital, escrito **estrictamente en Go (Golang) 1.21+** usando solo la biblioteca estándar. El programa arma el escenario de demostración del enunciado e imprime las cuatro consultas requeridas.

## Integrantes

| Integrante | Bloque |
| --- | --- |
| Tomás (completar nombre) | Personas, estados, pacientes, consultas 5.1 y 5.4, bono de concurrencia |
| Samuel (completar nombre) | Habitaciones, episodios, staff polimórfico, interfaz, bono de interfaz |
| Sebastián (completar nombre) | Hospital, escenario, tests, documentación e integración |

> **Nota sobre el grupo:** el enunciado indica `Groups: 1 student`, pero este grupo trabaja con una excepción autorizada por el profesor para ser de tres integrantes. Cada integrante se defiende por separado con su factor individual de nota.

## Cómo correr

Requiere Go 1.21 o superior instalado.

```bash
go run .
```

## Cómo probar

```bash
go test ./...
```

Formato y análisis estático (no deben reportar nada):

```bash
gofmt -l .
go vet ./...
```

## Tabla de nombres (enunciado → proyecto)

El enunciado permite trabajar en español siempre que sea consistente. Esta tabla deja ubicar cada método del enunciado en el código. **Los identificadores no llevan tildes ni eñes** (por eso `SufrirAtaqueDeSueno`); las tildes sí van en comentarios y textos impresos.

| Enunciado (inglés) | Proyecto (español) |
| --- | --- |
| `Person` | `Persona` |
| `Patient` | `Paciente` |
| `Doctor` | `Medico` |
| `Orderly` | `Camillero` |
| `Room` | `Habitacion` |
| `EpisodeRecord` | `RegistroEpisodio` |
| `Attender` | `Atendedor` |
| `PatientState` | `EstadoPaciente` |
| `NarcolepsyLevel` | `NivelNarcolepsia` |
| room state | `EstadoHabitacion` |
| `Awake` / `AsleepInHallway` / `AsleepInBed` | `Despierto` / `DormidoEnPasillo` / `DormidoEnCama` |
| `Mild` / `Moderate` / `Severe` | `Leve` / `Moderado` / `Severo` |
| `Available` / `Occupied` | `Disponible` / `Ocupada` |
| `AdmitPatient` | `AdmitirPaciente` |
| `HireDoctor` | `ContratarMedico` |
| `AssignRoom` | `AsignarHabitacion` |
| `RegisterEpisode` | `RegistrarEpisodio` |
| `SufferSleepAttack` | `SufrirAtaqueDeSueno` |
| `WakeUp` | `Despertar` |
| `State` | `Estado` |
| `DiagnosePatient` | `DiagnosticarPaciente` |
| `AttendSleepEmergency` | `AtenderEmergenciaSueno` |
| `MyEpisodes` | `MisEpisodios` |
| `IsAvailable` | `EstaDisponible` |
| `Occupy` / `Release` | `Ocupar` / `Liberar` |
| `Summary` | `Resumen` |
| `PatientsInHallway` | `PacientesEnPasillo` |
| `SevereReport` | `ReporteSeveros` |
| `Attend` (interfaz) | `Atender` |

## Estructura y responsabilidades

Un archivo por entidad: ningún archivo tiene dos dueños, para que los merges no peleen.

| Archivo | Responsabilidad | Dueño |
| --- | --- | --- |
| `hospital/persona.go` | `Persona` (id, nombre, edad) y sus lecturas. Es el struct que se embute. | Tomás |
| `hospital/estados.go` | `EstadoPaciente` y `NivelNarcolepsia` como constantes tipadas, con `String()`. | Tomás |
| `hospital/paciente.go` | `Paciente` con `Persona` embutida. `SufrirAtaqueDeSueno`, `Despertar`, `Estado`. | Tomás |
| `hospital/consultas.go` | `PacientesEnPasillo` (5.1) y `ReporteSeveros` (5.4). | Tomás |
| `hospital/habitacion.go` | `Habitacion` y `EstadoHabitacion`. `EstaDisponible`, `Ocupar`, `Liberar`. | Samuel |
| `hospital/episodio.go` | `RegistroEpisodio` con `time.Time` y `Resumen()`. | Samuel |
| `hospital/atendedor.go` | La interfaz `Atendedor`. | Samuel |
| `hospital/medico.go` | `Medico` con `Persona` embutida. `DiagnosticarPaciente`, `MisEpisodios`, `Atender`. | Samuel |
| `hospital/camillero.go` | `Camillero`, segundo tipo que cumple `Atendedor`. | Samuel |
| `hospital/errores.go` | Los errores previsibles del paquete. | Sebastián |
| `hospital/hospital.go` | `Hospital`: admitir, contratar, asignar, despachar y el `[]Atendedor`. | Sebastián |
| `main.go` | El escenario de la sección 6 y la impresión de las cuatro consultas. | Sebastián |
| `hospital/hospital_test.go` | Los tests obligatorios (asignación con cupo, sin cupo, y consultas). | Sebastián |

## Decisiones de diseño

- **Composición en vez de herencia.** `Persona` se embute por valor en `Paciente` y en `Medico`. Así ambos reusan `ID()`, `Nombre()` y `Edad()` sin duplicar campos, y con eso ya cumplen la interfaz `Atendedor` sin reescribir esos métodos. Se embute por **valor** (no por puntero) para que cada entidad tenga su propia identidad y no exista el riesgo de una `Persona` nula que haga estallar `Nombre()`.
- **Polimorfismo real.** El hospital guarda su personal en un `[]Atendedor` y `DespacharAtencion` recorre y llama a `Atender` **sin preguntar el tipo concreto** (no hay `switch` por tipo ni `.(*Medico)`). `ContratarPersonal(a Atendedor)` permite enganchar en vivo un tipo nuevo (por ejemplo una enfermera) sin tocar el resto.
- **Encapsulamiento por paquete.** Todos los campos van en minúscula y se leen o modifican por métodos. Un `main` fuera del paquete (ver `local/`) no puede tocar `p.estado`: no compila.
- **Estados como constantes tipadas.** `EstadoPaciente`, `NivelNarcolepsia` y `EstadoHabitacion` son tipos con constantes y `String()`, no cadenas sueltas, para que el compilador atrape errores de escritura y se impriman legibles.
- **Errores, no `panic`.** Los métodos que pueden fallar devuelven `error` como último valor. Los errores previsibles son variables del paquete (`errores.go`) para compararlos con `errors.Is`.
- **`map[*Paciente]int` en `ReporteSeveros`.** La clave es la identidad del paciente (el puntero), así la consulta por paciente es directa; con un slice de pares habría que recorrer. Es `*Paciente` y no `Paciente` para no copiar el struct y conservar la identidad como clave.
- **Receptor puntero en todo `Hospital`.** El struct lleva un `sync.Mutex` (para el bono de concurrencia), y copiar un mutex lo rompe; `go vet` lo reportaría. El método público `AsignarHabitacion` toma el candado una vez y delega en el privado `asignarHabitacionSinCandado`, porque `sync.Mutex` no es reentrante.

## Bonos

- **Concurrencia (+5%, Tomás):** `hospital/simulacion.go` lanza una goroutine por paciente y protege el estado con el `sync.Mutex` del hospital. Debe correr limpio con `go run -race .`.
  > En Windows, el detector de carreras (`-race`) necesita un compilador de C (gcc) en el PATH; por ejemplo el de MSYS2/MinGW. Sin él, `-race` no compila (el resto del proyecto sí corre normal).
- **Interfaz (+10%, Samuel):** API REST en `api/` y frontend en `web/` con Gin, declarado en `go.mod` únicamente para la interfaz. El paquete `hospital` no importa `net/http` ni Gin.
