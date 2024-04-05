package database

import (
	"errors"
	"fmt"
	"smartgrowth-connectors/configapi/model"
	"time"

	"github.com/google/uuid"
)

type inMemoryDB struct {
	users map[string]model.User
	sources map[string]model.Source
}

func NewInMemoryDB() (Database, error) {
	return &inMemoryDB {
		users: map[string]model.User{},
		sources: map[string]model.Source{},
	}, nil
}

// User
func (db *inMemoryDB) GetUserBySub(sub string) (model.User, error) {
	for _, val := range db.users {
		if val.Sub == sub {
			return val, nil
		}
	}
	var result model.User
	return result, fmt.Errorf("User with sub %s not found", sub)
}

func (db *inMemoryDB) GetUserById(id string) (model.User, error) {
	if val, ok := db.users[id]; ok {
		return val, nil
	}
	var result model.User
	return result, fmt.Errorf("User with id %s not found", id)
}

func (db *inMemoryDB) InsertUser(u model.User) (model.User, error) {

	var result model.User

	// User should not be identified
	if u.ID != "" {
		return result, errors.New("User should not be identified")
	}

	id := uuid.NewString()
	u.ID = id

	u.CreatedAt = time.Now()
	u.UpdatedAt = time.Now()

	db.users[id] = u
	return u, nil
}

func (db *inMemoryDB) ListUsers(offset int, limit int) ([]model.User, error) {
	
	var result []model.User
	for _, value := range db.users {
		result = append(result, value)
	}

	return result, nil
}

func (db *inMemoryDB) UpdateUser(id string, u model.User) (model.User, error) {
	
	var result model.User

	// User should exist
	if _, ok := db.users[id]; !ok {
		return result, errors.New("User not found")
	}

	u.ID = id
	u.UpdatedAt = time.Now()

	db.users[id] = u
	return u, nil
}

func (db *inMemoryDB) DeleteUserById(id string) (model.User, error) {
	
	var result model.User

	// User should exist
	if _, ok := db.users[id]; !ok {
		return result, errors.New("User not found")
	}

	result = db.users[id]
	delete(db.users, id)
	return result, nil
}

func (db *inMemoryDB) ListSources() ([]model.Source, error) {
	
	var result []model.Source
	for _, value := range db.sources {
		result = append(result, value)
	}

	return result, nil
}
