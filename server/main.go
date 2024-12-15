package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	db "rp/server/db"
)

type Server struct {
	port string
	sql  *db.Events
}

func createServer() *Server {
	events, err := db.CreateSchema()
	if err != nil {
		log.Fatalf("Error creating schema: %v", err)
	}

	port := os.Getenv("PORT") // Allow port to be configurable
	if port == "" {
		port = "8081"
	}

	server := &Server{
		port: port,
		sql:  events,
	}
	server.routes() // Set up routes during server creation
	return server
}

func main() {
	server := createServer()
	defer server.sql.Close()
	server.listen()
}

func returnJson(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

func (s *Server) routes() {
	http.HandleFunc("GET /events", s.corsMiddleware(s.handleGetEvents))
	http.HandleFunc("POST /events", s.corsMiddleware(s.handleCreateEvent))
	http.HandleFunc("PATCH /events/", s.corsMiddleware(s.handleUpdateEvent))
	http.HandleFunc("/", s.corsMiddleware(s.unhandled))

}

func (s *Server) listen() {
	port := s.port
	log.Printf("Server listening on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func (s *Server) corsMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, DELETE, PUT, PATCH")
		w.Header().Set("Access-Control-Allow-Headers", "content-type")
		w.Header().Set("Access-Control-Max-Age", "86400")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next(w, r)
	}
}

func (s *Server) handleGetEvents(w http.ResponseWriter, r *http.Request) {
	events, err := s.sql.GetEvents()
	if err != nil {
		http.Error(w, "Error getting events", http.StatusInternalServerError)
		return
	}
	returnJson(w, events)
}

func (s *Server) unhandled(w http.ResponseWriter, r *http.Request) { 
	println("Unhandled endpoint")
}
func (s *Server) handleCreateEvent(w http.ResponseWriter, r *http.Request) {
	var event db.Event
	err := json.NewDecoder(r.Body).Decode(&event)
	if err != nil {
		http.Error(w, "Error decoding event", http.StatusBadRequest)
		return
	}
	id, err := s.sql.AddEvent(event)
	if err != nil {
		http.Error(w, "Error adding event", http.StatusInternalServerError)
		return
	}

	event, err = s.sql.GetEventById(id)
	returnJson(w, event)
}

func (s *Server) handleUpdateEvent(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	id := path[len("/events/"):]
	if id == "" {
		http.Error(w, "ID not specified", http.StatusBadRequest)
		return
	}

	var event db.UpdateEvent
	err := json.NewDecoder(r.Body).Decode(&event)
	if err != nil {
		http.Error(w, "Error decoding event", http.StatusBadRequest)
		return
	}

	println("getting event with id: ", id)
	_, err = s.sql.GetEventById(id)
	if err != nil {
		http.Error(w, "Error while fetching existing event", http.StatusInternalServerError)
		println(err)
		return
	}

	err = s.sql.UpdateEvent(event, id)
	if err != nil {
		http.Error(w, "Error updating event", http.StatusInternalServerError)
		return
	}

	newEvent, err := s.sql.GetEventById(id)
	if err != nil {
		http.Error(w, "Error while fetching updated event", http.StatusInternalServerError)
		return
	}

	returnJson(w, newEvent)
}
