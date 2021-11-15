package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"mime"
	"net/http"
	"net/url"
	"os"
	"path/filepath"

	"code.sajari.com/docconv"
	"github.com/gnames/gnfinder"
	"github.com/gnames/gnfinder/config"
	"github.com/gnames/gnfinder/ent/nlp"
	"github.com/gnames/gnfinder/io/dict"
	"github.com/gnames/gnfmt"
	"github.com/gofiber/fiber/v2"
	gonanoid "github.com/matoous/go-nanoid/v2"
)

func getFilePath(response *http.Response) string {
	mimeType := response.Header.Get("Content-Type")
	extensions, err := mime.ExtensionsByType(mimeType)
	filePrefix, _ := gonanoid.New()

	if err != nil || len(extensions) == 0 {
		return filePrefix
	}

	return filepath.Join(os.TempDir(), filePrefix+extensions[0])
}

func downloadFile(URL string) (string, error) {
	//Get the response bytes from the url
	response, err := http.Get(URL)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	tmpFile := getFilePath(response)

	if response.StatusCode != 200 {
		return "", errors.New("received non 200 response code")
	}

	//Create a empty file
	file, err := os.Create(tmpFile)
	if err != nil {
		return "", err
	}
	defer file.Close()

	//Write the bytes to the fiel
	_, err = io.Copy(file, response.Body)
	if err != nil {
		return "", err
	}

	return tmpFile, nil
}

// find document and extract text from it
func parseDocument(filePath string) string {

	txt, err := docconv.ConvertPath(filePath)
	if err != nil {
		log.Fatal(err)
	}

	return parseText(txt.Body)
}

// parse names from text through gnfinder
func parseText(txt string) string {

	cfg := config.New()
	gnf := gnfinder.New(cfg, dict.LoadDictionary(), nlp.BayesWeights())
	output := gnf.Find("", txt)

	return output.Format(gnfmt.PrettyJSON)
}

// HTTP Server
func server(serverPort string) {
	app := fiber.New()

	app.Get("/parse", func(c *fiber.Ctx) error {
		fullFilePath := c.Query("file")
		queryText := c.Query("text", "_")

		if queryText != "_" {
			return c.Type("json").SendString(parseText(queryText))
		}

		_, err := url.ParseRequestURI(fullFilePath)
		if err != nil {
			return c.SendStatus(fiber.StatusBadRequest)
		}

		localPath, _ := downloadFile(fullFilePath)
		return c.Type("json").SendString(parseDocument(localPath))
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

	_, err := url.ParseRequestURI(*filePath)
	if err != nil {
		// local file
		fmt.Println(parseDocument(*filePath))
	} else {
		// remote server
		localPath, _ := downloadFile(*filePath)
		fmt.Println(parseDocument(localPath))
	}
}
