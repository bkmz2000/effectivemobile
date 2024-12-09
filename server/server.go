package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"effectivemobile/db"
)

// ErrorResponse represents the structure of the error response
// @Description Error response model
// @Property message string true "Error message" example("An error occurred")
// @Property code int true "Error code" example(400)
type ErrorResponse struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}

type Server struct {
	db *db.SongDB
}

func NewServer() *Server {
	db := db.Init()
	return &Server{db}
}

// Filter handles GET requests to filter songs
// @Summary Filter songs
// @Description Filter songs by group, name, and date range
// @Tags songs
// @Accept json
// @Produce json
// @Param group query string false "Group"
// @Param name query string false "Name"
// @Param date_after query string false "Date After"
// @Param date_before query string false "Date Before"
// @Success 200 {array} db.Song
// @Failure 500 {object} ErrorResponse
// @Router /song [get]
func (s *Server) Filter(w http.ResponseWriter, r *http.Request) {
	log.Println("Filter method called")
	group := r.URL.Query().Get("group")
	name := r.URL.Query().Get("name")
	dateAfter := r.URL.Query().Get("date_after")
	dateBefore := r.URL.Query().Get("date_before")

	log.Printf("Query parameters - Group: %s, Name: %s, Date After: %s, Date Before: %s", group, name, dateAfter, dateBefore)

	songs, err := s.db.Filter(group, name, dateAfter, dateBefore)
	if err != nil {
		log.Printf("Error filtering songs: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{Message: "Internal Server Error", Code: http.StatusInternalServerError})
		return
	}

	log.Printf("Filtered songs: %v", songs)
	resp, err := json.Marshal(songs)
	if err != nil {
		log.Printf("Error marshaling response: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{Message: "Failed to marshal response", Code: http.StatusInternalServerError})
		return
	}

	w.Write(resp)
	log.Println("Filter response sent")
}

// Text handles GET requests to retrieve song text
// @Summary Get song text
// @Description Retrieve the text of a song by group and name
// @Tags songs
// @Accept json
// @Produce plain
// @Param group query string true "Group"
// @Param name query string true "Name"
// @Param verse query string false "Verse Index"
// @Success 200 {string} string
// @Failure 404 {object} ErrorResponse
// @Failure 400 {object} ErrorResponse
// @Router /text [get]
func (s *Server) Text(w http.ResponseWriter, r *http.Request) {
	log.Println("Text method called")
	group := r.URL.Query().Get("group")
	name := r.URL.Query().Get("name")
	verse := r.URL.Query().Get("verse")

	log.Printf("Query parameters - Group: %s, Name: %s, Verse: %s", group, name, verse)

	song := s.db.GetSong(group, name)
	if song == nil {
		log.Printf("Song not found for Group: %s, Name: %s", group, name)
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(ErrorResponse{Message: "Song not found", Code: http.StatusNotFound})
		return
	}

	text := strings.Split(song.Text, "\n\n")
	log.Printf("Song text split into %d verses", len(text))

	idx := 0
	if verse != "" {
		fmt.Sscanf(verse, "%d", &idx)
		log.Printf("Verse index parsed: %d", idx)
	}

	if idx >= len(text) {
		log.Printf("Invalid verse index: %d", idx)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Message: "Invalid verse index", Code: http.StatusBadRequest})
		return
	}

	w.Write([]byte(text[idx]))
	log.Println("Text response sent")
}

// Update handles PUT requests to update a song
// @Summary Update a song
// @Description Update an existing song
// @Tags songs
// @Accept json
// @Produce json
// @Param song body db.Song true "Song"
// @Success 200 {object} db.Song
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /song [put]
func (s *Server) Update(w http.ResponseWriter, r *http.Request) {
	log.Println("Update method called")
	var song db.Song

	if err := json.NewDecoder(r.Body).Decode(&song); err != nil {
		log.Printf("Error decoding song: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	log.Printf("Song to update: %+v", song)

	if err := s.db.Update(&song); err != nil {
		log.Printf("Error updating song: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(song)
	log.Println("Update response sent")
}

// Add handles POST requests to add a new song
// @Summary Add a new song
// @Description Add a new song to the database
// @Tags songs
// @Accept json
// @Produce json
// @Param song body db.Song true "Song"
// @Success 201 {object} db.Song
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /song [post]
func (s *Server) Add(w http.ResponseWriter, r *http.Request) {
	log.Println("Add method called")
	var song db.Song

	if err := json.NewDecoder(r.Body).Decode(&song); err != nil {
		log.Printf("Error decoding song: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	log.Printf("Song to add: %+v", song)

	if err := s.db.Add(&song); err != nil {
		log.Printf("Error adding song: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(song)
	log.Println("Add response sent")
}

// Delete handles DELETE requests to remove a song
// @Summary Delete a song
// @Description Remove a song from the database
// @Tags songs
// @Accept json
// @Produce json
// @Param id query string true "Song ID"
// @Success 204
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /song [delete]
func (s *Server) Delete(w http.ResponseWriter, r *http.Request) {
	log.Println("Delete method called")
	idStr := r.URL.Query().Get("id")
	log.Printf("Received ID for deletion: %s", idStr)

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		log.Printf("Invalid ID format: %s", idStr)
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	if err := s.db.Remove(uint(id)); err != nil {
		log.Printf("Error deleting song with ID %d: %v", id, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
	log.Println("Delete response sent")
}

func (s *Server) Song(w http.ResponseWriter, r *http.Request) {
	log.Printf("Received %s request for /song", r.Method)

	switch r.Method {
	case "GET":
		log.Println("Filtering songs")
		s.Filter(w, r)
	case "POST":
		log.Println("Adding a new song")
		s.Add(w, r)
	case "PUT":
		log.Println("Updating a song")
		s.Update(w, r)
	case "DELETE":
		log.Println("Deleting a song")
		s.Delete(w, r)
	default:
		log.Printf("Method not allowed: %s", r.Method)
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (s *Server) ServeHTTP() error {
	http.HandleFunc("/song", s.Song)
	http.HandleFunc("/text", s.Text)

	return http.ListenAndServe(":8080", nil)
}
