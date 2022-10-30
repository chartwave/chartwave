package command

import (
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/moby/buildkit/frontend/dockerfile/parser"
)

type RemoveCommand struct {
	Path string
}

func (RemoveCommand) Name() string {
	return "REMOVE"
}

func (command *RemoveCommand) Run(basePath string, _ map[string]string) error {
	err := os.RemoveAll(path.Join(basePath, command.Path))
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

	// TODO: validate to prohibit going outside tmpDir via `/blalba` or `../../../`
	return nil
}
