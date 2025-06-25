package models

import (
	"database/sql"

	"github.com/capernix/gohttpx/database"
)

type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func CreateUser(name string) (User, error) {
	var user User
	query := "INSERT INTO users (name) VALUES (?) RETURNING id, name"
	err := database.DB.QueryRow(query, name).Scan(&user.ID, &user.Name)
	if err != nil {
		return User{}, err
	}
	return user, nil
}

func GetUser(id int) (User, bool) {
	var user User
	query := "SELECT id, name FROM users WHERE id = ?"
	err := database.DB.QueryRow(query, id).Scan(&user.ID, &user.Name)
	if err != nil {
		if err == sql.ErrNoRows {
			return User{}, false
		}
		return User{}, false
	}
	return user, true
}

func DeleteUser(id int) bool {
	query := "DELETE FROM users WHERE id = ?"
	result, err := database.DB.Exec(query, id)
	if err != nil {
		return false
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return false
	}
	return rowsAffected > 0
}

func ListUsers() []User {
	query := "SELECT id, name FROM users"
	rows, err := database.DB.Query(query)
	if err != nil {
		return []User{}
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User
		if err := rows.Scan(&user.ID, &user.Name); err != nil {
			continue
		}
		users = append(users, user)
	}
	return users
}
