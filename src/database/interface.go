package database

import (
	"smartgrowth-connectors/configapi/model"
)

type Database interface {
	// Users
	InsertUser(model.User) (model.User, error)

	// Sources
	ListSources() ([]model.Source, error)
}
