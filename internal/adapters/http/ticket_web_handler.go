package httpadapter

import (
	"html/template"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"ticketera/internal/adapters/storage"
	"time"
)

func TicketWebHandler(w http.ResponseWriter, r *http.Request) {
	// Consultar textos por defecto si algún campo viene vacío
	repo := storage.NewMongoRepository("mongodb://localhost:27017")
	textos, _ := repo.GetTextos()

	inicialStr := r.URL.Query().Get("inicial")
	inicial, _ := strconv.Atoi(inicialStr)
	if inicial < 1 {
		inicial = 1
	}
	header := r.URL.Query().Get("header")
	if header == "" {
		header = textos["header"]
	}
	title := r.URL.Query().Get("title")
	if title == "" {
		title = textos["title"]
	}
	content := r.URL.Query().Get("content")
	if content == "" {
		content = textos["content"]
	}
	pie := r.URL.Query().Get("pie")
	if pie == "" {
		pie = textos["pie"]
	}
	filtro := r.URL.Query().Get("filtro")
	if filtro == "" {
		if val, ok := textos["filtro"]; ok && val != "" {
			filtro = val
		} else {
			filtro = "Vale por:"
		}
	}
	cantidadStr := r.URL.Query().Get("cantidad")
	cantidad, _ := strconv.Atoi(cantidadStr)
	if cantidad < 1 {
		cantidad = 1
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	logo, _ := repo.GetLogoPath()
	// Si no hay logo, dejar vacío, pero si hay, usar el base64 de la base de datos
	color := r.URL.Query().Get("color")
	if color == "" {
		color = "#fff"
	}
	type TicketData struct {
		ENCABEZADO      string
		TITULO          string
		CONTENT         string // Para el lado cliente (sin filtrar)
		CONTROL_CONTENT string // Para el lado control (filtrado)
		PIE             string
		CORRELATIVO     int
		FECHA_ACTUAL    string
		URI_LOGO        template.URL
	}
	// Usar la fecha del formulario si viene, si no buscar en MongoDB, si no la actual
	fechaActual := r.URL.Query().Get("fecha")
	if fechaActual == "" {
		if val, ok := textos["fecha"]; ok && val != "" {
			fechaActual = val
		} else {
			fechaActual = time.Now().Format("2006-01-02")
		}
	}
	// Convertir a formato dd/mm/yyyy si viene en yyyy-mm-dd
	if len(fechaActual) == 10 && fechaActual[4] == '-' && fechaActual[7] == '-' {
		fechaActual = fechaActual[8:10] + "/" + fechaActual[5:7] + "/" + fechaActual[0:4]
	}
	tickets := make([]TicketData, cantidad)
	for i := 0; i < cantidad; i++ {
		tickets[i] = TicketData{
			ENCABEZADO:      header,
			TITULO:          title,
			CONTENT:         content,                      // sin filtrar
			CONTROL_CONTENT: filterFrase(content, filtro), // solo control, usando filtro dinámico
			PIE:             pie,
			CORRELATIVO:     inicial + i,
			FECHA_ACTUAL:    fechaActual,
			URI_LOGO:        template.URL(logo),
		}
	}
	// Estructura para pasar tickets y color al template
	// Leer ticketsPorFila del query string
	ticketsPorFilaStr := r.URL.Query().Get("ticketsPorFila")
	ticketsPorFila, _ := strconv.Atoi(ticketsPorFilaStr)
	if ticketsPorFila < 1 {
		ticketsPorFila = 4 // valor por defecto
	}
	tipoPagina := r.URL.Query().Get("tipoPagina")
	if tipoPagina == "" {
		tipoPagina = "carta"
	}
	data := struct {
		Tickets        []TicketData
		Color          string
		TicketsPorFila int
		TipoPagina     string
	}{
		Tickets:        tickets,
		Color:          color,
		TicketsPorFila: ticketsPorFila,
		TipoPagina:     tipoPagina,
	}

	wd, err := os.Getwd()
	if err != nil {
		http.Error(w, "Error obteniendo ruta base", 500)
		return
	}
	templateName := r.URL.Query().Get("template")
	if templateName == "" {
		templateName = "template_vale_por.html"
	}
	templatePath := wd + "/web/template/" + templateName
	funcMap := template.FuncMap{
		"mod": func(i, j int) int { return i % j },
		"add": func(i, j int) int { return i + j },
		"len": func(s interface{}) int {
			switch v := s.(type) {
			case []interface{}:
				return len(v)
			case []TicketData:
				return len(v)
			default:
				return 0
			}
		},
		// divf: división flotante para calcular el ancho de los tickets por fila
		"divf": func(a, b int) float64 {
			if b == 0 {
				return 0
			}
			return float64(a) / float64(b)
		},
	}
	baseName := templateName
	tmpl, err := template.New(baseName).Funcs(funcMap).ParseFiles(templatePath)
	if err != nil {
		http.Error(w, "Error cargando template", 500)
		return
	}
	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, "Error generando tickets", 500)
	}
}

// filterFrase elimina la frase indicada (insensible a mayúsculas) del texto
func filterFrase(s, frase string) string {
	if frase == "" {
		return s
	}
	re := regexp.MustCompile(`(?i)` + regexp.QuoteMeta(frase) + `[ ]*`)
	return re.ReplaceAllString(s, "")
}
