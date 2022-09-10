package utils

import (
	"fmt"
	"testing"
)

func compare(t *testing.T, got, want any) {
	if got != want {
		t.Errorf("got: %v, want: %v", got, want)
	}
}

func TestDoesFileExists(t *testing.T) {
	testCases := []struct {
		filename string
		exists   bool
	}{
		{"unknownfile.go", false},
		{"path_utils.go", true},
	}

	for _, testCase := range testCases {
		t.Run(
			fmt.Sprintf("Testing %v", testCase.filename),
			func(t *testing.T) {
				want := testCase.exists
				got := DoesFileExists(testCase.filename)
				compare(t, got, want)
			},
		)
	}

}
