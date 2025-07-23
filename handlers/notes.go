package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/capernix/gohttpx/models"
	"github.com/capernix/gohttpx/utils"
)

// CreateNote creates a new note
// @Summary Create a new note
// @Description Create a new note with a title and content
// @Tags notes
// @Accept  json
// @Produce  json
// @Param   note  body  models.Note  true  "Note to create"
// @Success 201 {object} models.Note
// @Failure 400 {object} map[string]string
// @Router /notes [post]
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

// GetNote returns a note by ID
// @Summary Get a note by ID
// @Description Get details of a note by its ID
// @Tags notes
// @Produce  json
// @Param   id   path   int  true  "Note ID"
// @Success 200 {object} models.Note
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /notes/{id} [get]
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

// DeleteNote deletes a note by ID
// @Summary Delete a note by ID
// @Description Delete a note by its ID
// @Tags notes
// @Param   id   path   int  true  "Note ID"
// @Success 204 {string} string "No Content"
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /notes/{id} [delete]
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

// ListNotes lists all notes
// @Summary List all notes
// @Description Get a list of all notes
// @Tags notes
// @Produce  json
// @Success 200 {array} models.Note
// @Router /notes [get]
func ListNotes(
	w http.ResponseWriter,
	r *http.Request,
) {
	notes := models.ListNotes()

	utils.WriteJSON(w, http.StatusOK, notes)
}
