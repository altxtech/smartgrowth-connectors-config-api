package model

import (
	"testing"
)

// The only thing we really need to test is the IntegrationConfig.Validate method

func TestValidEmptyConfig(t *testing.T){
	// Test an empty configuration
	config := IntegrationConfig{}
	schema := ConfigurationSchema{}
	err := config.Validate(schema)
	if err != nil {
		t.Errorf("Expected empty configuration to be valid, got %v", err)
	}
}

func TestValidFlatConfig(t *testing.T){
	config := IntegrationConfig{
		"key1": "value1",
		"key2": 1, 
		"key3": 3.14,
		"key4": true,
	}
	schema := ConfigurationSchema{
		SchemaField{"key1", "string", false, false, nil},
		SchemaField{"key2", "int", false, false, nil},
		SchemaField{"key3", "float", false, false, nil},
		SchemaField{"key4", "boolean", false, false, nil},
	}
	err := config.Validate(schema)
	if err != nil {
		t.Errorf("Expected flat configuration to be valid, got %v", err)
	}
}

func TestInvalidString(t *testing.T){
	config := IntegrationConfig{
		"key1": 1,
	}
	schema := ConfigurationSchema{
		SchemaField{"key1", "string", false, false, nil},
	}
	err := config.Validate(schema)
	if err == nil {
		t.Errorf("Expected error, got nil")
	}
}

func TestInvalidInt(t *testing.T){
	config := IntegrationConfig{
		"key1": "not an int",
	}
	schema := ConfigurationSchema{
		SchemaField{"key1", "int", false, false, nil},
	}
	err := config.Validate(schema)
	if err == nil {
		t.Errorf("Expected error, got nil")
	}
}

func TestInvalidFloat(t *testing.T){
	config := IntegrationConfig{
		"key1": "not a float",
	}
	schema := ConfigurationSchema{
		SchemaField{"key1", "float", false, false, nil},
	}
	err := config.Validate(schema)
	if err == nil {
		t.Errorf("Expected error, got nil")
	}
}

func TestInvalidBoolean(t *testing.T){
	config := IntegrationConfig{
		"key1": "not a boolean",
	}
	schema := ConfigurationSchema{
		SchemaField{"key1", "boolean", false, false, nil},
	}
	err := config.Validate(schema)
	if err == nil {
		t.Errorf("Expected error, got nil")
	}
}

func TestRequiredField(t *testing.T){
	schema := ConfigurationSchema{
		SchemaField{"key1", "string", true, false, nil},
		SchemaField{"key2", "int", false, false, nil},
	}
	config := IntegrationConfig{
		"key1": "value1",
	}
	err := config.Validate(schema)
	if err != nil {
		t.Errorf("Expected valid configuration, got %v", err)
	}

	config = IntegrationConfig{
		"key2": 1,
	}
	err = config.Validate(schema)
	if err == nil {
		t.Errorf("Expected error, got nil")
	}
}

func TestArrayField(t *testing.T){
	schema := ConfigurationSchema{
		SchemaField{"key1", "string", false, true, nil},
	}

	// Missing field pases (field not required)
	config := IntegrationConfig{}
	err := config.Validate(schema)
	if err != nil {
		t.Errorf("Expected valid configuration, got %v", err)
	}

	// Empty array passes (not required)
	config = IntegrationConfig{
		"key1": []interface{}{},
	}
	err = config.Validate(schema)
	if err != nil {
		t.Errorf("Expected valid configuration, got %v", err)
	}

	// Array with valid values passes
	config = IntegrationConfig{
		"key1": []interface{}{"value1", "value2"},
	}
	err = config.Validate(schema)
	if err != nil {
		t.Errorf("Expected valid configuration, got %v", err)
	}

	// Array of concrete type  (not interface)  fails
	config = IntegrationConfig{
		"key1": []string{"value1", "value2"},
	}
	err = config.Validate(schema)
	if err == nil {
		t.Errorf("Expected error, got nil")
	}

	// Array with invalid values fails
	config = IntegrationConfig{
		"key1": []interface{}{"value1", 1},
	}
	err = config.Validate(schema)
	if err == nil {
		t.Errorf("Expected error, got nil")
	}
}

