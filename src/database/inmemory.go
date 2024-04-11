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
	workspaces map[string]model.Workspace
}

func NewInMemoryDB() (Database, error) {
	return &inMemoryDB {
		users: map[string]model.User{},
		workspaces: map[string]model.Workspace{},
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

// Workspaces
func (db *inMemoryDB) InsertWorkspace(w model.Workspace) (model.Workspace, error) {

	var idW model.Workspace

	// Workspace should not be identified
	if w.ID != "" {
		return idW, errors.New("Workspace should not be identified")
	}

	id := uuid.NewString()
	w.ID = id


	db.workspaces[id] = w
	return w, nil
} 

func (db *inMemoryDB) ListWorkspacesForPrincipal(principal string) ([]model.Workspace, error) {

	results := []model.Workspace{}

	for _, val := range db.workspaces {
		if val.ViewableBy(principal) {
			results = append(results, val)
		}
	}

	return results, nil
}

func (db  *inMemoryDB) GetWorkspaceByID(id string) (model.Workspace, error) {
	
	workspace := db.workspaces[id]
	return workspace, nil
}

func (db *inMemoryDB) UpdateWorkspace(w model.Workspace) (model.Workspace, error) {

	var upW model.Workspace

	// Workspace should be identified
	if w.ID == "" {
		return upW, errors.New("Workspace should be identified")
	}

	// Workpace should exist
	_, ok := db.workspaces[w.ID]
	if !ok {
		return upW, fmt.Errorf("Workspace with id %s does not exist", w.ID) 
	}

	db.workspaces[w.ID] = w

	return w, nil
} 

func (db *inMemoryDB) DeleteWorkspaceByID(id string) (model.Workspace, error) {
	//  Should exists
	deleteResult, ok := db.workspaces[id]
	if !ok {
		return deleteResult, errors.New("Workspace does not exists")
	}

	delete(db.workspaces, id)
	return deleteResult, nil
}
