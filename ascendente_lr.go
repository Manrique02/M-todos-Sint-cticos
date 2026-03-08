package main

import (
	"fmt"
	"strings"
)

type ItemLR struct {
	produccion string
	posicion   int
	lookahead  string
}

type EstadoLR struct {
	items []ItemLR
	id    int
}

type AnalizadorLR struct {
	estados      []*EstadoLR
	tablaAccion  map[int]map[string]string
	tablaGoto    map[int]map[string]int
	producciones []string
}

func NewAnalizadorLR() *AnalizadorLR {
	a := &AnalizadorLR{
		estados:     make([]*EstadoLR, 0),
		tablaAccion: make(map[int]map[string]string),
		tablaGoto:   make(map[int]map[string]int),
		producciones: []string{
			"E' → E",
			"E → E + T",
			"E → T",
			"T → T * F",
			"T → F",
			"F → ( E )",
			"F → id",
		},
	}

	a.construirTablas()
	return a
}

func (a *AnalizadorLR) construirTablas() {
	// Construir estados LR(0) (simplificado para el ejemplo)
	a.estados = append(a.estados, &EstadoLR{
		items: []ItemLR{
			{produccion: "E' → • E", posicion: 0},
			{produccion: "E → • E + T", posicion: 0},
			{produccion: "E → • T", posicion: 0},
			{produccion: "T → • T * F", posicion: 0},
			{produccion: "T → • F", posicion: 0},
			{produccion: "F → • ( E )", posicion: 0},
			{produccion: "F → • id", posicion: 0},
		},
		id: 0,
	})

	// Configurar tabla de acción (simplificada)
	a.tablaAccion[0] = map[string]string{
		"id": "s5",
		"(":  "s4",
	}
	a.tablaAccion[1] = map[string]string{
		"+":   "s6",
		"EOF": "acc",
	}
	a.tablaAccion[2] = map[string]string{
		"+":   "r2",
		"*":   "s7",
		")":   "r2",
		"EOF": "r2",
	}
	a.tablaAccion[3] = map[string]string{
		"+":   "r4",
		"*":   "r4",
		")":   "r4",
		"EOF": "r4",
	}
	a.tablaAccion[4] = map[string]string{
		"id": "s5",
		"(":  "s4",
	}
	a.tablaAccion[5] = map[string]string{
		"+":   "r6",
		"*":   "r6",
		")":   "r6",
		"EOF": "r6",
	}
	a.tablaAccion[6] = map[string]string{
		"id": "s5",
		"(":  "s4",
	}
	a.tablaAccion[7] = map[string]string{
		"id": "s5",
		"(":  "s4",
	}
	a.tablaAccion[8] = map[string]string{
		"+": "s6",
		")": "s11",
	}
	a.tablaAccion[9] = map[string]string{
		"+": "r1",
		"*": "s7",
		")": "r1",
	}
	a.tablaAccion[10] = map[string]string{
		"+": "r3",
		"*": "r3",
		")": "r3",
	}
	a.tablaAccion[11] = map[string]string{
		"+": "r5",
		"*": "r5",
		")": "r5",
	}

	// Configurar tabla GOTO
	a.tablaGoto[0] = map[string]int{"E": 1, "T": 2, "F": 3}
	a.tablaGoto[4] = map[string]int{"E": 8, "T": 2, "F": 3}
	a.tablaGoto[6] = map[string]int{"T": 9, "F": 3}
	a.tablaGoto[7] = map[string]int{"F": 10}
}

type PilaLR struct {
	estados  []int
	simbolos []string
	valores  []int // Para valores semánticos (opcional)
}

func (p *PilaLR) Push(estado int, simbolo string, valor int) {
	p.estados = append(p.estados, estado)
	p.simbolos = append(p.simbolos, simbolo)
	p.valores = append(p.valores, valor)
}

func (p *PilaLR) Pop() (int, string, int) {
	if len(p.estados) == 0 {
		return 0, "", 0
	}
	estado := p.estados[len(p.estados)-1]
	simbolo := p.simbolos[len(p.simbolos)-1]
	valor := p.valores[len(p.valores)-1]

	p.estados = p.estados[:len(p.estados)-1]
	p.simbolos = p.simbolos[:len(p.simbolos)-1]
	p.valores = p.valores[:len(p.valores)-1]

	return estado, simbolo, valor
}

func (p *PilaLR) Top() int {
	if len(p.estados) == 0 {
		return 0
	}
	return p.estados[len(p.estados)-1]
}

func (a *AnalizadorLR) tokenizar(input string) []string {
	tokens := make([]string, 0)
	input = strings.ReplaceAll(input, " ", "")

	for i := 0; i < len(input); i++ {
		c := string(input[i])
		if c >= "0" && c <= "9" {
			tokens = append(tokens, "id")
		} else {
			tokens = append(tokens, c)
		}
	}
	tokens = append(tokens, "EOF")
	return tokens
}

