package chart

import (
	"errors"
)

type Chart struct {
	Targets map[string]*ChartTarget
}

type ChartTarget struct {
	Base     string
	Alias    string
	Commands []ChartCommand
}

var (
	ErrNoTargets = errors.New("chart has no targets")
)
