// SPDX-License-Identifier: MIT

package main

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"strings"

	"github.com/invopop/jsonschema"

	"github.com/xypron/go-print-type/image_definition"
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
		fmt.Println("### No schema generated")
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
	for j := 0; j < typ.NumField(); j++ {
		field := typ.Field(j)
		json := field.Tag.Get("json")
		if json == "" {
			fmt.Printf("No json tag for %s\n", json)
			json = field.Name
		} else {
			json = strings.Split(json, ",")[0]
		}
		match := false
		for pair := schema.Properties.Newest(); pair != nil; pair = pair.Prev() {

			if pair.Key != json {
				continue
			}
			match = true
			prop := pair.Value
			if prop.Extras == nil {
				prop.Extras = make(map[string]any)
			}
			if field.Name != json {
				fmt.Printf("Go %s does not match JSON %s\n",
					   field.Name, json);
				prop.Extras["go"] = field.Name
			}
			yaml := field.Tag.Get("yaml")
			if yaml == "" {
				fmt.Printf("No yaml tag for %s\n", json)
				yaml = strings.ToLower(field.Name)
			}
			prop.Extras["yaml"] = yaml
		}
		if !match {
			fmt.Printf("### Field not found %s\n", json)
		}
	}
}

func run(outputFile string) int {
	reflector := &CustomReflector{}
	schema := reflector.Reflect(imagedefinition.ImageDefinition{})

	schemaJSON, err := json.MarshalIndent(schema, "", "  ")
	if err != nil {
		fmt.Printf("Error marshaling schema to JSON: %v\n", err)
		return 1
	}

	file, err := os.Create(outputFile)
	if err != nil {
		fmt.Printf("Error creating file: %v\n", err)
		return 1
	}
	defer file.Close()

	if _, err = file.Write(schemaJSON); err != nil {
		fmt.Printf("Error writing to file: %v\n", err)
		return 1
	}

	fmt.Printf("JSON schema written to %s\n", outputFile)

	return 0
}

func usage() {
	fmt.Printf("Usage:\n" + os.Args[0] + " <filename>\n");
}

func main() {
	if len(os.Args) != 2 {
		usage()
		os.Exit(1)
	}

	ret := run(os.Args[1])

	os.Exit(ret)
}
