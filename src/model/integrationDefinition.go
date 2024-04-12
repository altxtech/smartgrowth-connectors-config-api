package model

import (
	"errors"
	"fmt"
)

type IntegrationDefinition struct {
	ID string `json:"id" firestore:"id"`
	Name string `json:"name" firestore:"name"`
	Type string `json:"type" firestore:"type"` // "source" or "destination"
	ConfigurationSchema ConfigurationSchema `json:"configuration_schema"  firestore:"configuration_schema"`
}

func NewIntegrationDefinition(name string, t string, schema ConfigurationSchema) (IntegrationDefinition, error) {
	def := IntegrationDefinition{ "", name, t, schema }
	err := def.Validate()
	if err != nil {
		return def, fmt.Errorf("Invalid definition: %v", err)
	}

	return def, nil
}

func (d IntegrationDefinition) Validate() error {

	// Type must be either "source" or "destination"
	if d.Type != "source" && d.Type != "destination" {
		return fmt.Errorf("Invalid type %s. Valid types are \"source\" and \"destination\"", d.Type)
	}

	err := d.ConfigurationSchema.Validate(3)
	if err != nil {
		return fmt.Errorf("Invalid configuration schema: %v", err)
	}

	return nil
}

type ConfigurationSchema []SchemaField

func (s ConfigurationSchema) Validate(remainingDepth  int) error {

	if remainingDepth < 0 {
		return errors.New("Max depth exceeded")
	}

	// There should be no name  collisions
	names := make(map[string]bool)

	for _, field := range s {

		if _, ok := names[field.Label]; ok {
			return fmt.Errorf("Field name %s is duplicated", field.Label)
		}
		names[field.Label] = true

		err := field.Validate(remainingDepth)
		if err != nil {
			return fmt.Errorf("Invalid field: %v", err)
		}
	}
	
	return nil
}

type SchemaField struct {
	Label string `json:"label" firestore:"label"`
	Type string  `json:"type" firestore:"type"` // "string", "int", "float",  "decimal", "boolean" or "object"
	Required bool  `json:"required" firestore:"required"`
	Array bool `json:"array" firesotre:"array"`
	Fields ConfigurationSchema `json:"fields" firestore:"fields"` // Only for "object" types
}

func (f SchemaField) Validate(remainingDepth int)  error {

	validTypes := map[string]bool{
		"string": true,
		"int": true,
		"float": true,
		"decimal": true,
		"boolean": true,
		"object": true,
	}
	_, ok := validTypes[f.Type]
	if !ok {
		return fmt.Errorf("Invalid type %s", f.Type)
	}

	if f.Type == "object" {
		
		err := f.Fields.Validate(remainingDepth - 1)
		if err != nil {
			return fmt.Errorf("Invalid object: %v", err)
		}
	}

	return nil
}
