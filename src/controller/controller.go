package controller

import (
	"fmt"
	"smartgrowth-connectors/configapi/database"
	"smartgrowth-connectors/configapi/model"
)

type Controller struct {
	db database.Database
	User *model.User
}

func NewController(db database.Database, user *model.User) (*Controller, error) {
	return &Controller{db, user}, nil
}

func (ctr *Controller) AsUser(sub string) (*Controller, error) {

	// New Controller
	var newCtr *Controller 

	// Fetch user from database
	user, err := ctr.db.GetUserBySub(sub)
	if err != nil {
		return newCtr, fmt.Errorf("Error fetching user with sub %s from db: %v", sub, err)
	}

	newCtr, err = NewController(ctr.db, &user)
	if err != nil {
		return newCtr, fmt.Errorf("Error creating user controller: %v", err)
	}

	return newCtr, nil
}
