package controller

import (
	"fmt"
	"smartgrowth-connectors/configapi/database"
	"smartgrowth-connectors/configapi/model"
)

type Controller struct {
	database database.Database 
}

func NewController(db database.Database) (*Controller, error) {
	return &Controller{db}, nil
}

// Sources
func (c *Controller) ListSources() ([]model.Source, error){

	var result []model.Source

	result, err := c.database.ListSources()
	if err != nil {
		return result, fmt.Errorf("Failed to fetch from database: %v", err)
	}

	return result, nil
}
