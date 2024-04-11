package model

import (
	"testing"
)

func TestNewIntegrationDefinition1(t *testing.T) {
	
	// Tests a source with empty schema. Should work

	_, err := NewIntegrationDefinition("name", "source", ConfigurationSchema{})
	if err != nil {
		t.Errorf("Error creating integration definition: %v", err)
	}

	_, err = NewIntegrationDefinition("name", "destination", ConfigurationSchema{})
	if err != nil {
		t.Errorf("Error creating integration definition: %v", err)
	}
}

func TestInvalidType(t *testing.T) {
	
	_, err := NewIntegrationDefinition("name", "this type is invalid", ConfigurationSchema{})
	if err == nil {
		t.Errorf("Expected error, got nil")
	}
}

func TestValidFlatSchema(t *testing.T) {

	schema := ConfigurationSchema{
		SchemaField{"field1", "string", false, false, nil},
		SchemaField{"field2", "number", false, false, nil},
		SchemaField{"field3", "boolean", false, false, nil},
	}

	_, err := NewIntegrationDefinition("name", "source", schema)
	if err != nil {
		t.Errorf("Error creating integration definition: %v", err)
	}
}

func TestInvalidFieldType(t *testing.T) {
	schema := ConfigurationSchema{
		SchemaField{"field1", "this type is not valid", false, false, nil},
	}

	_, err := NewIntegrationDefinition("name", "source", schema)
	if err == nil {
		t.Errorf("Expected error, got nil")
	}
}

func TestValidNestedObject(t *testing.T) {
	schema := ConfigurationSchema{
		SchemaField{"field1", "object", false, false, ConfigurationSchema{
			SchemaField{"field1.1", "string", false, false, nil},
			SchemaField{"field1.2", "number", false, false, nil},
			SchemaField{"field1.3", "boolean", false, false, nil},
		}},
	}

	_, err := NewIntegrationDefinition("name", "source", schema)
	if err != nil {
		t.Errorf("Error creating integration definition: %v", err)
	}
}

func TestInvalidNestedObject(t *testing.T) {
	schema := ConfigurationSchema{
		SchemaField{"field1", "object", false, false, ConfigurationSchema{
			SchemaField{"field1.1", "string", false, false, nil},
			SchemaField{"field1.2", "number", false, false, nil},
			SchemaField{"field1.3", "this type is invalid", false, false, nil},
		}},
	}

	_, err := NewIntegrationDefinition("name", "source", schema)
	if err == nil {
		t.Errorf("Expected error, got nil")
	}
}

func TestValidMaxDepth(t *testing.T) {

	schema := ConfigurationSchema{
		SchemaField{"level1", "object", false, false, ConfigurationSchema{
			SchemaField{"level2", "object", false, false, ConfigurationSchema{
				SchemaField{"level3", "object", false, false, ConfigurationSchema{}},
			}},
		}},
	}

	_, err := NewIntegrationDefinition("name", "source", schema)
	if err != nil {
		t.Errorf("Error creating integration definition: %v", err)
	}
}

func TestInvalidMaxDepth(t *testing.T) {
	
	schema := ConfigurationSchema{
		SchemaField{"level1", "object", false, false, ConfigurationSchema{
			SchemaField{"level2", "object", false, false, ConfigurationSchema{
				SchemaField{"level3", "object", false, false, ConfigurationSchema{
					SchemaField{"level4", "object", false, false, ConfigurationSchema{}},
				}},
			}},
		}},
	}

	_, err := NewIntegrationDefinition("name", "source", schema)
	if err == nil {
		t.Errorf("Expected error, got nil")
	}
}
