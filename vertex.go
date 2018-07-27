package dot

import "github.com/christat/search"

type Vertex struct {
	name  string
	graph *Graph
}

// New vertex allows to easily generate a vertex, providing the underlying graph instance and its unique name.
func NewVertex(name string, graph *Graph) *Vertex {
	return &Vertex{name: name, graph: graph}
}

// Name returns the unique identifier of the Vertex, i.e. its name.
func (v *Vertex) Name() string {
	return v.name
}

// Equals implements the search.State interface, comparing two instances of a Vertex by name (by dot standards,
// they should be unique).
func (v *Vertex) Equals(other search.State) bool {
	return v.name == other.(*Vertex).name
}

// Neighbors allows to obtain a map of adjacent vertices to the caller.
func (v *Vertex) Neighbors() (neighbors []search.State) {
	if v.graph.adjacencyMap != nil {
		_, found := v.graph.adjacencyMap[v.name]
		if  !found {
			return make([]search.State, 0)
		}
	}
	return v.graph.adjacencyMap[v.name]
}

const defaultCost = 10e9
const defaultHeuristic = 0


// Cost relies on the underlying graph structure to obtain either a cost function to traverse from v to target,
// or alternatively a cost key if the cost is coded into the graph description. Alternatively, it returns a
// default cost of 10e9 as a measure of caution.
func (v *Vertex) Cost(target search.State) float64 {
	if v.graph.CostFunc != nil {
		return v.graph.CostFunc(v, target.(*Vertex))
	}
	if v.graph.CostKey != "" {
		cost, err := v.graph.GetEdgeAttribute(v.name, target.(*Vertex).name, v.graph.CostKey)
		if err != nil {
			return defaultCost
		}
		floatCost, ok := cost.(float64)
		if !ok {
			intCost, ok := cost.(int)
			if !ok {
				return defaultCost
			}
			return float64(intCost)
		}
		return floatCost
	}
	return defaultCost
}

// Heuristic, similarly to the Cost method, relies on either a function or a key passed as an attribute of the
// underlying graph. As a fallback, a heuristic value of 0 is returned.
func (v *Vertex) Heuristic() float64 {
	if v.graph.HeuristicFunc != nil {
		return v.graph.HeuristicFunc(v)
	}
	if v.graph.HeuristicKey != "" {
		heuristic, err := v.graph.GetVertexAttribute(v.name, v.graph.HeuristicKey)
		if err != nil {
			return defaultHeuristic
		}
		floatHeuristic, ok := heuristic.(float64)
		if !ok {
			intHeuristic, ok := heuristic.(int)
			if !ok {
				return defaultHeuristic
			}
			return float64(intHeuristic)
		}
		return floatHeuristic
	}
	return defaultHeuristic
}
