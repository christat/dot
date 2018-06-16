package dot

import "fmt"

// Graph contains the topology and attributes of a Graph, including name, type, and adjacency map and vertex/edge attributes.
// Additionally, it is the backbone of the Vertex type, which implements the interface search.State from github.com/christat/search.
// This means we can use Graph to perform search with the algorithms provided in the aforementioned library.
type Graph struct {
	Name          string
	Type          string
	CostKey       string
	HeuristicKey  string
	CostFunc      func(origin, target *Vertex) float64
	HeuristicFunc func(vertex *Vertex) float64

	// maps a vertex accessor per each unique vertex name. For internal use only
	vertexMap map[string]*Vertex

	// adjacencyMap stores the adjacent vertices of each vertex. Each entry of the map consists of a list of Vertex pointers.
	adjacencyMap map[string][]*Vertex

	// vertexAttributes stores for every vertex a map of attributes in the form "name": "value".
	vertexAttributes map[string]map[string]interface{}

	// edgeAttributes stores, in a map for every vertex name, another map whose key is the target vertex name and
	// the value is a third map of attributes in the form "name": "value".
	edgeAttributes map[string]map[string]map[string]interface{}
}

// NewGraph creates and returns a pointer to a new Graph.
func NewGraph() (g *Graph) {
	g = new(Graph)
	g.vertexMap = make(map[string]*Vertex)
	g.adjacencyMap = make(map[string][]*Vertex)
	g.vertexAttributes = make(map[string]map[string]interface{})
	g.edgeAttributes = make(map[string]map[string]map[string]interface{})
	return g
}

// Build graph from existing data
func BuildGraph(
	g *Graph,
	vertexMap map[string]*Vertex,
	adjacencyMap map[string][]*Vertex,
	vertexAttributes map[string]map[string]interface{},
	edgeAttributes map[string]map[string]map[string]interface{}) {

	g.vertexMap = vertexMap
	g.adjacencyMap = adjacencyMap
	g.vertexAttributes = vertexAttributes
	g.edgeAttributes = edgeAttributes
}

// AdjacencyMap returns the adjacency map of the graph.
func (g *Graph) AdjacencyMap() map[string][]*Vertex {
	return g.adjacencyMap
}

// GetVertexAttributes allows obtaining the map of attributes for a given vertex.
func (g *Graph) GetVertexAttributes(vertex string) (value map[string]interface{}, err error) {
	attributes, exists := g.vertexAttributes[vertex]
	if !exists {
		return nil, fmt.Errorf("GetVertexAttributes() of vertex %v: vertex has no attributes", vertex)
	}
	return attributes, nil
}

// SetVertexAttribute allows to add an attribute to an existing map of attributes for a given vertex.
func (g *Graph) SetVertexAttribute(vertex string, attribute string, value interface{}) {
	if g.vertexAttributes == nil {
		g.vertexAttributes = make(map[string]map[string]interface{})
	}
	_, exists := g.vertexAttributes[vertex]
	if !exists {
		g.vertexAttributes[vertex] = make(map[string]interface{})
	}
	g.vertexAttributes[vertex][attribute] = value
}

// GetVertexAttributes obtains the desired attribute of vertex. If not found, an error value is returned instead
func (g *Graph) GetVertexAttribute(vertex string, attribute string) (value interface{}, err error) {
	_, exists := g.vertexAttributes[vertex]
	if !exists {
		return nil, fmt.Errorf("GetVertexAttribute() of vertex %v: vertex has no attributes", vertex)
	}
	value, exists = g.vertexAttributes[vertex][attribute]
	if !exists {
		return nil, fmt.Errorf("GetVertexAttribute() of vertex %v: attribute %v not found", vertex, attribute)
	}
	return value, nil
}

// SetEdgeAttributes provides an easy way to set a map of attributes for a specific edge (defined by the vertices origin -> target).
// if isDirectional is false, the same property will be set in both origin -> target and target -> origin.
func (g *Graph) SetEdgeAttributes(origin string, target string, isDirectional bool, edgeAttributes map[string]interface{}) {
	if len(edgeAttributes) > 0 {
		if g.edgeAttributes == nil {
			g.edgeAttributes = make(map[string]map[string]map[string]interface{})
		}
		_, exists := g.edgeAttributes[origin]
		if !exists {
			g.edgeAttributes[origin] = make(map[string]map[string]interface{})
		}
		g.edgeAttributes[origin][target] = edgeAttributes

		if !isDirectional {
			g.SetEdgeAttributes(target, origin, true, edgeAttributes)
		}
	}
}

// GetEdgeAttributes obtains all the attributes of the edge (defined by the vertices origin -> target).
// If the edge is undirected it is assumed that the map will hold the same properties in both directions, making one fetch enough.
func (g *Graph) GetEdgeAttributes(origin string, target string) (map[string]interface{}, error) {
	_, exists := g.edgeAttributes[origin]
	if !exists {
		return nil, fmt.Errorf("GetEdgeAttributes() of edge <%v> : failed to find origin\n", origin)
	}
	attributeMap, exists := g.edgeAttributes[origin][target]
	if !exists {
		return nil, fmt.Errorf("GetEdgeAttributes() of edges %v -> %v : failed to find connection\n", origin, target)
	}
	return attributeMap, nil
}

//  SetEdgeAttribute adds the desired attribute to an edge (defined by the vertices origin -> target)
// If isUndirected is true, the property is set for both directions of the edge.
func (g *Graph) SetEdgeAttribute(origin string, target string, isUndirected bool, attribute string, value interface{}) {
	if g.edgeAttributes == nil {
		g.edgeAttributes = make(map[string]map[string]map[string]interface{})
	}
	_, exists := g.edgeAttributes[origin]
	if !exists {
		g.edgeAttributes[origin] = make(map[string]map[string]interface{})
	}
	_, exists = g.edgeAttributes[origin][target]
	if !exists {
		g.edgeAttributes[origin][target] = make(map[string]interface{})
	}
	g.edgeAttributes[origin][target][attribute] = value

	if isUndirected {
		g.SetEdgeAttribute(target, origin, false, attribute, value)
	}
}

// GetEdgeAttribute obtains the desired attribute of an edge (defined by the vertices origin -> target).
// If the edge is undirected it is assumed that the map will hold the same properties in both directions, making one fetch enough.
func (g *Graph) GetEdgeAttribute(origin string, target string, attribute string) (interface{}, error) {
	_, exists := g.edgeAttributes[origin]
	if !exists {
		return nil, fmt.Errorf("GetEdgeAttribute() of edge <%v> : failed to find origin in map\n", origin)
	}
	_, exists = g.edgeAttributes[origin][target]
	if !exists {
		return nil, fmt.Errorf("GetEdgeAttribute() of edges %v -> %v : failed to find connection to target in map\n", origin, target)
	}
	value, exists := g.edgeAttributes[origin][target][attribute]
	if !exists {
		return nil, fmt.Errorf("GetEdgeAttribute() of edges %v -> %v, attribute %v : failed to find attribute\n",
			origin, target, attribute)
	}
	return value, nil
}

//
func (g *Graph) fetchOrCreateVertex(name string) *Vertex {
	vertex, exists := g.vertexMap[name]
	if !exists {
		vertex = NewVertex(name, g)
		g.vertexMap[name] = vertex
		return vertex
	}
	return vertex
}
