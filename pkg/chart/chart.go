package chart

import (
	"errors"

	"github.com/moby/buildkit/frontend/dockerfile/parser"
)

type Chart struct {
	Path    string
	Targets map[string]*ChartTarget
}

type ChartTarget struct {
	Base     string
	Alias    string
	Commands []ChartCommand
}

type ChartCommand interface {
	Name() string
	Run(string, map[string]string) error
	Parse(*parser.Node) error
}

var (
	ErrNoTargets = errors.New("chart has no targets")
)
