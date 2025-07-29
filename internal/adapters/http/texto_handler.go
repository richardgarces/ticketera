package httpadapter

import (
	"encoding/json"
	"net/http"
	"ticketera/internal/adapters/storage"
)

var repo = storage.NewMongoRepository("mongodb://localhost:27017")

func TextoHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}
	var body struct {
		Clave string `json:"clave"`
		Valor string `json:"valor"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "JSON inválido", http.StatusBadRequest)
		return
	}
	err := repo.SaveTexto(body.Clave, body.Valor)
	if err != nil {
		http.Error(w, "Error guardando texto", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
