package dot

import (
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strconv"
	"strings"
)

func printToken(str string) {
	if verbose {
		fmt.Printf("[ %v ]\n", str)
	}
}

func readFile(filePath string) (ok bool, fileStream []byte) {
	fileStream, err := ioutil.ReadFile(filePath)
	if err != nil {
		fmt.Fprint(os.Stderr, err)
		return false, nil
	}
	return true, fileStream
}

func stripAllComments(contents []byte) []byte {
	contents = commentRe.ReplaceAllLiteral(contents, []byte(""))
	return multiLineCommentRe.ReplaceAllLiteral(contents, []byte(""))
}

func sliceMatch(contents []byte, regexp regexp.Regexp) (match bool, src []byte, value string) {
	loc := regexp.FindIndex(contents)
	if loc != nil {
		// if match found, slice contents and return sliced value
		value = string(contents[loc[0]:loc[1]])
		contents = contents[loc[1]:]
		return true, contents, value
	}
	return false, contents, ""
}

func parseGraphType(g *Graph, contents []byte) (match bool, src []byte) {
	match, contents, graphType := sliceMatch(contents, *graphTypeRe)
	if !match {
		fmt.Fprintln(os.Stderr, " Syntax error: GRAPH TYPE could not be parsed")
		return false, contents
	}
	graphType = strings.Trim(graphType, " \t\n")
	g.Type = strings.ToLower(graphType)
	printToken("TYPE " + graphType)
	return true, contents
}

func parseGraphName(g *Graph, contents []byte) (match bool, src []byte) {
	match, contents, graphName := sliceMatch(contents, *graphNameRe)
	if !match {
		fmt.Fprint(os.Stderr, " Syntax error: GRAPH NAME could not be parsed\n")
		return false, contents
	}
	graphName = strings.Trim(graphName, " \t\n")
	g.Name = graphName
	printToken("NAME " + graphName)
	return true, contents
}

func parseBlockBegin(contents []byte) (match bool, src []byte) {
	match, contents, _ = sliceMatch(contents, *blockBeginRe)
	if !match {
		fmt.Fprintln(os.Stderr, " Syntax error: BLOCK BEGIN missing")
		return false, contents
	}
	printToken("--- BLOCK BEGIN found ---")
	return true, contents
}

func parseVertexName(contents []byte, matchOptional bool, isTarget bool) (match bool, src []byte, name string) {
	match, contents, name = sliceMatch(contents, *vertexNameRe)
	name = strings.Trim(name, " \t\n")
	if !match {
		if !matchOptional {
			fmt.Fprintln(os.Stderr, " Syntax error: VERTEX NAME could not be parsed")
		}
		return false, contents, ""
	}
	if isTarget {
		printToken("TARGET VERTEX NAME " + name)
	} else {
		printToken("VERTEX NAME " + name)
	}
	return true, contents, name
}

func parseAttributes(contents []byte) (match bool, src []byte, attributes map[string]interface{}) {
	match, contents, _ = sliceMatch(contents, *attributeBeginRe)
	if match {
		attributeMap := make(map[string]interface{})
		attributesEnd := false
		for !attributesEnd { // loop over attributes
			var attributeName, attributeString string
			match, contents, attributeName = sliceMatch(contents, *attributeNameRe)
			attributeName = strings.Trim(attributeName, " =\t\n")
			if !match {
				fmt.Fprintln(os.Stderr, " Syntax error: ATTRIBUTE section expected attribute name")
				return false, src, nil
			}
			printToken("\tATTRIBUTE " + attributeName)

			// There's two execution branches:
			match, contents, attributeString = sliceMatch(contents, *attributeValueEndRe)
			if match {
				// 1) value followed by end of section (']')
				attributeString = strings.Trim(attributeString, " \"\\]\t\n")
				printToken("\tVALUE " + attributeString)
				attributeMap[attributeName] = castAttributeValue(attributeString)
				attributesEnd = true
			} else {
				// 2) value followed by comma
				match, contents, attributeString = sliceMatch(contents, *attributeValueNextRe)
				attributeString = strings.Trim(attributeString, " ,\"\t\n")
				if !match {
					fmt.Fprintln(os.Stderr, " Syntax error: ATTRIBUTE section is neither ended nor continued")
					return false, src, nil
				}
				printToken("\tVALUE " + attributeString)
				attributeMap[attributeName] = castAttributeValue(attributeString)
			}
		}
		return len(attributeMap) > 0, contents, attributeMap
	}
	return false, contents, nil
}

func castAttributeValue(value string) (castedValue interface{}) {
	var err error
	if strings.ContainsRune(value, '.') {
		castedValue, err = strconv.ParseFloat(value, 64)
	} else {
		castedValue, err = strconv.Atoi(value)
	}
	if err != nil {
		castedValue, err = strconv.ParseBool(value)
		if err != nil {
			return value
		}
	}
	return castedValue
}

func parseVertexAttributes(contents []byte, g *Graph, vertexName string) (match bool, src []byte) {
	match, contents, vertexAttributes := parseAttributes(contents)
	if match {
		g.vertexAttributes[vertexName] = vertexAttributes
	}
	return match, contents
}

func parseEdgeType(contents []byte) (match bool, src []byte, isUndirected bool) {
	isUndirected = false
	match, contents, edgeType := sliceMatch(contents, *edgeTypeRe)
	edgeType = strings.Trim(edgeType, " \t\n")
	if !match {
		fmt.Fprintln(os.Stderr, " Syntax error: EDGE TYPE could not be parsed")
		return false, contents, isUndirected
	}
	printToken("EDGE TYPE " + edgeType)
	if edgeType == "--" {
		isUndirected = true
	}
	return true, contents, isUndirected
}

func parseTargetVertexName(contents []byte, g *Graph, sourceVertex string, isUndirected bool) (match bool, src []byte, targetVertex string) {
	match, contents, targetVertex = parseVertexName(contents, true, true)
	if match {
		g.adjacencyMap[sourceVertex] = append(g.adjacencyMap[sourceVertex], g.fetchOrCreateVertex(targetVertex))
		if isUndirected {
			g.adjacencyMap[targetVertex] = append(g.adjacencyMap[targetVertex], g.fetchOrCreateVertex(sourceVertex))
		}
		return true, contents, targetVertex
	}
	return false, contents, ""
}
