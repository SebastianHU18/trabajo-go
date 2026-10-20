# Uso de IA

Para este trabajo nos apoyamos mas que todo en Chat y en Claude sobre todo para resolver dudas de Go y claramente para que nos ayudara a resolver el proyecto,
porque la forma de Go es distinta y nos costo un poco. La usamos
para entender cómo reemplazar la herencia con composición e interfaces, para revisar que el
código fuera idiomático y para aclarar errores de compilación. El diseño, las decisiones y las
pruebas las fuimos mirando entre nosotros; la IA sirvió de apoyo, no reemplazó el trabajo.

## Prompts que nos sirvieron

- "¿Cómo hago en Go que un médico y un camillero se puedan tratar igual sin usar herencia?"
  De ahí salió la idea de la interfaz `Atendedor` y del `[]Atendedor` del hospital esta parte no se entendio casi pero fue importante para poder seguir con el desarrollo del trabajo.
- "En Go, ¿cuándo uso receptor por valor y cuándo por puntero?" Esto nos evitó varios errores
  donde los cambios no se guardaban y se entendio la diferencia de cuando usar uno o otro.

## Casos donde el código que generó estuvo mal

Estos son los dos puntos que más se nos complicaron y donde tuvimos que corregir lo que salía:

1. En un método que cambia el estado del paciente, la IA usó receptor por valor. El resultado
   era que el cambio no se reflejaba porque estaba trabajando sobre una copia. Lo corregimos
   pasando el método a receptor por puntero.

2. Cuando no había habitación libre, la primera versión usaba `panic`. El enunciado pide
   devolver un `error`, así que lo cambiamos por un error del paquete (`ErrSinHabitacion`) para
   poder compararlo con `errors.Is` y que el programa no se cayera.

## Qué aprendimos

- **Sebastián:** se entendio cómo Go maneja la privacidad con el paquete (los campos en minúscula) y
  por qué, al usar un mutex, los métodos tienen que ir con puntero.
- **Tomás:** me quedó claro para qué sirve guardar algo por valor o por puntero, y por qué los
  estados se hacen con constantes tipadas y no con textos sueltos.
- **Samuel:** aprendí cómo una interfaz deja usar distintos tipos de la misma forma sin tener
  que usar herencia, que es lo que se creia que se debia de hacer al principio.
