// Package dot provides a .dot parser implementation and a graph type
package dot

import (
	"fmt"
	"os"
)

var verbose = false

// Parse parses the fileStream, building a Graph instance or returning false otherwise.
func Parse(fileStream []byte, verboseFlag bool) (bool, *Graph) {
	verbose = verboseFlag
	g := NewGraph()

	fileStream = stripAllComments(fileStream)
	match, fileStream := parseGraphType(g, fileStream)
	if !match {
		return false, nil
	}

	match, fileStream = parseGraphName(g, fileStream)
	if !match {
		return false, nil
	}

	match, fileStream = parseBlockBegin(fileStream)
	if !match {
		return false, nil
	}

	// main vertex definitions block
	blockEnd := false
	for !blockEnd {
		blockEnd, fileStream, _ = sliceMatch(fileStream, *blockEndRe)
		if !blockEnd {
			var sourceVertex string
			match, fileStream, sourceVertex = parseVertexName(fileStream, false, false)
			if !match {
				return false, nil
			}

			_, fileStream = parseVertexAttributes(fileStream, g, sourceVertex)

			var edgeIsUndirected bool
			match, fileStream, edgeIsUndirected = parseEdgeType(fileStream)
			if !match {
				return false, nil
			}

			// fetch edge attributes. They cannot be set until we know the target vertices structure.
			var edgeAttributes map[string]interface{}
			match, fileStream, edgeAttributes = parseAttributes(fileStream)

			// target vertices can be specified either inline or in a nested block
			var targetVertex string
			match, fileStream, targetVertex = parseTargetVertexName(fileStream, g, sourceVertex, edgeIsUndirected)
			if match {
				// Inline target vertex
				g.SetEdgeAttributes(sourceVertex, targetVertex, edgeIsUndirected, edgeAttributes)
			} else {
				// Nested block with multiple targets
				match, fileStream, _ = sliceMatch(fileStream, *blockBeginRe)
				if !match {
					fmt.Fprintln(os.Stderr, " Syntax error: TARGET NAME could not be parsed")
					return false, nil
				}

				printToken(" --- Beginning multiple target specification ---")
				targetBlockEnd := false
				targetVertices := make([]string, 10)
				for !targetBlockEnd {
					// slice block end statement
					targetBlockEnd, fileStream, _ = sliceMatch(fileStream, *blockEndRe)
					if !targetBlockEnd {
						match, fileStream, targetVertex = parseTargetVertexName(fileStream, g, sourceVertex, edgeIsUndirected)
						if !match {
							fmt.Fprintln(os.Stderr, " Syntax error: TARGET NAME could not be parsed")
							return false, nil
						}
						targetVertices = append(targetVertices, targetVertex)
						_, fileStream = parseVertexAttributes(fileStream, g, targetVertex)
					} else {
						// apply edgeAttributes to all target vertices
						for _, vertex := range targetVertices {
							g.SetEdgeAttributes(sourceVertex, vertex, edgeIsUndirected, edgeAttributes)
						}
					}
				}
				printToken(" --- Ending multiple target specification ---")
			}
			_, fileStream = parseVertexAttributes(fileStream, g, targetVertex)
			// slice semicolons, if any
			match, fileStream, _ = sliceMatch(fileStream, *endOfStatementRe)
		}
	}
	printToken("--- BLOCK END found ---")
	if verbose {
		fmt.Println()
	}
	return true, g
}

// ParseFile wraps the Parse() function with a file reader to get a fileStream ([]byte) if the file exists.
// Returns a pointer to a Graph instance or false if reading the file or parsing failed.
func ParseFile(filePath string, verbose bool) (bool, *Graph) {
	ok, fileStream := readFile(filePath)
	if !ok {
		fmt.Fprintf(os.Stderr, "Failed to read file %v. Parsing aborted.", filePath)
	}
	return Parse(fileStream, verbose)
}
