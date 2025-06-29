package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/capernix/gohttpx/models"
	"github.com/capernix/gohttpx/utils"
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
		utils.WriteError(
			w,
			"Invalid Request",
			http.StatusBadRequest,
		)
		return
	}

	note, err := models.CreateNote(req.Title, req.Content)

	if err != nil {
		utils.WriteError(
			w,
			err.Error(),
			http.StatusBadRequest,
		)
		return
	}

	utils.WriteJSON(w, http.StatusCreated, note)
}

func GetNote(
	w http.ResponseWriter,
	r *http.Request,
) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		utils.WriteError(
			w,
			"invalid ID",
			http.StatusBadRequest)
		return
	}

	note, exists := models.GetNote(id)
	if !exists {
		utils.WriteError(
			w,
			"note not found",
			http.StatusNotFound,
		)
		return
	}

	utils.WriteJSON(w, http.StatusOK, note)
}

func DeleteNote(
	w http.ResponseWriter,
	r *http.Request,
) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		utils.WriteError(
			w,
			"invalid ID",
			http.StatusBadRequest,
		)
		return
	}

	if !models.DeleteNote(id) {
		utils.WriteError(
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

	utils.WriteJSON(w, http.StatusOK, notes)
}
