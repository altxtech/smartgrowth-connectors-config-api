package database

import (
	"smartgrowth-connectors/configapi/model"
)

type Database interface {
	// Users
	GetUserBySub(sub string) (model.User, error)
	GetUserById(id string) (model.User, error)
	InsertUser(model.User) (model.User, error)
	ListUsers(offset int, limit int) ([]model.User, error)
	UpdateUser(id string, user model.User) (model.User, error)
	DeleteUserById(id string) (model.User, error)

	// Workspace
	InsertWorkspace(model.Workspace) (model.Workspace, error)
	ListWorkspacesForPrincipal(string) ([]model.Workspace, error)
	GetWorkspaceByID(string) (model.Workspace, error)
	UpdateWorkspace(model.Workspace) (model.Workspace, error)
	DeleteWorkspaceByID(id string)  (model.Workspace, error)
}
