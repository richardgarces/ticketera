package main

import (
	"log"
	"net/http"

	httpadapter "ticketera/internal/adapters/http"
)

func main() {
	handler := httpadapter.NewHandler()
	log.Println("Servidor iniciado en :8080")
	http.ListenAndServe(":8080", handler)
}
