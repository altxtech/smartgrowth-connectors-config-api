package database

import (
	"smartgrowth-connectors/configapi/model"
)

type Database interface {
	// Users
	GetUserBySub(sub string) (model.User, error)
	InsertUser(model.User) (model.User, error)

}
