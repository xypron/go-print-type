package main

import (
	"go/doc"
	"go/parser"
	"go/token"
	"log"
)

func main() {
	// Create a new file set
	fset := token.NewFileSet()

	// Parse the Go source file
	node, err := parser.ParseDir(fset, "./image_definition", nil, parser.ParseComments)
	if err != nil {
		log.Fatal(err)
	}

	// Iterate through the files in the package
	for _, pkg := range node {
		// Create a new package documentation
		pkgDoc := doc.New(pkg, "./image_definition", doc.AllDecls)

		// Print the package name and documentation
		log.Println("Package:", pkgDoc.Name)
		log.Println("Documentation:", pkgDoc.Doc)

		// Print documentation for each type
		for _, t := range pkgDoc.Types {
			log.Println("Type:", t.Name)
			log.Println("Documentation:", t.Doc)
		}
	}
}
