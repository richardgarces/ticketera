package httpadapter

import (
	"net/http"

	"github.com/gorilla/mux"
)

type Handler struct {
	Router *mux.Router
}

func NewHandler() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/ticketweb", TicketWebHandler).Methods("GET")
	r.HandleFunc("/tickets", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "web/tickets.html")
	})
	r.HandleFunc("/api/tickets", TicketListHandler).Methods("GET")
	r.HandleFunc("/", HomeHandler).Methods("GET")
	r.HandleFunc("/api/ticket", TicketHandler).Methods("POST")
	r.HandleFunc("/admin", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "web/admin.html")
	})
	r.HandleFunc("/api/texto", TextoHandler).Methods("POST")
	r.HandleFunc("/api/textos", TextosListHandler).Methods("GET")
	r.HandleFunc("/api/texto/delete", TextoDeleteHandler).Methods("POST")
	r.HandleFunc("/api/logo", LogoHandler).Methods("POST")
	// Servir archivos estáticos (logo, imágenes, etc.)
	r.PathPrefix("/web/").Handler(http.StripPrefix("/web/", http.FileServer(http.Dir("web/"))))
	return r
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "web/index.html")
}
