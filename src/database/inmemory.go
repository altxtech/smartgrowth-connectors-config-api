package database

import (
	"errors"
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

func (db *inMemoryDB) ListSources() ([]model.Source, error) {
	
	var result []model.Source
	for _, value := range db.sources {
		result = append(result, value)
	}

	return result, nil
}
