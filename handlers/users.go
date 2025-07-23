package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/capernix/gohttpx/models"
	"github.com/capernix/gohttpx/utils"
)

// CreateUser creates a new user
// @Summary Create a new user
// @Description Register a new user with a name
// @Tags users
// @Accept  json
// @Produce  json
// @Param   user  body  models.User  true  "User to create"
// @Success 201 {object} models.User
// @Failure 400 {object} map[string]string
// @Router /users [post]
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

// GetUser returns a user by ID
// @Summary Get a user by ID
// @Description Get details of a user by their ID
// @Tags users
// @Produce  json
// @Param   id   path   int  true  "User ID"
// @Success 200 {object} models.User
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /users/{id} [get]
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

// DeleteUser deletes a user by ID
// @Summary Delete a user by ID
// @Description Delete a user by their ID
// @Tags users
// @Param   id   path   int  true  "User ID"
// @Success 204 {string} string "No Content"
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /users/{id} [delete]
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

// ListUsers lists all users
// @Summary List all users
// @Description Get a list of all users
// @Tags users
// @Produce  json
// @Success 200 {array} models.User
// @Router /users [get]
func ListUsers(
	w http.ResponseWriter,
	r *http.Request,
) {
	users := models.ListUsers()

	utils.WriteJSON(w, http.StatusOK, users)
}
