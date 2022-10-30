package chart

import (
	"fmt"
	"os"
	"path"
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
	From     string
	Src, Dst string
}

func (AddCommand) Name() string {
	return "ADD"
}

func (command *AddCommand) Run(string) error {
	if command.From != "" {
		return fmt.Errorf("--from is not supported yet")
	}
	return nil
}

func (command *AddCommand) Parse(n *parser.Node) error {
	if n.Next == nil || n.Next.Next == nil {
		return fmt.Errorf("failed to parse %q command: two arguments required - src and dst", command.Name())
	}
	for _, flag := range n.Flags {
		if strings.HasPrefix(flag, "--from=") {
			command.From = strings.TrimPrefix(flag, "--from=")
		}
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

func (command *RemoveCommand) Run(basePath string) error {
	err := os.Remove(path.Join(basePath, command.Path))
	if err != nil {
		return fmt.Errorf("failed to remove file: %w", err)
	}

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
