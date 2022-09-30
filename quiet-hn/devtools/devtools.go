package devtools

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/AksAman/gophercises/quietHN/utils"
	cHTML "github.com/alecthomas/chroma/formatters/html"
	"github.com/alecthomas/chroma/lexers"
	"github.com/alecthomas/chroma/styles"
)

var (
	filePathWithLineRegex = regexp.MustCompile(`(\t+)(.+)\/([^\/]+).go:(\d+)\s(.*)$`)
)

func ParseStackTraceToHTML(stackTrace string) string {
	traceLines := strings.Split(stackTrace, "\n")
	var htmlLines []string

	for _, line := range traceLines {
		if line == "" {
			continue
		}

		matches := filePathWithLineRegex.FindStringSubmatch(line)

		if len(matches) == 0 || len(matches) != 6 {
			htmlLines = append(htmlLines, line)
			continue
		}
		initialTabs := matches[1]
		filePath := strings.TrimSpace(matches[2])
		fileName := strings.TrimSpace(matches[3])
		lineNumber := strings.TrimSpace(matches[4])
		rest := matches[5]

		v := url.Values{}

		completeFilePath := filepath.Join(filePath, fileName+".go")
		v.Set("path", completeFilePath)
		v.Set("line", lineNumber)
		endpoint := fmt.Sprintf("/__internal__/view-source?%s#line%s", v.Encode(), lineNumber)

		href := fmt.Sprintf("%s<a href=\"%s\">%s:%s</a> %s", initialTabs, endpoint, completeFilePath, lineNumber, rest)
		htmlLines = append(htmlLines, href)
		lineNumberInt, err := strconv.Atoi(lineNumber)
		if err != nil {
			htmlLines = append(htmlLines, err.Error())
		}
		srcContents, err := GetHighlightedSourceCode(completeFilePath, lineNumberInt, "monokai", 3)
		if err != nil {
			htmlLines = append(htmlLines, err.Error())
		} else {
			htmlLines = append(htmlLines, srcContents)
		}

	}

	return strings.Join(htmlLines, "\n")
}

func GetHighlightedSourceCode(srcPath string, lineNumber int, theme string, lineCount int) (string, error) {
	file, err := os.Open(srcPath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	lexer := lexers.Get(".go")
	if lexer == nil {
		lexer = lexers.Fallback
	}

	fileContents := ""
	initialLineNumber := 0
	if lineCount > 0 {
		initialLineNumber = lineNumber

		n := 0
		startLine := lineNumber - lineCount
		if startLine < 0 {
			startLine = 0
		}

		endLine := lineNumber + lineCount
		scanner := bufio.NewScanner(file)
		totalLineCount, err := utils.CountLinesInReaderFromFilePath(srcPath)
		if err != nil {
			return "", err
		}

		if endLine > totalLineCount {
			endLine = totalLineCount
		}

		// startLine = 0
		// endLine = totalLineCount

		lineNumber = lineCount + 1

		fmt.Printf("startLine: %v\n", startLine)
		fmt.Printf("endLine: %v\n", endLine)

		fileContentBuilder := strings.Builder{}
		for scanner.Scan() {
			n++
			if n < startLine {
				continue
			}
			if n > endLine {
				break
			}
			fileContentBuilder.WriteString(scanner.Text() + "\n")
		}

		if err := scanner.Err(); err != nil {
			return "", err
		}

		fileContents = fileContentBuilder.String()
	} else {

		fileBuffer := bytes.NewBuffer(nil)
		_, err = io.Copy(fileBuffer, file)
		if err != nil {
			return "", err
		}
		fileContents = fileBuffer.String()

	}

	iterator, err := lexer.Tokenise(nil, fileContents)
	if err != nil {
		return "", err
	}

	style := styles.Get(theme)
	if style == nil {
		style = styles.Fallback
	}

	lineRange := [][2]int{}
	if lineNumber > 0 {
		lineRange = append(lineRange, [2]int{initialLineNumber, initialLineNumber})
	}

	fmt.Println("srcPath", srcPath, "lineRange: ", lineRange)

	formatter := cHTML.New(
		cHTML.BaseLineNumber(initialLineNumber),
		cHTML.WithLineNumbers(true),
		cHTML.LineNumbersInTable(true),
		cHTML.HighlightLines(lineRange),
		cHTML.LinkableLineNumbers(true, "line"),
	)
	buffer := new(bytes.Buffer)
	err = formatter.Format(buffer, style, iterator)
	if err != nil {
		return "", err
	}

	// if html, err := html.Parse(buffer.String()); err == nil {
	// 	// find line number

	// }

	return buffer.String(), nil

}
