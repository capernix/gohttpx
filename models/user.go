package models

import (
	"errors"
	"sync"
)

type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

var (
	userCache  = make(map[int]User)
	userMutex  sync.RWMutex
	nextUserID = 1
)

func CreateUser(name string) (User, error) {
	if name == "" {
		return User{}, errors.New("name is required")
	}

	userMutex.Lock()
	defer userMutex.Unlock()

	user := User{
		ID:   nextUserID,
		Name: name,
	}

	userCache[nextUserID] = user
	nextUserID++

	return user, nil
}

func GetUser(id int) (User, bool) {
	userMutex.Lock()
	defer userMutex.Unlock()

	user, exists := userCache[id]
	return user, exists
}

func DeleteUser(id int) bool {
	userMutex.Lock()
	defer userMutex.Unlock()

	if _, exists := userCache[id]; !exists {
		return false
	}

	delete(userCache, id)
	return true
}
