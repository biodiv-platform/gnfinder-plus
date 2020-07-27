package main

import (
	"flag"
	"fmt"
	"log"

	"code.sajari.com/docconv"
	"github.com/gnames/gnfinder"
)

func main() {
	// Parse CLI Arguements
	filePath := flag.String("file", "_", "File Path")
	flag.Parse()

	// Attempt to read file
	txt, err := docconv.ConvertPath(*filePath)
	if err != nil {
		log.Fatal(err)
	}

	// Run document contents through gnfinder
	gnf := gnfinder.NewGNfinder()
	output := gnf.FindNamesJSON([]byte(txt.Body))
	fmt.Println(string(output))
}
