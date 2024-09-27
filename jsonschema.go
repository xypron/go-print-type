package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/invopop/jsonschema"
	"go-print-type/image_definition"
)

// Define your Go type
type Person struct {
	Name    string `json:"name"`
	Age     int    `json:"age"`
	Email   string `json:"email"`
	Address string `json:"address"`
}

func main() {
	// Create a new reflector
	reflector := jsonschema.Reflector{}

	// Generate the JSON schema for the Person type
	schema := reflector.Reflect(imagedefinition.ImageDefinition{})

	// Marshal the schema to JSON
	schemaJSON, err := json.MarshalIndent(schema, "", "  ")
	if err != nil {
		fmt.Printf("Error marshaling schema to JSON: %v\n", err)
		return
	}

	// Write the JSON schema to a file
	file, err := os.Create("schema.json")
	if err != nil {
		fmt.Printf("Error creating file: %v\n", err)
		return
	}
	defer file.Close()

	_, err = file.Write(schemaJSON)
	if err != nil {
		fmt.Printf("Error writing to file: %v\n", err)
		return
	}

	fmt.Println("JSON schema written to schema.json")
}

