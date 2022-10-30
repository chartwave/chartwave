package command

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path"
	"strings"

	"github.com/moby/buildkit/frontend/dockerfile/parser"
	"github.com/otiai10/copy"
)

type AddCommand struct {
	From     string
	Src, Dst string
}

func (AddCommand) Name() string {
	return "ADD"
}

func (command *AddCommand) Run(basePath string, allPaths map[string]string) error {
	baseDir := "."
	if command.From != "" {
		baseDir = allPaths[command.From]
	}

	dst := path.Join(basePath, command.Dst)
	if strings.HasSuffix(command.Dst, "/") {
		dst += "/"
	}

	matches, err := fs.Glob(os.DirFS(baseDir), command.Src)
	if err != nil {
		return fmt.Errorf("failed to parse glob: %w", err)
	}

	if len(matches) == 0 {
		return fmt.Errorf("cannot find such file %q", command.Src)
	}
	for _, match := range matches {
		err = command.copy(path.Join(baseDir, match), dst)
		if err != nil {
			return fmt.Errorf("failed to copy %q: %w", match, err)
		}
	}

	return nil
}

func (command *AddCommand) copy(src, dst string) error {
	info, err := os.Stat(src)
	if err != nil {
		return fmt.Errorf("failed to stat src %q: %w", src, err)
	}
	if info.IsDir() {
		entries, err := os.ReadDir(src)
		if err != nil {
			return fmt.Errorf("cannot read directory %q: %w", src, err)
		}

		for _, entry := range entries {
			err = command.copy(
				path.Join(src, entry.Name()),
				path.Join(dst, entry.Name()),
			)
			if err != nil {
				return nil
			}
		}
	} else {
		info, err := os.Stat(dst)
		if err != nil && !errors.Is(err, os.ErrNotExist) {
			return fmt.Errorf("failed to stat dst %q: %w", dst, err)
		}

		if err == nil {
			// copy.Copy cannot copy file into directory, full path needs to be specified
			if info.IsDir() {
				dst = path.Join(dst, path.Base(src))
			}
		}

		err = copy.Copy(src, dst)
		if err != nil {
			return fmt.Errorf("failed to copy file %q: %w", src, err)
		}
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
