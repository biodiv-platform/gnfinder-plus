package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"code.sajari.com/docconv"
	"github.com/gnames/gnfinder"
	"github.com/gnames/gnfinder/config"
	"github.com/gnames/gnfinder/ent/nlp"
	"github.com/gnames/gnfinder/io/dict"
	"github.com/gnames/gnfmt"
	"github.com/gofiber/fiber/v2"
)

func parse(filePath string) string {
	// Attempt to read file
	txt, err := docconv.ConvertPath(filePath)
	if err != nil {
		log.Fatal(err)
	}

	// Run document contents through gnfinder
	cfg := config.New()
	gnf := gnfinder.New(cfg, dict.LoadDictionary(), nlp.BayesWeights())
	output := gnf.Find("", txt.Body)
	return output.Format(gnfmt.PrettyJSON)
}

func server(serverPort string) {
	app := fiber.New()

	app.Post("/parse", func(c *fiber.Ctx) error {

		file, err := c.FormFile("file")
		fullFilePath := os.TempDir() + string(os.PathSeparator) + file.Filename

		if err == nil {
			c.SaveFile(file, fullFilePath)
			parsedResponse := parse(fullFilePath)
			os.Remove(fullFilePath) // housekeeping

			return c.Type("json").SendString(parsedResponse)
		}

		return c.SendStatus(fiber.StatusBadRequest)
	})

	app.Listen(":" + serverPort)
}

func main() {
	// Parse CLI Arguements
	filePath := flag.String("file", "_", "File Path")
	serverPort := flag.String("port", "3006", "Server Port")
	flag.Parse()

	if *filePath == "_" {
		server(*serverPort)
	}

	fmt.Println(parse(*filePath))
}
