package model

import (
	"fmt"
)

type Integration struct {
	/*
		Represents a data integration. Can be either a data source or a data destination.
		Each data source belongs to a workspace, is tied to an Integration definition 
		and a configuration that should match the definition.

		For sake of simplicity, the permission  will only  exist for the workspace level.
		To view or edit an integration, the user must have the proper permissions for that workspace.
	*/
	ID string `json:"id" firestore:"id"`
	Name string `json:"name" firestore:"name"`
	WorkspaceID string `json:"workspace_id" firestore:"workspace_id"`
	DefinitionID string `json:"definition_id" firestore:"definition_id"`
	definition IntegrationDefinition // Definition denormalization
	Configuration IntegrationConfig `json:"configuration" firestore:"configuration"`
}

func NewIntegration(name string, workspaceID string, definition IntegrationDefinition, configuration IntegrationConfig) (Integration, error) {
	// Constructor be ignorant in respect to the state of the database
	integration := Integration{ "", name, workspaceID, definition.ID, definition, configuration }
	err := integration.Validate(integration.definition)
	if err != nil {
		return integration, fmt.Errorf("Invalid integration: %v", err)
	}

	return integration, nil
}

func  (i Integration) Validate(def IntegrationDefinition) error {
	// Checks if the configuration matches the definition
	err := i.Configuration.Validate(def.ConfigurationSchema)
	if err != nil {
		return fmt.Errorf("Invalid configuration: %v", err)
	}
	return nil
}

// Holds the concrete values for the configuration field
// Mapped as [field_name] => value
type IntegrationConfig map[string]interface{}

func NewIntegrationConfig(args map[string]interface{}, schema ConfigurationSchema) (IntegrationConfig, error) {
	config := IntegrationConfig(args)

	// Validate configuration
	err := config.Validate(schema)
	if err != nil {
		return config, fmt.Errorf("Invalid configuration: %v", err)
	}

	return config, nil
}

func (c IntegrationConfig) Validate(def ConfigurationSchema) error {
	// Checks if the configuration matches the schema

	for _, field := range def {
		value, ok := c[field.Name]
		if !ok {
			if field.Required {
				return fmt.Errorf("Field %s is required", field.Name)
			}
		} else {
			err := c.ValidateValue(field, value)
			if err != nil {
				return fmt.Errorf("Field %s is invalid: %v", field.Name, err)
			}
		}
	
	}

	return nil
}

func (c IntegrationConfig) ValidateValue(f ConfigurationSchemaField, value interface{}) error {
	
	if !f.Array {
		switch f.Type {
		case "string":
			_, ok := value.(string)
			if !ok {
				return fmt.Errorf("Expected string, got %T", value)
			}
		case "number":
			_, ok := value.(float64)
			if !ok {
				return fmt.Errorf("Expected number, got %T", value)
			}
		case "boolean":
			_, ok := value.(bool)
			if !ok {
				return fmt.Errorf("Expected boolean, got %T", value)
			}
		case "object":
			config, ok := value.(IntegrationConfig)
			if !ok {
				return fmt.Errorf("Expected object, got %T", value)
			}
			err := config.Validate(f.Fields)
			if err != nil {
				return fmt.Errorf("Invalid object: %v", err)
			}
		default:
			return fmt.Errorf("Invalid type %s", f.Type)
	}

	} else {
		// If array, value must be an Array
		items, ok := value.([]interface{})
		if !ok {
			return fmt.Errorf("Expected array, got %T", value)
		}

		if len(items) == 0 {
			if f.Required {
				return fmt.Errorf("Array is required")
			}
		} else {
			// Create a non-array copy of this field to evaluate each item
			nonArrayField := SchemaField{ f.Name, f.Type, f.Required, false, f.Fields }
			for _, item := range items {
				err := c.ValidateValue(nonArrayField, item)
				if err != nil {
					return fmt.Errorf("Invalid array item: %v", err)
				}
			}
		}
	}

	return nil
}
