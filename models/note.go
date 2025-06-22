package models

import (
	"errors"
	"sync"
)

type Note struct {
	ID      int    `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
}

var (
	notes      = make(map[int]Note)
	notesMutex sync.RWMutex
	nextID     = 1
)

func CreateNote(title, content string) (Note, error) {
	if title == "" || content == "" {
		return Note{}, errors.New("title and content are required")
	}

	notesMutex.Lock()
	defer notesMutex.Unlock()

	note := Note{
		ID:      nextID,
		Title:   title,
		Content: content,
	}

	notes[nextID] = note
	nextID++

	return note, nil
}

func GetNote(id int) (Note, bool) {
	notesMutex.Lock()
	defer notesMutex.Unlock()

	note, exists := notes[id]
	return note, exists
}

func DeleteNote(id int) bool {
	notesMutex.Lock()
	defer notesMutex.Unlock()

	if _, exists := notes[id]; !exists {
		return false
	}

	delete(notes, id)
	return true
}

func ListNotes() []Note {
	notesMutex.RLock()
	defer notesMutex.RUnlock()

	noteList := make([]Note, 0, len(notes))
	for _, note := range notes {
		noteList = append(noteList, note)
	}
	return noteList
}
