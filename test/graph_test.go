package dot_test

import (
	"testing"

	"github.com/christat/dot"
	"github.com/christat/search"
)

func generateGraph() *dot.Graph {
	g := dot.NewGraph()
	g.Name = "example"
	g.Type = "digraph"

	A := dot.NewVertex("A", g)
	C := dot.NewVertex("C", g)
	s := dot.NewVertex("s", g)
	t := dot.NewVertex("t", g)

	vertexMap := map[string]*dot.Vertex{
		"A": A,
		"C": C,
		"s": s,
		"t": t,
	}
	adjacencyMap := map[string][]search.State{
		"s": {A, C},
		"A": {s, t},
		"C": {s, t},
		"t": {},
	}
	vertexAttributes := map[string]map[string]interface{}{
		"s": {
			"h_cff": "2.0",
			"h_pdb": 2,
			"name":  "Start",
		},
		"A": {
			"h_cff": 10.1,
		},
		"C": {
			"h_ff":  1,
			"h_pdb": 10,
			"h_cff": 3.14159,
		},
	}
	edgeAttributes := map[string]map[string]map[string]interface{}{
		"s": {
			"A": {
				"k":     1,
				"h_cff": "2.0",
			},
			"C": {
				"k": 1,
			},
		},
		"A": {
			"s": {
				"k": 1,
			},
			"t": {
				"k": 3,
			},
		},
		"C": {
			"s": {
				"k": 1,
			},
			"t": {
				"k": 2,
			},
		},
	}
	dot.BuildGraph(g, vertexMap, adjacencyMap, vertexAttributes, edgeAttributes)
	return g
}

func TestSetVertexAttribute(t *testing.T) {
	g := new(dot.Graph)
	g.SetVertexAttribute("foo", "bar", "foobar")
	value, err := g.GetVertexAttribute("foo", "bar")
	if err != nil {
		t.Error("SetVertexAttribute() failed to store a value in the graph")
		return
	}
	if value != "foobar" {
		t.Error("SetVertexAttribute() set a value incorrectly in the graph")
	}
}

func TestGetVertexAttribute(t *testing.T) {
	g := generateGraph()
	value, err := g.GetVertexAttribute("s", "h_cff")
	if err != nil {
		t.Error(err)
	}
	if value != "2.0" {
		t.Error("GetVertexAttribute() fetched an incorrect attribute value")
		return
	}

	_, err = g.GetVertexAttribute("foo", "bar")
	if err == nil {
		t.Error("GetVertexAttribute() there's no way this could fail, or is there?")
	}
}

func TestSetEdgeAttributes(t *testing.T) {
	g := dot.NewGraph()
	edgeAttributes := map[string]interface{}{
		"a": 1,
		"b": true,
		"c": "foo",
	}
	g.SetEdgeAttributes("origin", "target", false, edgeAttributes)
	value, err := g.GetEdgeAttributes("origin", "target")
	if err != nil {
		t.Error("SetEdgeAttributes() failed to set an attribute map")
		return
	}
	if value["a"] != 1 || value["b"] != true || value["c"] != "foo" {
		t.Error("SetEdgeAttributes() failed to set attrubute values correctly")
		return
	}
	value, err = g.GetEdgeAttributes("target", "origin")
	if err != nil {
		t.Error("SetEdgeAttributes() failed to set an undirected attribute map")
		return
	}
	if value["a"] != 1 || value["b"] != true || value["c"] != "foo" {
		t.Error("SetEdgeAttributes() failed to set undirected attrubute values correctly")
		return
	}

	g = dot.NewGraph()
	g.SetEdgeAttributes("origin", "target", true, edgeAttributes)
	value, err = g.GetEdgeAttributes("origin", "target")
	if err != nil {
		t.Error("SetEdgeAttributes() failed to set an attribute map")
		return
	}
	if value["a"] != 1 || value["b"] != true || value["c"] != "foo" {
		t.Error("SetEdgeAttributes() failed to set attrubute values correctly")
		return
	}
	value, err = g.GetEdgeAttributes("target", "origin")
	if err == nil {
		t.Error("SetEdgeAttributes() set undirected edge properties when edge was directional")
	}
}

func TestGetEdgeAttributes(t *testing.T) {
	g := generateGraph()
	attributeMap, err := g.GetEdgeAttributes("C", "t")
	if err != nil {
		t.Error("GetEdgeAttributes() failed to fetch existing attribute map")
		return
	}
	if attributeMap["k"] != 2 {
		t.Error("GetEdgeAttributes() fetched attributes don't match with contents")
		return
	}

	attributeMap, err = g.GetEdgeAttributes("C", "bar")
	if err == nil {
		t.Error("GetEdgeAttributes fetched a non-existent target edge connection")
		return
	}
	attributeMap, err = g.GetEdgeAttributes("foo", "bar")
	if err == nil {
		t.Error("GetEdgeAttributes fetched a non-existent edge")
	}
}

func TestSetEdgeAttribute(t *testing.T) {
	g := dot.NewGraph()
	g.SetEdgeAttribute("origin", "target", true, "foo", 13.14)
	g.SetEdgeAttribute("origin", "target", false, "bar", false)
	value, err := g.GetEdgeAttribute("target", "origin", "foo")
	if err != nil {
		t.Error("SetEdgeAttribute() failed to set undirected string attribute 'foo'")
		return
	}
	if value != 13.14 {
		t.Error("SetEdgeAttribute() failed to set correct value for undirected string attribute 'foo'")
		return
	}
	value, err = g.GetEdgeAttribute("origin", "target", "bar")
	if err != nil {
		t.Error("SetEdgeAttribute() failed to set bool attribute 'bar'")
		return
	}
	if value != false {
		t.Error("SetEdgeAttribute() failed to set correct value for bool attribute 'bar'")
	}
}

func TestGetEdgeAttribute(t *testing.T) {
	g := generateGraph()
	attribute, err := g.GetEdgeAttribute("C", "t", "k")
	if err != nil {
		t.Error("GetEdgeAttributes() failed to fetch existing attribute map")
		return
	}
	if attribute != 2 {
		t.Error("GetEdgeAttributes() fetched attributes don't match with contents")
		return
	}

	_, err = g.GetEdgeAttribute("C", "t", "foo")
	if err == nil {
		t.Error("GetEdgeAttributes fetched a non-existent attribute value")
	}
}
