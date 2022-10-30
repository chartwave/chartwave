package command

import (
	"fmt"
	"strings"

	"github.com/moby/buildkit/frontend/dockerfile/parser"
)

type MetadataCommand struct {
	Key   string
	Value string
}

func (MetadataCommand) Name() string {
	return "METADATA"
}

func (MetadataCommand) Run(string, map[string]string) error {
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
