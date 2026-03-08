package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
)

// Estructura para guardar pasos del análisis
type PasoAnalisis struct {
	Paso        int      `json:"paso"`
	Pila        []string `json:"pila"`
	Entrada     string   `json:"entrada"`
	Accion      string   `json:"accion"`
	Descripcion string   `json:"descripcion"`
}

type ResultadoAnalisis struct {
	Exito   bool           `json:"exito"`
	Pasos   []PasoAnalisis `json:"pasos"`
	Mensaje string         `json:"mensaje"`
	Arbol   string         `json:"arbol"`
}

func handleIndex(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprint(w, htmlContent)
}

func handleAnalizar(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(map[string]string{"error": "Método no permitido"})
		return
	}

	var req struct {
		Expresion string `json:"expresion"`
		Metodo    string `json:"metodo"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Datos inválidos"})
		return
	}

	req.Expresion = strings.TrimSpace(req.Expresion)
	if req.Expresion == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Expresión vacía"})
		return
	}

	var resultado ResultadoAnalisis

	switch req.Metodo {
	case "recursivo":
		resultado = analizarPredictivoRecursivo(req.Expresion)
	case "norecursivo":
		resultado = analizarPredictivoNoRecursivo(req.Expresion)
	case "lr":
		resultado = analizarLR(req.Expresion)
	default:
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Método desconocido"})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resultado)
}

func analizarPredictivoRecursivo(expresion string) ResultadoAnalisis {
	parser := NewPredictivoRecursivo(expresion)
	exito := parser.Parse()

	pasos := make([]PasoAnalisis, 0)

	// Agregar cada paso registrado
	for i, paso := range parser.pasos {
		pasos = append(pasos, PasoAnalisis{
			Paso:        i + 1,
			Entrada:     expresion,
			Accion:      "Paso " + fmt.Sprintf("%d", i+1),
			Descripcion: paso,
		})
	}

	resultado := ResultadoAnalisis{
		Exito: exito,
		Pasos: pasos,
	}

	if exito {
		resultado.Mensaje = "✅ Análisis exitoso!"
		detalles := "\n\n📝 EXPLICACIÓN:\n"
		detalles += "El método Predictivo Recursivo analizó la expresión \"" + expresion + "\" de forma descendente, "
		detalles += "siguiendo las reglas de la gramática de forma directa mediante funciones recursivas:\n\n"
		detalles += "1. Comenzó con la regla E (expresión)\n"
		detalles += "2. Descendió a través de T (término) y F (factor)\n"
		detalles += "3. Consumió todos los símbolos válidos\n"
		detalles += "4. Llegó exitosamente al fin de entrada (EOF)\n\n"
		detalles += "✓ El análisis fue completamente válido según la gramática LL(1)."
		resultado.Mensaje += detalles
	} else {
		resultado.Mensaje = "❌ Error de sintaxis"
		detalles := "\n\n📝 EXPLICACIÓN DEL ERROR:\n"
		detalles += "El método Predictivo Recursivo NO pudo analizar la expresión \"" + expresion + "\" porque:\n\n"
		if parser.error != "" {
			detalles += "• " + parser.error + "\n\n"
		}
		detalles += "RAZONES DEL FALLO:\n"
		detalles += "1. La expresión contiene caracteres inválidos o mal ubicados\n"
		detalles += "2. Los paréntesis no están balanceados\n"
		detalles += "3. Hay una violación de la sintaxis esperada (números y operadores +, *)\n\n"
		detalles += "El análisis descendente llegó a una rama donde no hay una producción válida para el siguiente token."
		resultado.Mensaje += detalles
	}

	return resultado
}

func analizarPredictivoNoRecursivo(expresion string) ResultadoAnalisis {
	parser := NewPredictivoNoRecursivo()
	exito := parser.Parse(expresion)

	pasos := []PasoAnalisis{
		{
			Paso:        1,
			Entrada:     expresion,
			Accion:      "Inicio",
			Descripcion: "Método Predictivo No Recursivo - Usa una tabla predictiva y pila explícita",
		},
	}

	resultado := ResultadoAnalisis{
		Exito: exito,
		Pasos: pasos,
	}

	if exito {
		resultado.Mensaje = "✅ Análisis exitoso!"
		detalles := "\n\n📝 EXPLICACIÓN:\n"
		detalles += "El método Predictivo No Recursivo analizó la expresión \"" + expresion + "\" utilizando:\n\n"
		detalles += "1. UNA TABLA PREDICTIVA: Precalculada con la gramática LL(1)\n"
		detalles += "2. UNA PILA EXPLÍCITA: Mantiene estados y símbolos el análisis\n"
		detalles += "3. ANÁLISIS DESCENDENTE: De arriba hacia abajo sin recursión\n\n"
		detalles += "VENTAJAS DE ESTE MÉTODO:\n"
		detalles += "• Evita el costo de las llamadas recursivas\n"
		detalles += "• Más eficiente en memoria\n"
		detalles += "• La tabla permite predictibilidad exacta\n\n"
		detalles += "✓ La expresión es válida según la tabla predictiva."
		resultado.Mensaje += detalles
	} else {
		resultado.Mensaje = "❌ Error de sintaxis"
		detalles := "\n\n📝 EXPLICACIÓN DEL ERROR:\n"
		detalles += "El método Predictivo No Recursivo NO validó la expresión \"" + expresion + "\" porque:\n\n"
		detalles += "RAZONES DEL FALLO:\n"
		detalles += "1. El símbolo encontrado no está en la tabla predictiva para el estado actual\n"
		detalles += "2. Hay una acción de error (sync) en la tabla\n"
		detalles += "3. La expresión viola las reglas de la gramática LL(1)\n\n"
		detalles += "DIFERENCIA CON RECURSIVO:\n"
		detalles += "• Este método detecta errores usando la tabla de forma directa\n"
		detalles += "• No necesita recorrer el árbol de decisión recursivo"
		resultado.Mensaje += detalles
	}

	return resultado
}

func analizarLR(expresion string) ResultadoAnalisis {
	parser := NewAnalizadorLR()

	// Limpiar espacios
	expresion = strings.ReplaceAll(expresion, " ", "")

	arbol := parser.Parse(expresion)

	pasos := []PasoAnalisis{
		{
			Paso:        1,
			Entrada:     expresion,
			Accion:      "Inicio",
			Descripcion: "Método Ascendente LR - Análisis bottom-up con tabla LR(0)",
		},
	}

	resultado := ResultadoAnalisis{
		Exito: arbol != nil,
		Pasos: pasos,
	}

	if arbol != nil {
		resultado.Mensaje = "✅ Análisis completado exitosamente!"
		detalles := "\n\n📝 EXPLICACIÓN:\n"
		detalles += "El método Ascendente LR analizó la expresión \"" + expresion + "\" de forma ascendente:\n\n"
		detalles += "PROCESO DE ANÁLISIS:\n"
		detalles += "1. DESPLAZAMIENTOS: Consume tokens de la entrada\n"
		detalles += "2. REDUCCIONES: Aplica reglas gramaticales de abajo hacia arriba\n"
		detalles += "3. CONSTRUCCIÓN DEL ÁRBOL: Crea la estructura sintáctica en orden inverso\n\n"
		detalles += "CARACTERÍSTICAS DEL MÉTODO LR:\n"
		detalles += "• Potencia: Maneja gramáticas más complejas que LL(1)\n"
		detalles += "• Efectividad: Detecta errores rápidamente (shift/reduce)\n"
		detalles += "• Tabla LR(0): Máquina de estados con transiciones definidas\n\n"
		detalles += "✓ El análisis Bottom-up completó exitosamente el árbol sintáctico."
		resultado.Mensaje += detalles
		resultado.Arbol = arbolToString(arbol, 0)
	} else {
		resultado.Mensaje = "❌ Error en el análisis"
		detalles := "\n\n📝 EXPLICACIÓN DEL ERROR:\n"
		detalles += "El método Ascendente LR NO pudo procesar la expresión \"" + expresion + "\" porque:\n\n"
		detalles += "RAZONES DEL FALLO:\n"
		detalles += "1. CONFLICTO SHIFT/REDUCE: No hay transición válida en la tabla LR\n"
		detalles += "2. ERROR EN ENTRADA: La expresión contiene símbolos inesperados\n"
		detalles += "3. TABLA INCOMPLETA: No hay reducción para los símbolos presentes\n\n"
		detalles += "DIFERENCIA CON MÉTODOS DESCENDENTES:\n"
		detalles += "• LR analiza consiguiendo símbolos antes de decidir qué hacer\n"
		detalles += "• Los métodos descendentes predicen sin ver suficientes símbolos\n"
		detalles += "• LR puede ser más preciso pero más complejo"
		resultado.Mensaje += detalles
	}

	return resultado
}

func arbolToString(arbol *ArbolSintactico, nivel int) string {
	indent := strings.Repeat("  ", nivel)
	var sb strings.Builder

	if len(arbol.hijos) == 0 {
		sb.WriteString(fmt.Sprintf("%s%s: %s\n", indent, arbol.tipo, arbol.valor))
	} else {
		sb.WriteString(fmt.Sprintf("%s%s\n", indent, arbol.tipo))
		for _, hijo := range arbol.hijos {
			sb.WriteString(arbolToString(hijo, nivel+1))
		}
	}

	return sb.String()
}

func main() {
	// Rutas
	http.HandleFunc("/", handleIndex)
	http.HandleFunc("/api/analizar", handleAnalizar)

	// Servir archivos estáticos si es necesario
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	fmt.Println("🚀 Servidor iniciado en http://localhost:3000")
	fmt.Println("Abre tu navegador en http://localhost:3000")

	if err := http.ListenAndServe(":3000", nil); err != nil {
		log.Fatal(err)
	}
}

const htmlContent = `
<!DOCTYPE html>
<html lang="es">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Analizador Sintáctico</title>
    <link rel="stylesheet" href="https://cdn.jsdelivr.net/npm/bootstrap-icons@1.11.3/font/bootstrap-icons.css">
    <style>
        :root {
            --color-dark: #2D261D;
            --color-brown: #A79675;
            --color-gold: #C1B18B;
            --color-cream: #DCCFAA;
            --color-light-cream: #F5F3EF;
        }

        * {
            margin: 0;
            padding: 0;
            box-sizing: border-box;
        }

        body {
            font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
            background: linear-gradient(135deg, var(--color-dark) 0%, var(--color-brown) 100%);
            min-height: 100vh;
            padding: 20px;
            color: var(--color-dark);
        }

        .container {
            max-width: 1200px;
            margin: 0 auto;
        }

        header {
            text-align: center;
            color: var(--color-cream);
            margin-bottom: 40px;
            animation: fadeIn 0.5s ease-in;
        }

        header h1 {
            font-size: 2.5em;
            margin-bottom: 10px;
            text-shadow: 2px 2px 4px rgba(0,0,0,0.3);
            display: flex;
            align-items: center;
            justify-content: center;
            gap: 15px;
        }

        header h1 i {
            font-size: 1.1em;
        }

        header p {
            font-size: 1.1em;
            opacity: 0.9;
        }

        .main-content {
            display: grid;
            grid-template-columns: 1fr 1fr;
            gap: 20px;
            margin-bottom: 20px;
        }

        .card {
            background: var(--color-light-cream);
            border-radius: 10px;
            padding: 25px;
            box-shadow: 0 10px 40px rgba(0,0,0,0.3);
            animation: slideUp 0.5s ease-out;
            border-left: 5px solid var(--color-gold);
        }

        .card h2 {
            color: var(--color-dark);
            margin-bottom: 20px;
            display: flex;
            align-items: center;
            gap: 10px;
            font-size: 1.5em;
        }

        .card h2 i {
            color: var(--color-gold);
            font-size: 1.3em;
        }

        @keyframes fadeIn {
            from { opacity: 0; }
            to { opacity: 1; }
        }

        @keyframes slideUp {
            from {
                opacity: 0;
                transform: translateY(20px);
            }
            to {
                opacity: 1;
                transform: translateY(0);
            }
        }

        .form-group {
            margin-bottom: 20px;
        }

        label {
            display: block;
            margin-bottom: 8px;
            font-weight: 600;
            color: var(--color-dark);
            display: flex;
            align-items: center;
            gap: 8px;
        }

        label i {
            color: var(--color-gold);
            font-size: 1.1em;
        }

        input[type="text"],
        select {
            width: 100%;
            padding: 12px;
            border: 2px solid var(--color-gold);
            border-radius: 5px;
            font-size: 1em;
            transition: border-color 0.3s;
            background: white;
            color: var(--color-dark);
        }

        input[type="text"]:focus,
        select:focus {
            outline: none;
            border-color: var(--color-brown);
            box-shadow: 0 0 8px rgba(167, 150, 117, 0.3);
        }

        .button-group {
            display: flex;
            gap: 10px;
        }

        button {
            flex: 1;
            padding: 12px;
            background: linear-gradient(135deg, var(--color-gold) 0%, var(--color-brown) 100%);
            color: white;
            border: none;
            border-radius: 5px;
            font-size: 1em;
            font-weight: 600;
            cursor: pointer;
            transition: transform 0.2s, box-shadow 0.2s;
            display: flex;
            align-items: center;
            justify-content: center;
            gap: 8px;
        }

        button i {
            font-size: 1.1em;
        }

        button:hover {
            transform: translateY(-2px);
            box-shadow: 0 5px 20px rgba(167, 150, 117, 0.4);
        }

        button:active {
            transform: translateY(0);
        }

        .clear-btn {
            background: linear-gradient(135deg, #E8B4A8 0%, #D4856D 100%);
        }

        .clear-btn:hover {
            box-shadow: 0 5px 20px rgba(232, 180, 168, 0.4);
        }

        .results {
            grid-column: 1 / -1;
        }

        .resultado-section {
            margin-bottom: 30px;
        }

        .resultado-section h2 {
            color: var(--color-dark);
            margin-bottom: 15px;
            border-bottom: 2px solid var(--color-gold);
            padding-bottom: 10px;
            display: flex;
            align-items: center;
            gap: 10px;
            font-size: 1.3em;
        }

        .resultado-section h2 i {
            color: var(--color-gold);
        }

        .mensaje {
            padding: 15px;
            border-radius: 5px;
            margin-bottom: 20px;
            font-weight: 600;
            animation: slideIn 0.3s ease-out;
            display: flex;
            align-items: center;
            gap: 10px;
        }

        .mensaje i {
            font-size: 1.5em;
        }

        .mensaje.exito {
            background: #D9EDDA;
            color: #155724;
            border-left: 4px solid #28a745;
        }

        .mensaje.exito i {
            color: #28a745;
        }

        .mensaje.error {
            background: #F8D7DA;
            color: #721c24;
            border-left: 4px solid #F5576C;
        }

        .mensaje.error i {
            color: #F5576C;
        }

        @keyframes slideIn {
            from {
                opacity: 0;
                transform: translateX(-20px);
            }
            to {
                opacity: 1;
                transform: translateX(0);
            }
        }

        .pasos-tabla {
            width: 100%;
            border-collapse: collapse;
            margin-bottom: 20px;
        }

        .pasos-tabla thead {
            background: linear-gradient(135deg, var(--color-gold) 0%, var(--color-brown) 100%);
            color: white;
        }

        .pasos-tabla th,
        .pasos-tabla td {
            padding: 12px;
            text-align: left;
            border-bottom: 1px solid #ddd;
        }

        .pasos-tabla th {
            font-weight: 600;
        }

        .pasos-tabla tbody tr:hover {
            background: rgba(167, 150, 117, 0.1);
        }

        .arbol-sintactico {
            background: #FAFAF8;
            padding: 15px;
            border-radius: 5px;
            font-family: 'Courier New', monospace;
            overflow-x: auto;
            border-left: 4px solid var(--color-gold);
            color: var(--color-dark);
        }

        .loading {
            display: none;
            text-align: center;
            padding: 40px 20px;
        }

        .spinner {
            border: 4px solid rgba(167, 150, 117, 0.2);
            border-top: 4px solid var(--color-gold);
            border-radius: 50%;
            width: 50px;
            height: 50px;
            animation: spin 1s linear infinite;
            margin: 0 auto 20px;
        }

        @keyframes spin {
            0% { transform: rotate(0deg); }
            100% { transform: rotate(360deg); }
        }

        .ejemplos {
            display: flex;
            flex-wrap: wrap;
            gap: 10px;
            margin-top: 20px;
        }

        .ejemplo-btn {
            padding: 8px 12px;
            background: white;
            border: 2px solid var(--color-gold);
            border-radius: 5px;
            cursor: pointer;
            transition: all 0.3s;
            font-size: 0.9em;
            color: var(--color-dark);
            font-weight: 600;
        }

        .ejemplo-btn:hover {
            background: var(--color-gold);
            color: white;
            border-color: var(--color-brown);
            transform: translateY(-2px);
        }

        .info-box {
            background: rgba(167, 150, 117, 0.1);
            padding: 15px;
            border-radius: 5px;
            border-left: 3px solid var(--color-gold);
        }

        .info-box p {
            margin-bottom: 10px;
            line-height: 1.8;
        }

        .info-box p:last-child {
            margin-bottom: 0;
        }

        .info-box strong {
            color: var(--color-dark);
        }

        @media (max-width: 768px) {
            .main-content {
                grid-template-columns: 1fr;
            }

            header h1 {
                font-size: 1.8em;
            }

            .button-group {
                flex-direction: column;
            }

            button {
                width: 100%;
            }
        }
    </style>
</head>
<body>
    <div class="container">
        <header>
            <h1>
                <i class="bi bi-diagram-3"></i>
                Analizador Sintáctico
            </h1>
            <p>Analiza expresiones matemáticas con diferentes métodos</p>
        </header>

        <div class="main-content">
            <div class="card">
                <h2><i class="bi bi-gear"></i> Configuración</h2>
                
                <div class="form-group">
                    <label for="expresion"><i class="bi bi-input-cursor-text"></i> Expresión:</label>
                    <input type="text" id="expresion" placeholder="Ejemplo: 3+5*2">
                    <div class="ejemplos">
                        <button class="ejemplo-btn" onclick="setEjemplo('3+5*2')">3+5*2</button>
                        <button class="ejemplo-btn" onclick="setEjemplo('3+5*(2+1)')">3+5*(2+1)</button>
                        <button class="ejemplo-btn" onclick="setEjemplo('(3+5)*2')"> (3+5)*2</button>
                        <button class="ejemplo-btn" onclick="setEjemplo('2*3+4')">2*3+4</button>
                    </div>
                </div>

                <div class="form-group">
                    <label for="metodo"><i class="bi bi-list-check"></i> Método:</label>
                    <select id="metodo">
                        <option value="recursivo">Predictivo Recursivo</option>
                        <option value="norecursivo">Predictivo No Recursivo</option>
                        <option value="lr">Ascendente LR</option>
                    </select>
                </div>

                <div class="button-group">
                    <button onclick="analizar()"><i class="bi bi-search"></i> Analizar</button>
                    <button class="clear-btn" onclick="limpiar()"><i class="bi bi-trash3"></i> Limpiar</button>
                </div>
            </div>

            <div class="card">
                <h2><i class="bi bi-info-circle"></i> Métodos</h2>
                <div class="info-box">
                    <p>
                        <strong><i class="bi bi-arrow-down-right"></i> Predictivo Recursivo:</strong><br>
                        Análisis descendente directo de la gramática.
                    </p>
                    <p>
                        <strong><i class="bi bi-arrow-down-right"></i> Predictivo No Recursivo:</strong><br>
                        Análisis descendente usando tabla predictiva.
                    </p>
                    <p>
                        <strong><i class="bi bi-arrow-up-left"></i> Ascendente LR:</strong><br>
                        Análisis ascendente construcción bottom-up.
                    </p>
                </div>
            </div>

            <div class="results card">
                <div class="loading">
                    <div class="spinner"></div>
                    <p><i class="bi bi-hourglass-split"></i> Analizando...</p>
                </div>
                
                <div id="resultado" style="display: none;"></div>
            </div>
        </div>
    </div>

    <script>
        function setEjemplo(expr) {
            document.getElementById('expresion').value = expr;
        }

        function limpiar() {
            document.getElementById('expresion').value = '';
            document.getElementById('resultado').innerHTML = '';
            document.getElementById('resultado').style.display = 'none';
        }

        async function analizar() {
            const expresion = document.getElementById('expresion').value.trim();
            const metodo = document.getElementById('metodo').value;

            if (!expresion) {
                alert('Por favor ingresa una expresión');
                return;
            }

            const loading = document.querySelector('.loading');
            const resultado = document.getElementById('resultado');

            loading.style.display = 'block';
            resultado.style.display = 'none';

            try {
                const response = await fetch('/api/analizar', {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json',
                    },
                    body: JSON.stringify({
                        expresion: expresion,
                        metodo: metodo
                    })
                });

                const data = await response.json();
                mostrarResultado(data);
            } catch (error) {
                resultado.innerHTML = '<div class="mensaje error"><i class="bi bi-exclamation-circle"></i> Error: ' + error.message + '</div>';
                resultado.style.display = 'block';
            } finally {
                loading.style.display = 'none';
            }
        }

        function mostrarResultado(data) {
            const resultado = document.getElementById('resultado');
            let html = '';

            html += '<div class="resultado-section">';
            if (data.exito) {
                html += '<div class="mensaje exito"><i class="bi bi-check-circle-fill"></i> ' + data.mensaje + '</div>';
            } else {
                html += '<div class="mensaje error"><i class="bi bi-x-circle-fill"></i> ' + data.mensaje + '</div>';
            }

            if (data.pasos && data.pasos.length > 0) {
                html += '<h2><i class="bi bi-list-ol"></i> Pasos del Análisis</h2>';
                html += '<table class="pasos-tabla">';
                html += '<thead><tr><th><i class="bi bi-hash"></i> Paso</th><th><i class="bi bi-file-text"></i> Descripción</th></tr></thead>';
                html += '<tbody>';
                
                data.pasos.forEach(paso => {
                    html += '<tr>';
                    html += '<td>' + paso.paso + '</td>';
                    html += '<td>' + paso.descripcion + '</td>';
                    html += '</tr>';
                });
                
                html += '</tbody></table>';
            }

            if (data.arbol) {
                html += '<h2><i class="bi bi-diagram-2"></i> Árbol Sintáctico</h2>';
                html += '<div class="arbol-sintactico">';
                html += '<pre>' + escapeHtml(data.arbol) + '</pre>';
                html += '</div>';
            }

            html += '</div>';
            resultado.innerHTML = html;
            resultado.style.display = 'block';
        }

        function escapeHtml(text) {
            const div = document.createElement('div');
            div.textContent = text;
            return div.innerHTML;
        }

        // Permitir analizar con Enter
        document.getElementById('expresion').addEventListener('keypress', function(event) {
            if (event.key === 'Enter') {
                analizar();
            }
        });
    </script>
</body>
</html>
`
