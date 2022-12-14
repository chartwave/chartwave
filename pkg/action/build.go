package action

import (
	"context"
	"os"

	"github.com/chartwave/chartwave/pkg/chart"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

type Build struct {
	path string
}

func (i *Build) Run(ctx context.Context) error {
	f, err := os.Open(i.path)
	if err != nil {
		return err
	}
	defer f.Close()

	chart, err := chart.ParseChartfile(f)
	if err != nil {
		log.WithError(err).Error("failed to parse chartfile")
		return err
	}

	err = chart.Build("dist")
	if err != nil {
		log.WithError(err).Error("failed to build chart")
		return err
	}

	return nil
}

func (i *Build) Cmd() *cli.Command {
	return &cli.Command{
		Name:   "build",
		Usage:  "Build chart",
		Flags:  i.flags(),
		Action: toCtx(i.Run),
	}
}

func (i *Build) flags() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:        "path",
			Destination: &i.path,
			Required:    true,
		},
	}
}
