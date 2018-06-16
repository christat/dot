package dot_test

import (
	"path/filepath"
	"strconv"
	"strings"
	"testing"

	"github.com/christat/dot"
)

// TestParseFile takes example files as input to be parsed.
// It calls dot.Parse internally to parse the resulting []byte from reading files
func TestParseFile(t *testing.T) {
	const filePathTemplate = "./test_files/%v.dot"

	for i := range make([]int, 5) {
		filePath := strings.Replace(filePathTemplate, "%v", "graph"+strconv.Itoa(i+1), 1)
		filePath, _ = filepath.Abs(filePath)
		ok, _ := dot.ParseFile(filePath, false)
		if !ok {
			t.Errorf("Failed to parse test file graph%v.dot", i+1)
		}
	}

	filePath := strings.Replace(filePathTemplate, "%v", "kanagawa", 1)
	filePath, _ = filepath.Abs(filePath)
	ok, _ := dot.ParseFile(filePath, false)
	if ok {
		t.Error("Parsed a graph inside kanagawa O_o")
	}
}
