package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/capernix/gohttpx/models"
)

func CreateNote(
	w http.ResponseWriter,
	r *http.Request,
) {
	var req struct {
		Title   string `json:"title"`
		Content string `json:"content"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(
			w,
			"Invalid Request",
			http.StatusBadRequest,
		)
		return
	}

	note, err := models.CreateNote(req.Title, req.Content)

	if err != nil {
		http.Error(
			w,
			err.Error(),
			http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(note)
}

func GetNote(
	w http.ResponseWriter,
	r *http.Request,
) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(
			w,
			"invalid ID",
			http.StatusBadRequest)
		return
	}

	note, exists := models.GetNote(id)
	if !exists {
		http.Error(
			w,
			"note not found",
			http.StatusNotFound,
		)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(note)
}

func DeleteNote(
	w http.ResponseWriter,
	r *http.Request,
) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(
			w,
			"invalid ID",
			http.StatusBadRequest,
		)
		return
	}

	if !models.DeleteNote(id) {
		http.Error(
			w,
			"note not found",
			http.StatusNotFound,
		)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func ListNotes(
	w http.ResponseWriter,
	r *http.Request,
) {
	notes := models.ListNotes()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(notes)
}
