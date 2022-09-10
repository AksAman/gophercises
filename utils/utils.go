package utils

import (
	"fmt"
	"strings"
)

func PrintLine() {
	line := strings.Repeat("-", 50)
	fmt.Println(line)
}
