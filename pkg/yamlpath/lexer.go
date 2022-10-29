package yamlpath

import (
	"fmt"
	"strings"

	"github.com/alecthomas/participle/v2"
	"github.com/mohae/deepcopy"
)

// YamlPath is a string to go through YAML (but internally through golang interface{}).
// It understands these expressions:
// `$deployment`
// `$deployment.key`
// `$deployment.metadata.annotations."app/name"`
// `$deployment.spec.volumes[0]`
type YamlPath struct {
	Root *SubPath `parser:"'$' @@"`
}

func (p *YamlPath) Name() string {
	if p.Root == nil {
		return ""
	}

	return "$" + p.Root.Key
}

func (p *YamlPath) Copy() *YamlPath {
	return deepcopy.Copy(p).(*YamlPath)
}

func (p *YamlPath) AddSub(sub *SubPath) {
	s := p.Root
	for s.Sub != nil {
		s = s.Sub
	}

	s.Sub = sub
}

func (p *YamlPath) String() string {
	parts := []string{}
	i := p.Root

	for i != nil {
		s := i.Key
		if i.Index != nil {
			s += fmt.Sprintf("[%d]", *i.Index)
		}
		parts = append(parts, s)

		i = i.Sub
	}

	return strings.Join(parts, ".")
}

type SubPath struct {
	Key   string   `parser:"@(String|Char|RawString|Ident)"`
	Index *int     `parser:"('[' @Int ']')?"`
	Sub   *SubPath `parser:"('.' @@)?"`
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
