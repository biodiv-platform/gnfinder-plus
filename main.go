package main

import (
	"flag"
	"fmt"

	"github.com/gnames/gnfinder"
	"github.com/lu4p/cat"
)

func main() {
	filePath := flag.String("file", "_", "File Path")
	flag.Parse()
	txt, _ := cat.File(*filePath)
	gnf := gnfinder.NewGNfinder()
	output := gnf.FindNamesJSON([]byte(txt))
	fmt.Println(string(output))
}
