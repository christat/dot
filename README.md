[![GoDoc](https://godoc.org/github.com/christat/dot?status.svg)](https://godoc.org/github.com/christat/dot)
[![Build Status](https://travis-ci.org/christat/dot.svg?branch=master)](https://travis-ci.org/christat/dot)
# DOT

A graph parsing library for graph definitions using the [.dot specification.](http://www.graphviz.org/doc/info/lang.html)

Includes:
- type **Graph** to represent all the connections and attributes of the graph, along with utility functions to manipulate vertices, edges and attributes for both.
- Two library functions:
    -  **dot.Parse()**: parses a []byte with a .dot graph definition.
    - **dot.ParseFile()**: a wrapper to read an input file and invoke _dot.Parse()_
- An tiny executable to test the parsing functionality. It takes an argument _-f [fileName...]_ and optionally flags _-v_ (verbose mode, prints chain of tokens detected during parsing) and _-i_ (inspection, prints all connections and attributes for vertices and edges).

**Note**: the parser implements a subset of the full specification, with the following limitations:
- HTML strings (<...>) are not allowed in _IDs_.
- Escaped quotes (\") are not allowed in _IDs_.
- The keyword **strict** is not recognized.
- The keyword **subgraph** is not recognized.

## Download/Installation

In your Go project's root directory, open a terminal and paste the following:

```
go get github.com/christat/dot
```

## License

Licensed under the MIT license.