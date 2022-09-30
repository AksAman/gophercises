package utils

import (
	"bufio"
	"bytes"
	"io"
	"os"
)

func CountLinesInReader(f *os.File) (int, error) {
	r := bufio.NewReader(f)
	buf := make([]byte, 32*1024)
	count := 0
	lineSep := []byte{'\n'}

	for {
		c, err := r.Read(buf)
		count += bytes.Count(buf[:c], lineSep)

		switch {
		case err == io.EOF:
			return count, nil

		case err != nil:
			return count, err
		}
	}
}

func CountLinesInReaderFromFilePath(srcPath string) (int, error) {
	f, err := os.Open(srcPath)
	if err != nil {
		return 0, err
	}
	defer f.Close()

	return CountLinesInReader(f)
}
