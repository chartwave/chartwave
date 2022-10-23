package yamlpath

import (
	"github.com/alecthomas/participle/v2"
)

// YamlPath is a string to go through YAML (but internally through golang interface{}).
// It understands these expressions:
// `$deployment`
// `$deployment.key`
// `$deployment.metadata.annotations."app/name"`
// `$deployment.spec.volumes[0]`
// `$deployment.spec.volumes[0][1]`
type YamlPath struct {
	Root *path `parser:"'$' @@"`
}

type path struct {
	Key     string `parser:"@(String|Char|RawString|Ident)"`
	Index   []int  `parser:"('[' @Int ']')*"`
	SubPath *path  `parser:"('.' @@)?"`
}

var (
	parser *participle.Parser[YamlPath]
)

func init() {
	parser = getParser()
}

func getParser() *participle.Parser[YamlPath] {
	return participle.MustBuild[YamlPath]()
}

func ParsePath(input string) (*YamlPath, error) {
	return parser.ParseString("input", input)
}
