package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/christat/dot/parser"
)

const (
	exitError   = 1
	exitSuccess = 0
)

func main() {
	// definition of CLI parameters
	filePath := flag.String("f", "", "path to .dot file containing the graph definition\n")
	inspect := flag.Bool("i", false, "inspection mode. prints the parsed graph's attributes\n")
	verbose := flag.Bool("v", false, "verbose mode. If set, control statements are printed during parsing\n")
	flag.Parse()

	//get CLI program exec name
	programName := os.Args[0]

	// parameter checking
	if !(len(*filePath) > 0) {
		fmt.Fprintf(os.Stderr, "Please, provide a .dot file through argument -file.\nsee %v -help for more details\n",
			programName)
		os.Exit(exitError)
	}

	// run parser
	ok, g := dot.ParseFile(*filePath, *verbose)
	if !ok {
		fmt.Fprintf(os.Stderr, "Failed to parse file %v. Please check file and try verbose mode (-v) to assess any errors.", filePath)
		os.Exit(exitError)
	}

	// if attribute inspection has been selected
	if *inspect {
		vertices := g.AdjacencyMap

		// loop over all vertices of the graph
		for vertex := range vertices {
			attributes, ok := g.VertexAttributes[vertex]
			fmt.Printf(" Vertex %v:\n", vertex)
			if ok {
				for attribute := range attributes {
					attributeValue, ok := g.VertexAttributes[vertex][attribute]
					if !ok {
						fmt.Fprintf(os.Stderr, "Failed to fetch vertex attribute %v", attribute)
						os.Exit(exitError)
					}
					fmt.Printf("\t%v: %v\n", attribute, attributeValue)
				}
				fmt.Println()
			} else {
				fmt.Printf("\t<no attributes>\n\n")
			}
		}

		// loop over all edges of the graph (vertex src -> vertex target)
		for vertex := range vertices {
			// get all immediately adjacent vertices
			neighbors := g.AdjacencyMap[vertex]
			if neighbors == nil {
				fmt.Fprintf(os.Stderr, "Failed to fetch neighbors of vertex %v", vertex)
				os.Exit(exitError)
			}
			for _, neighbor := range neighbors {
				attributes, _ := g.GetEdgeAttributes(vertex, neighbor.(string))
				if len(attributes) > 0 {
					fmt.Printf(" Edge %v -> %v:\n", vertex, neighbor)
					if len(attributes) > 0 {
						for attribute := range attributes {
							value, err := g.GetEdgeAttribute(vertex, neighbor.(string), attribute)
							if err != nil {
								fmt.Fprintf(os.Stderr, "Failed to fetch edge attributes %v", attribute)
								os.Exit(exitError)
							}
							fmt.Printf("\t %v: %v\n", attribute, value)
						}
						fmt.Println()
					}
				}
			}
		}
	}
	os.Exit(exitSuccess)
}
