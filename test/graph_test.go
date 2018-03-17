package dot_test

import (
	"testing"

	"github.com/christat/dot/graph"
	"errors"
)

func generateGraph() *dot.Graph {
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
	g.EdgeAttributes = map[string]map[string]map[string]interface{}{
		"s": {
			"A": {
				"k": 1,
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
	return g
}

func TestNeighbors(t *testing.T) {
	g := generateGraph()
	neighbors, err := g.Neighbors("C")
	if err != nil {
		t.Error("Neighbors() failed on valid string vertex key")
	}
	if len(neighbors) != 2 || neighbors[0].(string) != "s" || neighbors[1].(string) != "t" {
		t.Error("Neighbors() fetched incorrect map entity")
	}

	neighbors, _ = g.Neighbors("foo")
	if neighbors != nil {
		t.Error("Neighbors did not fail on non-existent vertex key")
	}
}

func TestG(t *testing.T) {
	g := generateGraph()

	// test custom function
	g.CostFunc = func (origin, target interface{}) (float64, error) {
		return 1234, nil
	}
	cost, _ := g.G("s", "A")
	if cost != 1234 {
		t.Error("CostFunc function attribute did not execute in G function")
	}
	g.CostFunc = func (origin, target interface{}) (float64, error) {
		return 0, errors.New("failed on purpose! :)")
	}
	_, err := g.G("s", "A")
	if err == nil {
		t.Error("CostFunc function with forced error did not trigger")
	}
	g.CostFunc = nil

	// test costKeys
	g.CostKey = "h_cff"
	cost, err = g.G("s", "A")
	if err != nil {
		t.Error("G failed to fetch cost for valid CostKey")
	}
	if cost != 2 {
		t.Error("G fetched incorrect cost for CostKey")
	}

	g.CostKey = "invalid"
	_, err = g.G("s", "A")
	if err == nil {
		t.Error("G did not fail on invalid CostKey")
	}

	g.CostKey = "name"
	_, err = g.G("s", "A")
	if err == nil {
		t.Error("G did not fail on non-float CostKey value")
	}

	// test default
	g.CostKey = ""
	cost, _ = g.G("s", "A")
	if cost != 1 {
		t.Error("G without cost function nor key did not return unitary cost")
	}
}

func TestH(t *testing.T) {
	g := generateGraph()

	// test custom function
	g.HeuristicFunc = func (node interface{}) (float64, error) {
		return 1234, nil
	}
	heuristic, _ := g.H("s")
	if heuristic != 1234 {
		t.Error("HeuristicFunc function attribute did not execute in H function")
	}
	g.HeuristicFunc = func (node interface{}) (float64, error) {
		return 0, errors.New("failed on purpose! :)")
	}
	_, err := g.H("s")
	if err == nil {
		t.Error("HeuristicFunc function with forced error did not trigger")
	}
	g.HeuristicFunc = nil

	// test costKeys
	g.HeuristicKey = "h_cff"
	heuristic, err = g.H("s")
	if err != nil {
		t.Error("H failed to fetch heuristic for valid HeuristicKey")
	}
	if heuristic != 2 {
		t.Error("H fetched incorrect heuristic for HeuristicKey")
	}

	g.HeuristicKey = "invalid"
	_, err = g.H("s")
	if err == nil {
		t.Error("H did not fail on invalid HeuristicKey")
	}

	g.HeuristicKey = "name"
	_, err = g.H("s")
	if err == nil {
		t.Error("H did not fail on non-float HeuristicKey value")
	}

	// test default
	g.HeuristicKey = ""
	heuristic, _ = g.H("s")
	if heuristic != 0 {
		t.Error("H without heuristic function nor key did not return zero heuristic")
	}
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
	if value != "2.0" {
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
