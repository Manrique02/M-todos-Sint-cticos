package main

import (
	"fmt"
	"unicode"
)

type PredictivoRecursivo struct {
	input        string
	position     int
	currentToken string
	pasos        []string
	error        string
	errorPos     int
}

func NewPredictivoRecursivo(input string) *PredictivoRecursivo {
	p := &PredictivoRecursivo{
		input:    input,
		position: 0,
		pasos:    make([]string, 0),
	}
	p.nextToken()
	return p
}

func (p *PredictivoRecursivo) nextToken() {
	// Saltar espacios en blanco
	for p.position < len(p.input) && unicode.IsSpace(rune(p.input[p.position])) {
		p.position++
	}

	if p.position >= len(p.input) {
		p.currentToken = "EOF"
		return
	}

	// Reconocer números
	if unicode.IsDigit(rune(p.input[p.position])) {
		for p.position < len(p.input) && unicode.IsDigit(rune(p.input[p.position])) {
			p.position++
		}
		p.currentToken = "num"
		return
	}

	// Reconocer operadores y paréntesis
	switch p.input[p.position] {
	case '+', '*', '(', ')':
		p.currentToken = string(p.input[p.position])
		p.position++
	default:
		p.currentToken = "ERROR"
	}
}

func (p *PredictivoRecursivo) match(expected string) bool {
	if p.currentToken == expected {
		p.nextToken()
		return true
	}
	return false
}

// E → T E'
func (p *PredictivoRecursivo) E() bool {
	p.pasos = append(p.pasos, "Aplicando regla E → T E' (Entra a T para factor)")
	fmt.Println("E → T E'")
	return p.T() && p.El()
}

// E' → + T E' | ε
func (p *PredictivoRecursivo) El() bool {
	if p.currentToken == "+" {
		p.pasos = append(p.pasos, "Encontró '+': aplicando E' → + T E'")
		fmt.Println("E' → + T E'")
		return p.match("+") && p.T() && p.El()
	}
	p.pasos = append(p.pasos, "No hay '+' visible, cerrando con producción vacía E' → ε")
	fmt.Println("E' → ε")
	return true // ε producción
}

// T → F T'
func (p *PredictivoRecursivo) T() bool {
	p.pasos = append(p.pasos, "Aplicando T → F T' (Entra a F para término)")
	fmt.Println("T → F T'")
	return p.F() && p.Tl()
}

// T' → * F T' | ε
func (p *PredictivoRecursivo) Tl() bool {
	if p.currentToken == "*" {
		p.pasos = append(p.pasos, "Encontró '*': aplicando T' → * F T'")
		fmt.Println("T' → * F T'")
		return p.match("*") && p.F() && p.Tl()
	}
	p.pasos = append(p.pasos, "No hay '*' visible, cerrando con producción vacía T' → ε")
	fmt.Println("T' → ε")
	return true // ε producción
}

// F → ( E ) | num
func (p *PredictivoRecursivo) F() bool {
	if p.currentToken == "(" {
		p.pasos = append(p.pasos, "Encontró '(': aplicando F → ( E ) para subexpresión")
		fmt.Println("F → ( E )")
		return p.match("(") && p.E() && p.match(")")
	} else if p.currentToken == "num" {
		p.pasos = append(p.pasos, "Encontró número: aplicando F → num (factor base)")
		fmt.Println("F → num")
		return p.match("num")
	}
	p.error = "Caracter invalido encontrado: " + p.currentToken + " (esperaba número o '(')"
	p.errorPos = p.position
	return false
}

func (p *PredictivoRecursivo) Parse() bool {
	result := p.E() && p.currentToken == "EOF"
	if result {
		p.pasos = append(p.pasos, "Alcanzado fin de entrada (EOF) - Análisis válido")
		fmt.Println("✅ Análisis exitoso!")
	} else {
		if p.error == "" {
			if p.currentToken != "EOF" {
				p.error = "Entrada no completamente consumida. Token pendiente: " + p.currentToken
			}
		}
		p.pasos = append(p.pasos, "Análisis fallido: "+p.error)
		fmt.Println("❌ Error de sintaxis")
	}
	return result
}
