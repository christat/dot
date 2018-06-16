package dot

type Vertex struct {
	name  string
	graph *Graph
}

func NewVertex(name string, graph *Graph) *Vertex {
	return &Vertex{name: name, graph: graph}
}

func (v *Vertex) Name() string {
	return v.name
}

func (v *Vertex) Equals(other *Vertex) bool {
	return v.name == other.name
}

func (v *Vertex) Neighbors() (neighbors []*Vertex) {
	return v.graph.adjacencyMap[v.name]
}

const defaultCost = 10e9
const defaultHeuristic = 0

func (v *Vertex) Cost(target *Vertex) float64 {
	if v.graph.CostFunc != nil {
		return v.graph.CostFunc(v, target)
	}
	if v.graph.CostKey != "" {
		cost, err := v.graph.GetEdgeAttribute(v.name, target.name, v.graph.CostKey)
		if err != nil {
			return defaultCost
		}
		floatCost, ok := cost.(float64)
		if !ok {
			return defaultCost
		}
		return floatCost
	}
	return defaultCost
}

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
			return defaultHeuristic
		}
		return floatHeuristic
	}
	return defaultHeuristic
}
