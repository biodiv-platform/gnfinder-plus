package main

import (
	"flag"
	"fmt"
	"log"

	"code.sajari.com/docconv"
	"github.com/gnames/gnfinder"
	"github.com/gnames/gnfinder/config"
	"github.com/gnames/gnfinder/ent/nlp"
	"github.com/gnames/gnfinder/io/dict"
	"github.com/gnames/gnfmt"
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
	cfg := config.New()
	gnf := gnfinder.New(cfg, dict.LoadDictionary(), nlp.BayesWeights())
	output := gnf.Find(txt.Body)
	fmt.Println(output.Format(gnfmt.PrettyJSON))
}
