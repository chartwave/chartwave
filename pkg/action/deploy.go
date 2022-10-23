package action

import (
	"context"

	"github.com/chartwave/chartwave/pkg/yamlpath"
	"github.com/davecgh/go-spew/spew"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

type Deploy struct {
	yamlpath string
}

func (i *Deploy) Run(ctx context.Context) error {
	res, err := yamlpath.ParsePath(i.yamlpath)
	if err != nil {
		log.WithError(err).Error("failed to parse yamlpath")
		return err
	}
	spew.Dump(res)
	return nil
}

func (i *Deploy) Cmd() *cli.Command {
	return &cli.Command{
		Name:   "deploy",
		Usage:  "Deploy changes",
		Flags:  i.flags(),
		Action: toCtx(i.Run),
	}
}

func (i *Deploy) flags() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:        "path",
			Destination: &i.yamlpath,
			Required:    true,
		},
	}
}
