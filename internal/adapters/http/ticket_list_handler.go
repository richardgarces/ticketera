package httpadapter

import (
	"encoding/json"
	"net/http"
	"ticketera/internal/adapters/storage"
)

func TicketListHandler(w http.ResponseWriter, r *http.Request) {
	repo := storage.NewMongoRepository("mongodb://localhost:27017")
	tickets, err := repo.GetTickets()
	if err != nil {
		http.Error(w, "Error obteniendo tickets", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(tickets)
}
