package httpadapter

import (
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"ticketera/internal/adapters/storage"
)

func LogoHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}
	file, handler, err := r.FormFile("logo")
	if err != nil {
		http.Error(w, "Logo no recibido", http.StatusBadRequest)
		return
	}
	defer file.Close()
	logoPath := filepath.Join("web", "logo_"+handler.Filename)
	f, _ := os.Create(logoPath)
	defer f.Close()
	ioData, _ := ioutil.ReadAll(file)
	f.Write(ioData)
	storage.SaveLogoPath(logoPath)
	w.WriteHeader(http.StatusOK)
}
