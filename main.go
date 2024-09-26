/*
Package go-print-type uses reflection to print complex types.
*/
package main

import (
	"fmt"
	"reflect"
	"go-print-type/image_definition"
)

// PrintTypeOf prints the definition of any complex type
//
// # Parameters
// - `typ`: the type to be printed
// - `indent`: a string printed in front of all output lines
func PrintTypeOf(typ reflect.Type, indent string) {
	if indent == "" {
		fmt.Printf("(%s)\n", typ)
	}
	switch typ.Kind() {
	case reflect.Ptr, reflect.Slice:
		PrintTypeOf(typ.Elem(), indent);
	case reflect.Struct:
		for j := 0; j < typ.NumField(); j++ {
			field := typ.Field(j)
			yaml := field.Tag.Get("yaml")
			if yaml == "" {
				yaml = "????"
			}
			schema := field.Tag.Get("jsonschema")
			fmt.Printf("%s|-%s (%s) %s\n", indent, yaml, field.Type, schema)
			PrintTypeOf(field.Type, indent + "| ")
		}
	}
}

func main() {
	PrintTypeOf(reflect.TypeOf(imagedefinition.ImageDefinition{}), "")
}
