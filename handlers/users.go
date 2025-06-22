package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/capernix/gohttpx/models"
	"github.com/capernix/gohttpx/utils"
)

func CreateUser(
	w http.ResponseWriter,
	r *http.Request,
) {
	var req struct {
		Name string `json:"name"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteError(
			w,
			"Invalid Request",
			http.StatusBadRequest,
		)
		return
	}

	user, err := models.CreateUser(req.Name)

	if err != nil {
		utils.WriteError(
			w,
			err.Error(),
			http.StatusBadRequest,
		)
		return
	}

	utils.WriteJSON(w, http.StatusCreated, user)
}

func GetUser(
	w http.ResponseWriter,
	r *http.Request,
) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		utils.WriteError(
			w,
			"invalid user ID",
			http.StatusBadRequest)
		return
	}

	user, exists := models.GetUser(id)
	if !exists {
		utils.WriteError(
			w,
			"user not found",
			http.StatusNotFound,
		)
		return
	}

	utils.WriteJSON(w, http.StatusOK, user)
}

func DeleteUser(
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

	if !models.DeleteUser(id) {
		utils.WriteError(
			w,
			"user not found",
			http.StatusNotFound,
		)
		return
	}

	w.WriteHeader(http.StatusNoContent)

}

func ListUsers(
	w http.ResponseWriter,
	r *http.Request,
) {
	users := models.ListUsers()

	utils.WriteJSON(w, http.StatusOK, users)
}
