package builder

import (
	"fmt"
	"io"
	"strings"

	"github.com/moby/buildkit/frontend/dockerfile/parser"
)

type Chartfile struct {
	Targets map[string]*ChartTarget
}

type ChartTarget struct {
	Base     string
	Alias    string
	Commands []ChartCommand
}

func ParseChartfile(input io.Reader) (*Chartfile, error) {
	mobyAST, err := parser.Parse(input)
	if err != nil {
		return nil, fmt.Errorf("failed to parse Chartfile into AST: %w", err)
	}

	result := &Chartfile{
		Targets: make(map[string]*ChartTarget),
	}
	var target *ChartTarget
	for _, n := range mobyAST.AST.Children {
		switch strings.ToUpper(n.Value) {
		case "FROM":
			if n.Next == nil {
				return nil, fmt.Errorf("failed to parse Chartfile: need to specify base chart in FROM command")
			}
			alias := ""
			target = &ChartTarget{
				Base: n.Next.Value,
			}
			if n.Next.Next != nil {
				if strings.ToUpper(n.Next.Next.Value) == "AS" {
					if n.Next.Next.Next == nil {
						return nil, fmt.Errorf("failed to parse Chartfile: need to specify target alias in 'FROM base AS alias' command")
					}
					alias = n.Next.Next.Next.Value
				} else {
					return nil, fmt.Errorf("failed to parse Chartfile: unknown modifier for 'FROM' command: %s", n.Next.Next.Value)
				}
			}

			target.Alias = alias
			if _, ok := result.Targets[alias]; ok {
				return nil, fmt.Errorf("failed to parse Chartfile: target alias %q is used more than once", alias)
			}

			result.Targets[alias] = target
		case MetadataCommand{}.Name():
			if target == nil {
				return nil, fmt.Errorf("failed to parse Chartfile: cannot use command %q outside of target", n.Value)
			}
			command := &MetadataCommand{}
			err := command.Parse(n)
			if err != nil {
				return nil, err
			}

			target.Commands = append(target.Commands, command)
		case AddCommand{}.Name():
			if target == nil {
				return nil, fmt.Errorf("failed to parse Chartfile: cannot use command %q outside of target", n.Value)
			}
			command := &AddCommand{}
			err := command.Parse(n)
			if err != nil {
				return nil, err
			}

			target.Commands = append(target.Commands, command)
		case RemoveCommand{}.Name():
			if target == nil {
				return nil, fmt.Errorf("failed to parse Chartfile: cannot use command %q outside of target", n.Value)
			}
			command := &RemoveCommand{}
			err := command.Parse(n)
			if err != nil {
				return nil, err
			}

			target.Commands = append(target.Commands, command)
		}
	}

	return result, nil
}