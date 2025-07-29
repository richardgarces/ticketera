package httpadapter

import (
	"encoding/json"
	"net/http"
	"ticketera/internal/adapters/storage"
)

func TextosListHandler(w http.ResponseWriter, r *http.Request) {
	repo := storage.NewMongoRepository("mongodb://localhost:27017")
	textos, err := repo.GetTextos()
	if err != nil {
		http.Error(w, "Error obteniendo textos", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(textos)
}
