package models

import (
	"database/sql"
	"errors"

	"github.com/capernix/gohttpx/database"
)

type Note struct {
	ID      int    `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
}

func CreateNote(title, content string) (Note, error) {
	if title == "" || content == "" {
		return Note{}, errors.New("title and content are required")
	}

	var note Note
	query := "INSERT INTO notes (title, content) VALUES (?, ?) RETURNING id, title, content"
	err := database.DB.QueryRow(query, title, content).Scan(&note.ID, &note.Title, &note.Content)
	if err != nil {
		return Note{}, err
	}
	return note, nil
}

func GetNote(id int) (Note, bool) {
	var note Note
	query := "SELECT id, title, content FROM notes WHERE id = ?"
	err := database.DB.QueryRow(query, id).Scan(&note.ID, &note.Title, &note.Content)
	if err != nil {
		if err == sql.ErrNoRows {
			return Note{}, false
		}
		return Note{}, false
	}
	return note, true
}

func DeleteNote(id int) bool {
	query := "DELETE FROM notes WHERE id = ?"
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

func ListNotes() []Note {
	query := "SELECT id, title, content FROM notes"
	rows, err := database.DB.Query(query)
	if err != nil {
		return []Note{}
	}
	defer rows.Close()

	var notes []Note
	for rows.Next() {
		var note Note
		if err := rows.Scan(&note.ID, &note.Title, &note.Content); err != nil {
			continue
		}
		notes = append(notes, note)
	}
	return notes
}
