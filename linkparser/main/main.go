package main

import (
	"os"
	"path/filepath"
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

	parseLinks(filepath.Join(htmlSrcDir, "ex1.html"))

}

func parseLinks(htmlFilePath string) {
	htmlBytes := os.ReadFile(htmlFilePath)
}
