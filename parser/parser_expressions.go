package dot

import (
	"regexp"
)

// Regular expressions used to parse the input files according to the DOT spec.
// For a full description of the DOT language, see: http://www.graphviz.org/doc/info/lang.html
var (
	lineStartRe          = "^"
	whitespaceRe         = "[[:space:]]*"
	matchPerLineRe       = "(?m)"
	caseInsensitiveRe    = "(?i)"
	startWhitespaceRe    = lineStartRe + whitespaceRe
	graphTypeRe          = regexp.MustCompile(startWhitespaceRe + caseInsensitiveRe + "(digraph|Graph)")
	graphNameRe          = regexp.MustCompile(startWhitespaceRe + "([[:alnum:]]+)")
	blockEndRe           = regexp.MustCompile(startWhitespaceRe + "\\};?" + whitespaceRe)
	blockBeginRe         = regexp.MustCompile(startWhitespaceRe + "\\{" + whitespaceRe)
	vertexNameRe         = regexp.MustCompile(startWhitespaceRe + "([[:alnum:]]+)")
	edgeTypeRe           = regexp.MustCompile(startWhitespaceRe + "(--|->)")
	endOfStatementRe     = regexp.MustCompile(startWhitespaceRe + ";" + whitespaceRe)
	commentRe            = regexp.MustCompile(matchPerLineRe + whitespaceRe + "//.*(\n|$)")
	multiLineCommentRe   = regexp.MustCompile(whitespaceRe + "/\\*[^*]*\\*+(?:[^/*][^*]*\\*+)*/(\n|$)?")
	attributeBeginRe     = regexp.MustCompile(startWhitespaceRe + "\\[")
	attributeNameRe      = regexp.MustCompile("(" + startWhitespaceRe + ")(([[:alnum:]]|_)+)([[:space:]]*)=")
	attributeValueNextRe = regexp.MustCompile(startWhitespaceRe + "((([+-]?(\\.[[:digit:]]|[[:digit:]]+(\\.[[:digit:]]*)?))|([[:alnum:]]+)|(_|\"[^\"]*\")+))([[:space:]]*),")
	attributeValueEndRe  = regexp.MustCompile(startWhitespaceRe + "((([+-]?(\\.[[:digit:]]|[[:digit:]]+(\\.[[:digit:]]*)?))|([[:alnum:]]+)|(_|\"[^\"]*\")+))([[:space:]]*)]")
)
