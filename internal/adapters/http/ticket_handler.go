package httpadapter

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	docxadapter "ticketera/internal/adapters/docx"
	"ticketera/internal/adapters/storage"
	"ticketera/internal/domain"
)

func TicketHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		http.Error(w, "Error al procesar el formulario", http.StatusBadRequest)
		return
	}
	title := r.FormValue("title")
	content := r.FormValue("content")
	pie := r.FormValue("pie")

	// Si se sube un logo, guárdalo temporalmente y usa ese, si no usa el de la base
	logoPath := ""
	file, handler, err := r.FormFile("logo")
	if err == nil {
		defer file.Close()
		logoPath = filepath.Join(os.TempDir(), handler.Filename)
		f, _ := os.Create(logoPath)
		defer f.Close()
		io.Copy(f, file)
	} else {
		repo := storage.NewMongoRepository("mongodb://localhost:27017")
		logoPath, _ = repo.GetLogoPath()
	}

	repo := storage.NewMongoRepository("mongodb://localhost:27017")
	textos, _ := repo.GetTextos()

	// Guardar ticket y obtener el número correlativo
	ticket := domain.Ticket{
		Title:   title,
		Content: content,
		Pie:     pie,
		LogoURL: logoPath,
	}
	_ = repo.SaveTicket(&ticket)

	data := map[string]string{
		"{{TITLE}}":   ticket.Title,
		"{{CONTENT}}": ticket.Content,
		"{{PIE}}":     ticket.Pie,
		"{{NUMERO}}":  fmt.Sprintf("%d", ticket.Numero),
		"{{LOGO}}":    ticket.LogoURL,
	}
	for k, v := range textos {
		data["{{"+k+"}}"] = v
	}

	templatePath := "web/ticket_template.docx"
	outputPath := filepath.Join(os.TempDir(), "ticket.docx")
	err = docxadapter.GenerateTicketDocx(templatePath, outputPath, data)
	if err != nil {
		http.Error(w, "Error generando el ticket", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Disposition", "attachment; filename=ticket.docx")
	w.Header().Set("Content-Type", "application/vnd.openxmlformats-officedocument.wordprocessingml.document")
	http.ServeFile(w, r, outputPath)
}
