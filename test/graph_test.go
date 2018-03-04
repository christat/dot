package dot_test

import (
	"github.com/christat/dot/parser"
	"testing"
)

func generateGraph() (*dot.Graph) {
	g := dot.NewGraph()
	g.Name = "example"
	g.Type = "digraph"
	g.AdjacencyMap = map[string][]interface{}{
		"s": {"A", "C"},
		"A": {"s", "t"},
		"C": {"s", "t"},
		"t": {},
	}
	g.VertexAttributes = map[string]map[string]interface{}{
		"s": {
			"h_cff": 2,
			"h_pdb": 2,
			"name": "Start",
		},
		"A": {
			"h_cff": 10.1,
		},
		"C": {
			"h_ff": 1,
			"h_pdb": 10,
			"h_cff": 3.14159,
		},
	}
	g.EdgeAttributes = map[string]map[string]map[string]interface{}{
		"s": {
			"A": {
				"k": 1,
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
	return g
}

func TestSetVertexAttribute(t *testing.T) {
	g := new(dot.Graph)
	g.SetVertexAttribute("foo", "bar", "foobar")
	value, ok := g.VertexAttributes["foo"]["bar"]
	if !ok {
		t.Error("SetVertexAttribute() failed to store a value in the graph")
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
	if value != 2 {
		t.Error("GetVertexAttribute() fetched an incorrect attribute value")
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
	g.SetEdgeAttributes("origin", "target", true, edgeAttributes)
	value, ok := g.EdgeAttributes["origin"]["target"]
	if !ok {
		t.Error("SetEdgeAttributes() failed to set an attribute map")
	}
	if value["a"] != 1 || value["b"] != true || value["c"] != "foo" {
		t.Error("SetEdgeAttributes() failed to set attrubute values correctly")
	}
	value, ok = g.EdgeAttributes["target"]["origin"]
	if !ok {
		t.Error("SetEdgeAttributes() failed to set an undirected attribute map")
	}
	if value["a"] != 1 || value["b"] != true || value["c"] != "foo" {
		t.Error("SetEdgeAttributes() failed to set undirected attrubute values correctly")
	}

	g = dot.NewGraph()
	g.SetEdgeAttributes("origin", "target", false, edgeAttributes)
	value, ok = g.EdgeAttributes["origin"]["target"]
	if !ok {
		t.Error("SetEdgeAttributes() failed to set an attribute map")
	}
	if value["a"] != 1 || value["b"] != true || value["c"] != "foo" {
		t.Error("SetEdgeAttributes() failed to set attrubute values correctly")
	}
	value, ok = g.EdgeAttributes["target"]["origin"]
	if ok {
		t.Error("SetEdgeAttributes() set undirected edge properties when the undirection was set false")
	}
}

func TestGetEdgeAttributes(t *testing.T) {
	g := generateGraph()
	attributeMap, err := g.GetEdgeAttributes("C", "t")
	if err != nil {
		t.Error("GetEdgeAttributes() failed to fetch existing attribute map")
	}
	if attributeMap["k"] != 2 {
		t.Error("GetEdgeAttributes() fetched attributes don't match with contents")
	}

	attributeMap, err = g.GetEdgeAttributes("C", "bar")
	if err == nil {
		t.Error("GetEdgeAttributes fetched a non-existent target edge connection")
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
	value, ok := g.EdgeAttributes["target"]["origin"]["foo"]
	if !ok {
		t.Error("SetEdgeAttribute() failed to set undirected string attribute 'foo'")
	}
	if value != 13.14 {
		t.Error("SetEdgeAttribute() failed to set correct value for undirected string attribute 'foo'")
	}
	value, ok = g.EdgeAttributes["origin"]["target"]["bar"]
	if !ok {
		t.Error("SetEdgeAttribute() failed to set bool attribute 'bar'")
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
	}
	if attribute != 2 {
		t.Error("GetEdgeAttributes() fetched attributes don't match with contents")
	}

	_, err = g.GetEdgeAttribute("C", "t", "foo")
	if err == nil {
		t.Error("GetEdgeAttributes fetched a non-existent attribute value")
	}
}