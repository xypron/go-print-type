package main

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"

	"github.com/invopop/jsonschema"

	"go-print-type/image_definition"
)

type CustomReflector struct {
	reflector jsonschema.Reflector
}

// CollectTypes collects types from the given reflect.Type and stores them in the provided map.
func (r *CustomReflector) CollectTypes(m *map[string]reflect.Type, typ reflect.Type) {
	(*m)[typ.Name()] = typ
	switch typ.Kind() {
	case reflect.Ptr, reflect.Slice, reflect.Map:
		r.CollectTypes(m, typ.Elem())
	case reflect.Struct:
		for j := 0; j < typ.NumField(); j++ {
			field := typ.Field(j)
			r.CollectTypes(m, field.Type)
		}
	}
}

// Reflect generates a JSON schema for the given value.
func (r *CustomReflector) Reflect(value interface{}) *jsonschema.Schema {
	schema := r.reflector.Reflect(value)
	if schema == nil {
		fmt.Println("No schema generated")
		return nil
	}

	typeMap := make(map[string]reflect.Type)
	r.CollectTypes(&typeMap, reflect.TypeOf(value))

	for key, def := range schema.Definitions {
		if typ, exists := typeMap[key]; exists {
			r.processDefinition(def, typ)
		}
	}

	return schema
}

// processDefinition processes the schema definition and adds yaml tag as Extra
func (r *CustomReflector) processDefinition(schema *jsonschema.Schema, typ reflect.Type) {
	for element := schema.Properties.Newest(); element != nil; element = element.Prev() {
		for j := 0; j < typ.NumField(); j++ {
			field := typ.Field(j)
			if prop, exists := schema.Properties.Get(field.Name); exists {
				if prop.Extras == nil {
					prop.Extras = make(map[string]any)
				}
				if yaml := field.Tag.Get("yaml"); yaml != "" {
					prop.Extras["yaml"] = yaml
				}
			}
		}
	}
}

func main() {
	reflector := &CustomReflector{}
	schema := reflector.Reflect(imagedefinition.ImageDefinition{})

	schemaJSON, err := json.MarshalIndent(schema, "", "  ")
	if err != nil {
		fmt.Printf("Error marshaling schema to JSON: %v\n", err)
		return
	}

	file, err := os.Create("schema.json")
	if err != nil {
		fmt.Printf("Error creating file: %v\n", err)
		return
	}
	defer file.Close()

	if _, err = file.Write(schemaJSON); err != nil {
		fmt.Printf("Error writing to file: %v\n", err)
		return
	}

	fmt.Println("JSON schema written to schema.json")
}
