package httpadapter

import (
	"encoding/json"
	"net/http"
	"ticketera/internal/adapters/storage"
)

func TextoDeleteHandler(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Clave string `json:"clave"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "JSON inválido", http.StatusBadRequest)
		return
	}
	repo := storage.NewMongoRepository("mongodb://localhost:27017")
	err := repo.DeleteTexto(body.Clave)
	if err != nil {
		http.Error(w, "Error eliminando texto", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
