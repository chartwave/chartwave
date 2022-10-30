package chart

import (
	"fmt"
	"io/fs"
	"os"
	"path"
	"strings"

	"github.com/chartwave/chartwave/pkg/helper"
	log "github.com/sirupsen/logrus"
	"helm.sh/helm/v3/pkg/action"
)

func (c *Chart) Build(dstDir string) error {
	err := os.RemoveAll(dstDir)
	if err != nil {
		return fmt.Errorf("failed to drop dst directory: %w", err)
	}

	if len(c.Targets) == 0 {
		return ErrNoTargets
	}

	tmpPaths := make(map[string]string)
	for name := range c.Targets {
		tmp, err := os.MkdirTemp("", "*")
		if err != nil {
			return err
		}
		// defer func(path string) {
		// 	err := os.RemoveAll(path)
		// 	if err != nil {
		// 		log.WithError(err).WithField("path", path).Error("failed to clean up temporary directory")
		// 	}
		// }(tmp)

		tmpPaths[name] = tmp
	}

	depGraph, err := c.buildDependenciesGraph()
	if err != nil {
		return err
	}

	targets := depGraph.Run()
	for target := range targets {
		name := target.Data.Alias
		err := target.Data.Build(tmpPaths[name], tmpPaths)
		if err != nil {
			target.SetFailed()
			return err
		}

		target.SetSucceeded()
	}

	tmpPath := tmpPaths[""]
	err = os.Rename(
		tmpPath,
		dstDir,
	)
	if err != nil {
		return fmt.Errorf("failed to move built chart: %w", err)
	}

	return nil
}

func (t *ChartTarget) Build(basePath string, allPaths map[string]string) error {
	l := log.WithField("target", t.Alias)
	l.WithField("path", basePath).Debug("starting building target")

	err := t.downloadBase(basePath)
	if err != nil {
		return err
	}

	for _, command := range t.Commands {
		err = command.Run(basePath, allPaths)
		if err != nil {
			return fmt.Errorf("failed to run %q command: %w", command.Name(), err)
		}
	}

	return nil
}

func (t *ChartTarget) downloadBase(basePath string) error {
	if t.Base == "scratch" {
		err := os.Mkdir(path.Join(basePath, "templates"), 0o755)
		if err != nil {
			return fmt.Errorf("failed to create templates dir in scratch: %w", err)
		}

		return nil
	}

	parts := strings.SplitN(t.Base, ":", 2)
	chart := parts[0]
	version := ""
	if len(parts) == 2 {
		version = parts[1]
	}
	log.WithField("chart", chart).WithField("version", version).Info("downloading base chart")

	// TODO: update repositories and set up registries
	cfg, err := helper.NewCfg()
	if err != nil {
		return err
	}

	pull := action.NewPullWithOpts(action.WithConfig(cfg))
	pull.Settings = helper.NewHelm()
	pull.Version = version

	pull.DestDir = basePath
	pull.Untar = true
	pull.UntarDir = basePath

	logs, err := pull.Run(chart)
	if logs != "" {
		log.StandardLogger().Print(logs)
	}

	if err != nil {
		return fmt.Errorf("failed to download and unarchive chart: %w", err)
	}

	// helm makes a subdirectory for unarchived chart
	// need to find it and move to the right place

	d := os.DirFS(basePath)
	entries, err := fs.ReadDir(d, ".")
	if err != nil {
		return fmt.Errorf("failed to process unarchived chart: %w", err)
	}
	if len(entries) != 1 {
		return fmt.Errorf("cannot find unarchived chart")
	}
	s := entries[0]

	entries, err = fs.ReadDir(d, s.Name())
	if err != nil {
		return fmt.Errorf("read unarchived chart: %w", err)
	}
	for _, entry := range entries {
		err = os.Rename(
			path.Join(basePath, s.Name(), entry.Name()),
			path.Join(basePath, entry.Name()),
		)
		if err != nil {
			return fmt.Errorf("failed to move files: %w", err)
		}
	}

	err = os.RemoveAll(path.Join(basePath, s.Name()))
	if err != nil {
		log.WithError(err).Warn("failed to remove old empty directory")
	}

	return nil
}
