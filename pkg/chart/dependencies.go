package chart

import (
	"fmt"

	"github.com/chartwave/chartwave/pkg/chart/command"
	"github.com/helmwave/helmwave/pkg/release/dependency"
	log "github.com/sirupsen/logrus"
)

func (c *Chart) buildDependenciesGraph() (*dependency.Graph[string, *ChartTarget], error) {
	graph := dependency.NewGraph[string, *ChartTarget]()

	for name, target := range c.Targets {
		err := graph.NewNode(name, target)
		if err != nil {
			return nil, fmt.Errorf("failed to add node %q to dependency graph: %w", name, err)
		}

		for _, dep := range target.getDependencies() {
			log.WithField("dependant", name).WithField("dependency", dep).Debug("adding dependency between targets")
			graph.AddDependency(name, dep)
		}
	}

	err := graph.Build()
	if err != nil {
		return nil, fmt.Errorf("failed to build dependency graph: %w", err)
	}

	return graph, nil
}

func (t *ChartTarget) getDependencies() []string {
	deps := make(map[string]bool)
	for _, cmd := range t.Commands {
		switch cmd := cmd.(type) {
		case *command.AddCommand:
			if cmd.From != "" {
				deps[cmd.From] = true
			}
		}
	}

	result := make([]string, 0, len(deps))
	for name := range deps {
		result = append(result, name)
	}

	return result
}