func TestRequiredArrayField(t *testing.T){
	
	schema := ConfigurationSchema{
		SchemaField{"key1", "string", true, true, nil},
	}

	// Array of valid values passes
	config := IntegrationConfig{
		"key1": []interface{}{"value1", "value2"},
	}
	err := config.Validate(schema)
	if err != nil {
		t.Errorf("Expected valid configuration, got %v", err)
	}

	// Missing field fails
	config = IntegrationConfig{}
	err = config.Validate(schema)
	if err == nil {
		t.Errorf("Expected error, got nil")
	}

	// Empty array fails
	config = IntegrationConfig{
		"key1": []interface{}{},
	}
	err = config.Validate(schema)
	if err == nil {
		t.Errorf("Expected error, got nil")
	}
}

func TestNestedConfig(t *testing.T){

	// A somewhat complex schema to test the nesting features
	// A 2D Matrix of objects
	schema := ConfigurationSchema{
		SchemaField{"x", "object", true, true, ConfigurationSchema{
			SchemaField{"y", "object", true, true, ConfigurationSchema{
				// Basic Optional values
				SchemaField{"key1", "string", false, false, nil},
				SchemaField{"key2", "int", false, false, nil},
				SchemaField{"key3", "float", false, false, nil},
				SchemaField{"key4", "boolean", false, false, nil},
				// Required value
				SchemaField{"key5", "int", true, false, nil},
				// Arrays
				SchemaField{"key6", "int", true, true, nil},
				SchemaField{"key7", "int", false, true, nil},
			}},
		}},
	}

	// Missing X
	config := IntegrationConfig{}
	err := config.Validate(schema)
	if err == nil {
		t.Errorf("Expected error, got nil")
	}

	// X is not []interface{}
	config = IntegrationConfig{
		"x": "not an object",
	}
	err = config.Validate(schema)
	if err == nil {
		t.Errorf("Expected error, got nil")
	}

	// X is empty
	config = IntegrationConfig{
		"x": []interface{}{},
	}

	// X has and invalid type
	config = IntegrationConfig{
		"x": []interface{}{1},
	}
	err = config.Validate(schema)
	if err == nil {
		t.Errorf("Expected error, got nil")
	}

	// Missing Y
	config = IntegrationConfig{
		"x": []interface{}{
			IntegrationConfig{},
		},
	}
	err = config.Validate(schema)
	if err == nil {
		t.Errorf("Expected error, got nil")
	}

	// Y is not []interface{}
	config = IntegrationConfig{
		"x": []interface{}{
			IntegrationConfig{
				"y": "not an object",
			},
		},
	}
	err = config.Validate(schema)
	if err == nil {
		t.Errorf("Expected error, got nil")
	}

	// Y is Empty
	config = IntegrationConfig{
		"x": []interface{}{
			IntegrationConfig{
				"y": []interface{}{},
			},
		},
	}

	// Y has an invalid type
	config = IntegrationConfig{
		"x": []interface{}{
			IntegrationConfig{
				"y": []interface{}{1},
			},
		},
	}
	err = config.Validate(schema)
	if err == nil {
		t.Errorf("Expected error, got nil")
	}

	// Missing required field in y
	config = IntegrationConfig{
		"x": []interface{}{
			IntegrationConfig{
				"y": []interface{}{
					IntegrationConfig{},
				},
			},
		},
	}
	err = config.Validate(schema)
	if err == nil {
		t.Errorf("Expected error, got nil")
	}

	// Empty required array in y
	config = IntegrationConfig{
		"x": []interface{}{
			IntegrationConfig{
				"y": []interface{}{
					IntegrationConfig{
						"key5": 1, 
						"key6": []interface{}{},
					},
				},
			},
		},
	}
	err = config.Validate(schema)
	if err == nil {
		t.Errorf("Expected error, got nil")
	}

	// Invalid type in Y
	config = IntegrationConfig{
		"x": []interface{}{
			IntegrationConfig{
				"y": []interface{}{
					IntegrationConfig{
						"key2": "not an int",
						"key5": 1,
						"key6": []interface{}{1, 2},
					},
				},
			},
		},
	}
	err = config.Validate(schema)
	if err == nil {
		t.Errorf("Expected error, got nil")
	}

	// Valid config
	config = IntegrationConfig{
		"x": []interface{}{
			IntegrationConfig{
				"y": []interface{}{
					IntegrationConfig{
						"key1": "value1",
						"key2": 1,
						"key3": 3.14,
						"key4": true,
						"key5": 1,
						"key6": []interface{}{1, 2},
						"key7": []interface{}{1, 2},
					},
				},
			},
		},
	}
	err = config.Validate(schema)
	if err != nil {
		t.Errorf("Expected valid configuration, got %v", err)
	}
}
