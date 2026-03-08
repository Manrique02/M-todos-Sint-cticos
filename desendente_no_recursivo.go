package main

import (
	"fmt"
)

type Token struct {
	Tipo  string
	Valor string
}

type Pila struct {
	items []string
}

func (p *Pila) Push(item string) {
	p.items = append(p.items, item)
}

func (p *Pila) Pop() string {
	if len(p.items) == 0 {
		return ""
	}
	item := p.items[len(p.items)-1]
	p.items = p.items[:len(p.items)-1]
	return item
}

func (p *Pila) Top() string {
	if len(p.items) == 0 {
		return ""
	}
	return p.items[len(p.items)-1]
}

type PredictivoNoRecursivo struct {
	tablaAccion map[string]map[string]string
	tablaGoto   map[string]map[string]string
	pila        *Pila
	entrada     []Token
	posicion    int
}

func NewPredictivoNoRecursivo() *PredictivoNoRecursivo {
	p := &PredictivoNoRecursivo{
		tablaAccion: make(map[string]map[string]string),
		tablaGoto:   make(map[string]map[string]string),
		pila:        &Pila{},
		entrada:     make([]Token, 0),
		posicion:    0,
	}

	// Inicializar tabla de acción (simplificada)
	p.tablaAccion["0"] = map[string]string{
		"num": "s5",
		"(":   "s4",
		"EOF": "",
	}
	p.tablaAccion["1"] = map[string]string{
		"+":   "s6",
		"*":   "",
		")":   "",
		"EOF": "acc",
	}
	p.tablaAccion["2"] = map[string]string{
		"+":   "r2",
		"*":   "s7",
		")":   "r2",
		"EOF": "r2",
	}
	p.tablaAccion["3"] = map[string]string{
		"+":   "r4",
		"*":   "r4",
		")":   "r4",
		"EOF": "r4",
	}
	p.tablaAccion["4"] = map[string]string{
		"num": "s5",
		"(":   "s4",
	}
	p.tablaAccion["5"] = map[string]string{
		"+":   "r6",
		"*":   "r6",
		")":   "r6",
		"EOF": "r6",
	}
	p.tablaAccion["6"] = map[string]string{
		"num": "s5",
		"(":   "s4",
	}
	p.tablaAccion["7"] = map[string]string{
		"num": "s5",
		"(":   "s4",
	}
	p.tablaAccion["8"] = map[string]string{
		"+": "s6",
		")": "s11",
	}
	p.tablaAccion["9"] = map[string]string{
		"+":   "r1",
		"*":   "s7",
		")":   "r1",
		"EOF": "r1",
	}
	p.tablaAccion["10"] = map[string]string{
		"+":   "r3",
		"*":   "r3",
		")":   "r3",
		"EOF": "r3",
	}
	p.tablaAccion["11"] = map[string]string{
		"+":   "r5",
		"*":   "r5",
		")":   "r5",
		"EOF": "r5",
	}

	// Inicializar tabla GOTO
	p.tablaGoto["0"] = map[string]string{"E": "1", "T": "2", "F": "3"}
	p.tablaGoto["4"] = map[string]string{"E": "8", "T": "2", "F": "3"}
	p.tablaGoto["6"] = map[string]string{"T": "9", "F": "3"}
	p.tablaGoto["7"] = map[string]string{"F": "10"}

	return p
}

func (p *PredictivoNoRecursivo) tokenizar(input string) []Token {
	tokens := make([]Token, 0)
	i := 0

	for i < len(input) {
		// Saltar espacios
		if input[i] == ' ' {
			i++
			continue
		}

		// Números
		if input[i] >= '0' && input[i] <= '9' {
			start := i
			for i < len(input) && input[i] >= '0' && input[i] <= '9' {
				i++
			}
			tokens = append(tokens, Token{"num", input[start:i]})
			continue
		}

		// Operadores y paréntesis
		switch input[i] {
		case '+', '*', '(', ')':
			tokens = append(tokens, Token{string(input[i]), string(input[i])})
		}
		i++
	}

	tokens = append(tokens, Token{"EOF", ""})
	return tokens
}

func (p *PredictivoNoRecursivo) Parse(input string) bool {
	fmt.Println("=== ANALIZADOR PREDICTIVO NO RECURSIVO (LR) ===")
	fmt.Printf("Analizando: %s\n\n", input)

	p.entrada = p.tokenizar(input)
	p.pila = &Pila{}
	p.pila.Push("0") // Estado inicial
	p.posicion = 0

	fmt.Println("Pila\t\tEntrada\t\tAcción")
	fmt.Println("----\t\t-------\t\t------")

	for {
		estado := p.pila.Top()
		token := p.entrada[p.posicion]

		// Mostrar estado actual
		fmt.Printf("%v\t\t%s\t\t", p.pila.items, token.Tipo)

		// Consultar tabla de acción
		accion, exists := p.tablaAccion[estado][token.Tipo]
		if !exists {
			fmt.Println("ERROR - Acción no definida")
			return false
		}

		if accion == "" {
			fmt.Println("ERROR - Acción vacía")
			return false
		}

		fmt.Printf("%s\n", accion)

		if accion[0] == 's' { // Shift
			// Mover a nuevo estado
			nuevoEstado := accion[1:]
			p.pila.Push(token.Tipo)
			p.pila.Push(nuevoEstado)
			p.posicion++

		} else if accion[0] == 'r' { // Reduce
			// Aplicar reducción según producción
			switch accion {
			case "r1": // E → E + T
				p.pila.Pop() // Pop 3 veces (estado, +, estado, T, estado, E)
				p.pila.Pop()
				p.pila.Pop()
				p.pila.Pop()
				p.pila.Pop()
				p.pila.Pop()
				fmt.Println("\tReduciendo: E → E + T")
				p.gotoReduce("E")

			case "r2": // E → T
				p.pila.Pop() // Pop 2 veces (estado, T)
				p.pila.Pop()
				fmt.Println("\tReduciendo: E → T")
				p.gotoReduce("E")

			case "r3": // T → T * F
				p.pila.Pop() // Pop 6 veces
				p.pila.Pop()
				p.pila.Pop()
				p.pila.Pop()
				p.pila.Pop()
				p.pila.Pop()
				fmt.Println("\tReduciendo: T → T * F")
				p.gotoReduce("T")

			case "r4": // T → F
				p.pila.Pop() // Pop 2 veces
				p.pila.Pop()
				fmt.Println("\tReduciendo: T → F")
				p.gotoReduce("T")

			case "r5": // F → ( E )
				p.pila.Pop() // Pop 6 veces
				p.pila.Pop()
				p.pila.Pop()
				p.pila.Pop()
				p.pila.Pop()
				p.pila.Pop()
				fmt.Println("\tReduciendo: F → ( E )")
				p.gotoReduce("F")

			case "r6": // F → num
				p.pila.Pop() // Pop 2 veces
				p.pila.Pop()
				fmt.Println("\tReduciendo: F → num")
				p.gotoReduce("F")
			}

		} else if accion == "acc" { // Accept
			fmt.Println("\n✅ Análisis exitoso!")
			return true
		}
	}
}

func (p *PredictivoNoRecursivo) gotoReduce(noTerminal string) {
	estado := p.pila.Top()
	nuevoEstado := p.tablaGoto[estado][noTerminal]
	p.pila.Push(noTerminal)
	p.pila.Push(nuevoEstado)
}