type ArbolSintactico struct {
	tipo  string
	valor string
	hijos []*ArbolSintactico
}

func NewArbol(tipo string, valor string) *ArbolSintactico {
	return &ArbolSintactico{
		tipo:  tipo,
		valor: valor,
		hijos: make([]*ArbolSintactico, 0),
	}
}

func (a *ArbolSintactico) AddHijo(hijo *ArbolSintactico) {
	a.hijos = append(a.hijos, hijo)
}

func (a *ArbolSintactico) Print(nivel int) {
	indent := strings.Repeat("  ", nivel)
	if len(a.hijos) == 0 {
		fmt.Printf("%s%s: %s\n", indent, a.tipo, a.valor)
	} else {
		fmt.Printf("%s%s\n", indent, a.tipo)
		for _, hijo := range a.hijos {
			hijo.Print(nivel + 1)
		}
	}
}

func (a *AnalizadorLR) Parse(input string) *ArbolSintactico {
	fmt.Println("\n=== ANALIZADOR ASCENDENTE LR ===")
	fmt.Printf("Analizando: %s\n\n", input)

	tokens := a.tokenizar(input)
	pila := &PilaLR{}
	pila.Push(0, "", 0)
	pos := 0

	fmt.Println("Paso\tPila\t\tEntrada\t\tAcción")
	fmt.Println("----\t----\t\t-------\t\t------")

	paso := 1

	for {
		estado := pila.Top()
		token := tokens[pos]

		// Mostrar estado actual
		fmt.Printf("%d\t%v\t\t%s", paso, pila.estados, token)

		// Consultar acción
		accion, exists := a.tablaAccion[estado][token]
		if !exists {
			fmt.Println("\tERROR: Acción no definida")
			return nil
		}

		fmt.Printf("\t\t%s", accion)

		if accion[0] == 's' { // Shift
			var nuevoEstado int
			fmt.Sscanf(accion[1:], "%d", &nuevoEstado)

			pila.Push(nuevoEstado, token, len(pila.valores))
			pos++
			fmt.Printf(" (desplazar %s)", token)

		} else if accion[0] == 'r' { // Reduce
			var numProd int
			fmt.Sscanf(accion[1:], "%d", &numProd)

			// Aplicar reducción según producción
			prod := a.producciones[numProd]
			partes := strings.Split(prod, " → ")
			cabeza := partes[0]
			cuerpo := partes[1]

			// Calcular longitud de la reducción
			longitud := len(strings.Split(cuerpo, " "))
			if cuerpo == "ε" {
				longitud = 0
			}

			// Construir nodo para la producción
			// nodo := NewArbol(cabeza, "") // Simplificado para este ejemplo

			// Construir árbol semántico
			fmt.Printf(" (reducir %s)", prod)

			if longitud > 0 {
				// Pop y construir hijos (de derecha a izquierda)
				// hijos := make([]*ArbolSintactico, longitud) // Simplificado
				for i := longitud - 1; i >= 0; i-- {
					pila.Pop() // Simplificado: solo descartamos los valores
					// Nota: En una implementación real, aquí recuperaríamos el nodo del árbol
					// y lo añadiríamos como hijo
				}
			}

			// Push nuevo estado
			nuevoEstado := a.tablaGoto[pila.Top()][cabeza]
			pila.Push(nuevoEstado, cabeza, 0)

		} else if accion == "acc" { // Accept
			fmt.Println("\t✅ ACEPTADO")
			fmt.Println("\nAnálisis completado exitosamente!")

			// Construir árbol sintáctico final
			return NewArbol("Programa", input)
		}

		fmt.Println()
		paso++
	}
}

func (a *AnalizadorLR) evaluar(arbol *ArbolSintactico) int {
	// Evaluador semántico simple
	if arbol.tipo == "id" {
		val := 0
		fmt.Sscanf(arbol.valor, "%d", &val)
		return val
	}

	switch arbol.tipo {
	case "E":
		if len(arbol.hijos) == 3 && arbol.hijos[1].valor == "+" {
			return a.evaluar(arbol.hijos[0]) + a.evaluar(arbol.hijos[2])
		}
		if len(arbol.hijos) == 1 {
			return a.evaluar(arbol.hijos[0])
		}
	case "T":
		if len(arbol.hijos) == 3 && arbol.hijos[1].valor == "*" {
			return a.evaluar(arbol.hijos[0]) * a.evaluar(arbol.hijos[2])
		}
		if len(arbol.hijos) == 1 {
			return a.evaluar(arbol.hijos[0])
		}
	case "F":
		if len(arbol.hijos) == 3 && arbol.hijos[0].valor == "(" {
			return a.evaluar(arbol.hijos[1])
		}
		if len(arbol.hijos) == 1 {
			return a.evaluar(arbol.hijos[0])
		}
	}

	return 0
}
