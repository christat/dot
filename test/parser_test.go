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
	const filePathTemplate = "./dot_files/%v.dot"

	for i := range make([]int, 5) {
		filePath := strings.Replace(filePathTemplate, "%v", "graph"+strconv.Itoa(i+1), 1)
		filePath, _ = filepath.Abs(filePath)
		ok, _ := dot.ParseFile(filePath)
		if !ok {
			t.Errorf("Failed to parse test file graph%v.dot", i+1)
		}
	}

	filePath := strings.Replace(filePathTemplate, "%v", "kanagawa", 1)
	filePath, _ = filepath.Abs(filePath)
	ok, _ := dot.ParseFile(filePath)
	if ok {
		t.Error("Parsed a graph inside kanagawa O_o")
	}
}

func TestInspectParsedFile(t *testing.T) {
	undirName := "cyclic_undirected_graph.dot"
	filePath, _ := filepath.Abs("./dot_files/" + undirName)
	ok, graph := dot.ParseFile(filePath)
	if !ok {
		t.Errorf("Failed to parse test file %v", undirName)
	}

	value, err := graph.GetEdgeAttribute("a", "b", "w")
	if err != nil {
		t.Errorf("Failed to fetch existing edge attribute")
	}
	if value != 7 {
		t.Errorf("Value contained in edge attribute is invalid")
	}
	value, err = graph.GetEdgeAttribute("b", "a", "w")
	if err != nil {
		t.Errorf("Failed to fetch edge attribute in undirected edge (set implicitly)")
	}
	if value != 7 {
		t.Errorf("Value contained in inverse direction of edge is invalid")
	}
}