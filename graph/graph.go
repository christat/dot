package dot

import (
	"fmt"
	"strconv"
)

// Graph contains the topology and attributes of a Graph, including name, type, and adjacency map and vertex/edge attributes.
// Additionally, it implements the Domain interface defined in github.com/christat/search. This means we can use this
// type to perform search with the algorithms provided in the aforementioned library.
type Graph struct {
	Name string
	Type string

	// AdjacencyMap stores the adjacent vertices of each vertex. Each entry of the map consists of a list of vertices.
	AdjacencyMap map[string][]interface{}

	// VertexAttributes stores for every vertex a map of attributes in the form "name": "value".
	VertexAttributes map[string]map[string]interface{}

	// EdgeAttributes stores, in a map for every vertex name, another map whose key is the target vertex name and
	// the value is a third map of attributes in the form "name": "value".
	EdgeAttributes map[string]map[string]map[string]interface{}

	CostFunc      func(origin, target interface{}) (cost float64, err error) // If set, called internally by G
	HeuristicFunc func(node interface{}) (cost float64, err error) // If set, called internally by H

	CostKey      string // If set (and no CostFunc func is set), assigns cost given by this key in node of G
	HeuristicKey string // If set (and no HeuristicFunc func is set), assigns, heuristic given by this key in node of H
}

/*
	The following three functions implement the search.Domain interface.
*/

func (g *Graph) Neighbors(node interface{}) (neighbors []interface{}, err error) {
	return g.AdjacencyMap[node.(string)], nil
}

func (g *Graph) G(origin, target interface{}) (cost float64, err error) {
	// if function override
	if g.CostFunc != nil {
		return g.CostFunc(origin, target)
	}
	// if cost key set
	if g.CostKey != "" {
		value, err := g.GetEdgeAttribute(origin.(string), target.(string), g.CostKey)
		if err != nil {
			return 0, err
		}
		return domainTypeSwitch(value)
	}
	// fallback unit cost
	return 1, nil
}

func (g *Graph) H(node interface{}) (heuristic float64, err error) {
	// if function override
	if g.HeuristicFunc != nil {
		return g.HeuristicFunc(node)
	}
	// if heuristic key set
	if g.HeuristicKey != "" {
		value, err := g.GetVertexAttribute(node.(string), g.HeuristicKey)
		if err != nil {
			return 0, err
		}
		return domainTypeSwitch(value)
	}
	// fallback without heuristic
	return 0, nil
}

func domainTypeSwitch(i interface{}) (value float64, err error) {
	switch v := i.(type) {
	case float64:
		return v, nil
	case int:
		return float64(v), nil
	case string:
		value, err = strconv.ParseFloat(i.(string), 64)
		if err != nil {
			return 0, err
		}
		return value, nil
	default:
		return 0, fmt.Errorf("cannot convert '%v' to float64", i)
	}
}

// NewGraph creates and returns a pointer to a new Graph
func NewGraph() (g *Graph) {
	g = new(Graph)
	g.AdjacencyMap = make(map[string][]interface{})
	g.VertexAttributes = make(map[string]map[string]interface{})
	g.EdgeAttributes = make(map[string]map[string]map[string]interface{})
	return g
}

// SetVertexAttribute allows to add an attribute to an existing map of attributes for a given vertex.
func (g *Graph) SetVertexAttribute(vertex string, attribute string, value interface{}) {
	if g.VertexAttributes == nil {
		g.VertexAttributes = make(map[string]map[string]interface{})
	}
	_, exists := g.VertexAttributes[vertex]
	if !exists {
		g.VertexAttributes[vertex] = make(map[string]interface{})
	}
	g.VertexAttributes[vertex][attribute] = value
}

