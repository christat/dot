package dot

import (
	"bytes"
	"testing"
)

// This file unit tests all the components of the .dot parser.
// Since we're accessing private functions, the body must remain within the package itself.

func TestStripAllComments(t *testing.T) {
	fileStream := []byte(
		`// strip this comment

		/* this one too */

		keep this line though

		/*and...

			...this one...

			...too*/`)
	// spaces, tabs and newlines are irrelevant as they're always captured by parsing expressions
	result := bytes.Trim(stripAllComments(fileStream), " \t\n")
	expected := bytes.Trim([]byte("keep this line though"), " \t\n")
	if !bytes.Equal(expected, result) {
		t.Errorf("stripAllComments() expected: %v, got: %v", string(expected), string(result))
	}
}

func TestParseGraphType(t *testing.T) {
	g := NewGraph()
	fileStream := []byte("DiGrAph test {...") //tests ignore case flag in regexp
	match, fileStream := parseGraphType(g, fileStream)
	if !match {
		t.Error("parseGraphType() failed to match MiXeDcAsE digraph")
	}
	fileStream = []byte("foo")
	match, _ = parseGraphType(g, fileStream)
	if match {
		t.Errorf("parseGraphType() accepted '%v' as graphType", string(fileStream))
	}
}

func TestParseGraphName(t *testing.T) {
	g := NewGraph()
	fileStream := []byte("9name123 {...")
	match, fileStream := parseGraphName(g, fileStream)
	if !match {
		t.Error("parseGraphName() failed to match alphanumeric name")
	}
	fileStream = []byte("_Inv@l1d! {...")
	match, _ = parseGraphType(g, fileStream)
	if match {
		t.Errorf("parseGraphNamee() accepted '%v' as graphName", string(fileStream))
	}
}

func TestParseBlockBegin(t *testing.T) {
	fileStream := []byte(`
		{
			start [ cost = 3, distance = 7 ] -> [ k = 0.12 ] a1;
		}
	`)
	match, fileStream := parseBlockBegin(fileStream)
	if !match {
		t.Error("parseBlockBegin() didn't match a block begin bracket")
	}
	fileStream = []byte("This is not the byte slice you're looking for")
	match, _ = parseBlockBegin(fileStream)
	if match {
		t.Errorf("parseBlockBegin() matched '%v' as a block begin bracket", string(fileStream))
	}
}

func TestParseVertexName(t *testing.T) {
	g := NewGraph()
	fileStream := []byte(`
		    start [ cost = 3, distance = 7 ] -> [ k = 0.12 ] a1;
	`)
	match, fileStream, _ := parseVertexName(fileStream, g,false, false)
	if !match {
		t.Error("parseVertexName() didn't match a correct vertex name")
	}
	fileStream = []byte("{ this bracket shouldn't be here")
	match, fileStream, name := parseVertexName(fileStream, g,true, false)
	if match {
		t.Errorf("parseVertexName() matched '%v' as a block begin bracket", string(name))
	}
}

func TestParseAttributes(t *testing.T) {
	fileStream := []byte("[\tfoo = 0.12, bar=26, foobar =12.26, quote\t=\"sth\", bool\n=true, string=\ttest ]")
	match, _, attr := parseAttributes(fileStream)
	if !match {
		t.Error("parseAttributes() failed to match a correct attributes section")
	}
	foo := attr["foo"] == 0.12
	bar := attr["bar"] == 26
	foobar := attr["foobar"] == 12.26
	quote := attr["quote"] == "sth"
	boolean := attr["bool"] == true
	str := attr["string"] == "test"
	if !foo || !bar || !foobar || !quote || !boolean || !str {
		t.Error("parseAttributes() failed to set the attributes map correctly")
	}

	fileStream = []byte("[ foo\n=1 bar\t=test ]")
	match, _, _ = parseAttributes(fileStream)
	if match {
		t.Error("parseAttributes() parsed an incorrect attribute map")
	}

	fileStream = []byte("[ 15=false, ]")
	match, _, _ = parseAttributes(fileStream)
	if match {
		t.Error("parseAttributes() parsed an incorrect attribute map")
	}
}

func TestCastAttributeValue(t *testing.T) {
	source := "29.10"
	value := castAttributeValue(source)
	_, isType := value.(float64)
	if !isType {
		t.Errorf("castAttributeValue() of '%v' didn't yield type float64", source)
	}

	source = "2017"
	value = castAttributeValue(source)
	_, isType = value.(int)
	if !isType {
		t.Errorf("castAttributeValue() of '%v' didn't yield type int", source)
	}

	source = "true"
	value = castAttributeValue(source)
	_, isType = value.(bool)
	if !isType {
		t.Errorf("castAttributeValue() of '%v' didn't yield type bool", source)
	}

	source = "life-changer"
	value = castAttributeValue(source)
	_, isType = value.(string)
	if !isType {
		t.Errorf("castAttributeValue() of '%v' didn't yield type string", source)
	}
}

func TestParseVertexAttributes(t *testing.T) {
	g := NewGraph()
	fileStream := []byte("[ a=3.1496, b= false, c =	foo ]")
	match, _ := parseVertexAttributes(fileStream, g, "origin")
	if !match {
		t.Error("parseVertexAttributes() failed to match correct attributes section")
	}
	aCorrect := g.vertexAttributes["origin"]["a"] == 3.1496
	bCorrect := g.vertexAttributes["origin"]["b"] == false
	cCorrect := g.vertexAttributes["origin"]["c"] == "foo"
	if !aCorrect || !bCorrect || !cCorrect {
		t.Error("parseVertexAttributes() failed to set attributes section correctly")
	}
}

func TestParseEdgeType(t *testing.T) {
	fileStream := []byte("--")
	match, _, isDirectional := parseEdgeType(fileStream)
	if !match || isDirectional {
		t.Error("parseEdgeType() failed to match or detect undirected edge (--)")
	}

	fileStream = []byte("->")
	match, _, isDirectional = parseEdgeType(fileStream)
	if !match || !isDirectional {
		t.Error("parseEdgeType() failed to match or detect directed edge (->)")
	}

	fileStream = []byte(">>---->")
	match, _, _ = parseEdgeType(fileStream)
	if match {
		t.Error("parseEdgeType() matches with strings other than edges")
	}
}

func TestParseTargetVertexName(t *testing.T) {
	g := NewGraph()
	fileStream := []byte("target;")
	match, _, targetName := parseTargetVertexName(fileStream, g, "origin", true)
	if !match {
		t.Error("parseTargetVertexName() failed to parse valid name 'target'")
	}
	found := false
	for i := range g.adjacencyMap["origin"] {
		if g.adjacencyMap["origin"][i].(*Vertex).name == targetName {
			found = true
			break
		}
	}

	foundUndirected := false
	for i := range g.adjacencyMap[targetName] {
		if g.adjacencyMap[targetName][i].(*Vertex).name == "origin" {
			foundUndirected = true
			break
		}
	}
	if !found || !foundUndirected {
		t.Error("parseTargetVertexName() failed to get stored in adjacency map")
	}
}
