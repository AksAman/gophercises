package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/AksAman/gophercises/linkparser"
	"github.com/AksAman/gophercises/utils"
)

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	// get current working directory
	currentDir, err := os.Getwd()
	check(err)

	htmlSrcDir := filepath.Join(currentDir, "html")

	htmlFiles, _ := os.ReadDir(htmlSrcDir)

	for _, htmlFileInfo := range htmlFiles {
		utils.PrintLine()
		htmlFileName := htmlFileInfo.Name()
		fmt.Printf("htmlFileName: %v\n", htmlFileName)
		links, err := parseLinks(filepath.Join(htmlSrcDir, htmlFileName))
		check(err)

		saveToJSON(links, htmlFileName)
		utils.PrintLine()
	}

}

func saveToJSON(links []linkparser.Link, srcHtmlFileName string) {
	data, err := json.MarshalIndent(links, "", "    ")
	check(err)

	// save json
	jsonFileStem := strings.Split(srcHtmlFileName, filepath.Ext(srcHtmlFileName))[0]
	jsonFileName := jsonFileStem + ".ignore.json"
	jsonFilePath := filepath.Join("data", jsonFileName)
	jsonFile, err := os.Create(jsonFilePath)
	check(err)

	_, err = jsonFile.Write(data)
	check(err)
}

func parseLinks(htmlFilePath string) ([]linkparser.Link, error) {
	htmlFile, err := os.Open(htmlFilePath)
	check(err)

	r := bufio.NewReader(htmlFile)

	parsedLinks, err := linkparser.Parse(r)
	if err != nil {
		return nil, err
	}
	return parsedLinks, nil
}