// GetVertexAttributes obtains the desired attribute of vertex. If not found, an error value is returned instead
func (g *Graph) GetVertexAttribute(vertex string, attribute string) (value interface{}, err error) {
	_, exists := g.VertexAttributes[vertex]
	if !exists {
		return nil, fmt.Errorf("GetVertexAttribute() of vertex %v: vertex has no attributes", vertex)
	}
	value, exists = g.VertexAttributes[vertex][attribute]
	if !exists {
		return nil, fmt.Errorf("GetVertexAttribute() of vertex %v: attribute %v not found", vertex, attribute)
	}
	return value, nil
}

// SetEdgeAttributes provides an easy way to set a map of attributes for a specific edge (defined by the vertices origin -> target).
// if isUndirected is true, the same property will be set in both origin -> target and target -> origin.
func (g *Graph) SetEdgeAttributes(origin string, target string, isUndirected bool, edgeAttributes map[string]interface{}) {
	if len(edgeAttributes) > 0 {
		if g.EdgeAttributes == nil {
			g.EdgeAttributes = make(map[string]map[string]map[string]interface{})
		}
		_, exists := g.EdgeAttributes[origin]
		if !exists {
			g.EdgeAttributes[origin] = make(map[string]map[string]interface{})
		}
		g.EdgeAttributes[origin][target] = edgeAttributes

		if isUndirected {
			g.SetEdgeAttributes(target, origin, false, edgeAttributes)
		}
	}
}

// GetEdgeAttributes obtains all the attributes of the edge (defined by the vertices origin -> target).
// If the edge is undirected it is assumed that the map will hold the same properties in both directions, making one fetch enough.
func (g *Graph) GetEdgeAttributes(origin string, target string) (map[string]interface{}, error) {
	_, exists := g.EdgeAttributes[origin]
	if !exists {
		return nil, fmt.Errorf("GetEdgeAttributes() of edge <%v> : failed to find origin\n", origin)
	}
	attributeMap, exists := g.EdgeAttributes[origin][target]
	if !exists {
		return nil, fmt.Errorf("GetEdgeAttributes() of edges %v -> %v : failed to find connection\n", origin, target)
	}
	return attributeMap, nil
}

//  SetEdgeAttribute adds the desired attribute to an edge (defined by the vertices origin -> target)
// If isUndirected is true, the property is set for both directions of the edge.
func (g *Graph) SetEdgeAttribute(origin string, target string, isUndirected bool, attribute string, value interface{}) {
	if g.EdgeAttributes == nil {
		g.EdgeAttributes = make(map[string]map[string]map[string]interface{})
	}
	_, exists := g.EdgeAttributes[origin]
	if !exists {
		g.EdgeAttributes[origin] = make(map[string]map[string]interface{})
	}
	_, exists = g.EdgeAttributes[origin][target]
	if !exists {
		g.EdgeAttributes[origin][target] = make(map[string]interface{})
	}
	g.EdgeAttributes[origin][target][attribute] = value

	if isUndirected {
		g.SetEdgeAttribute(target, origin, false, attribute, value)
	}
}

// GetEdgeAttribute obtains the desired attribute of an edge (defined by the vertices origin -> target).
// If the edge is undirected it is assumed that the map will hold the same properties in both directions, making one fetch enough.
func (g *Graph) GetEdgeAttribute(origin string, target string, attribute string) (interface{}, error) {
	_, exists := g.EdgeAttributes[origin]
	if !exists {
		return nil, fmt.Errorf("GetEdgeAttribute() of edge <%v> : failed to find origin in map\n", origin)
	}
	_, exists = g.EdgeAttributes[origin][target]
	if !exists {
		return nil, fmt.Errorf("GetEdgeAttribute() of edges %v -> %v : failed to find connection to target in map\n", origin, target)
	}
	value, exists := g.EdgeAttributes[origin][target][attribute]
	if !exists {
		return nil, fmt.Errorf("GetEdgeAttribute() of edges %v -> %v, attribute %v : failed to find attribute\n",
			origin, target, attribute)
	}
	return value, nil
}
