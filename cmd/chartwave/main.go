package main

import (
	"fmt"
	"os"

	"github.com/chartwave/chartwave/pkg/action"
	logSetup "github.com/chartwave/chartwave/pkg/log"
	chartwave "github.com/chartwave/chartwave/pkg/version"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

//nolint:gochecknoglobals // we need global list of commands
var commands = []*cli.Command{
	new(action.Build).Cmd(),
	new(action.Deploy).Cmd(),
	version(),
	completion(),
}

func main() {
	c := CreateApp()

	defer recoverPanic()

	if err := c.Run(os.Args); err != nil {
		log.Fatal(err) //nolint:gocritic // we try to recover panics, not regural command errors
	}
}

func recoverPanic() {
	if r := recover(); r != nil {
		switch r.(type) {
		case CommandNotFoundError:
			log.Error(r)
			log.Exit(127)
		default:
			log.Panic(r)
		}
	}
}

// CreateApp creates *cli.App with all commands.
func CreateApp() *cli.App {
	c := cli.NewApp()

	c.EnableBashCompletion = true
	// c.Usage = "is like docker-compose for helm"
	c.Version = chartwave.Version
	// c.Description = "This tool helps you compose your helm releases!\n" +
	// 	"0. $ chartwave yml\n" +
	// 	"1. $ chartwave build\n" +
	// 	"2. $ chartwave up\n"

	logSet := logSetup.Settings{}
	c.Before = logSet.Run
	c.Flags = logSet.Flags()

	c.Commands = commands
	c.CommandNotFound = command404

	return c
}

// CommandNotFoundError is return when CLI command is not found.
type CommandNotFoundError struct {
	Command string
}

func (e CommandNotFoundError) Error() string {
	return fmt.Sprintf("👻 Command %q not found", e.Command)
}

func command404(c *cli.Context, s string) {
	err := CommandNotFoundError{
		Command: s,
	}
	panic(err)
}

func version() *cli.Command {
	return &cli.Command{
		Name:    "version",
		Aliases: []string{"ver"},
		Usage:   "Show shorts version",
		Action: func(c *cli.Context) error {
			fmt.Println(chartwave.Version) //nolint:forbidigo // we need to use fmt.Println here

			return nil
		},
	}
}
