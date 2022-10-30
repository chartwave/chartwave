package builder

import (
	"fmt"
	"strings"

	"github.com/moby/buildkit/frontend/dockerfile/parser"
)

type ChartCommand interface {
	Name() string
	Run(string) error
	Parse(*parser.Node) error
}

type MetadataCommand struct {
	Key   string
	Value string
}

func (MetadataCommand) Name() string {
	return "METADATA"
}

func (MetadataCommand) Run(string) error {
	return nil
}

func (command *MetadataCommand) Parse(n *parser.Node) error {
	s := strings.TrimPrefix(n.Original, command.Name())
	s = strings.TrimSpace(s)
	parts := strings.SplitN(s, "=", 2)
	if len(parts) != 2 {
		return fmt.Errorf("failed to parse %q command: invalid syntax", command.Name())
	}
	command.Key = parts[0]
	command.Value = parts[1]

	return nil
}

type AddCommand struct {
	Src, Dst string
}

func (AddCommand) Name() string {
	return "ADD"
}

func (AddCommand) Run(string) error {
	return nil
}

func (command *AddCommand) Parse(n *parser.Node) error {
	if n.Next == nil || n.Next.Next == nil {
		return fmt.Errorf("failed to parse %q command: two arguments required - src and dst", command.Name())
	}

	command.Src = n.Next.Value
	command.Dst = n.Next.Next.Value

	return nil
}

type RemoveCommand struct {
	Path string
}

func (RemoveCommand) Name() string {
	return "REMOVE"
}

func (RemoveCommand) Run(string) error {
	return nil
}

func (command *RemoveCommand) Parse(n *parser.Node) error {
	s := strings.TrimPrefix(n.Original, command.Name())
	s = strings.TrimSpace(s)
	if s == "" {
		return fmt.Errorf("failed to parse %q command: path required", command.Name())
	}
	command.Path = s

	return nil
}
