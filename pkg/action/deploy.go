package action

import (
	"context"

	"github.com/chartwave/chartwave/pkg/k8s"
	"github.com/urfave/cli/v2"
)

type Deploy struct {
	yamlpath string
}

func (i *Deploy) Run(ctx context.Context) error {
	return k8s.ApplyManifest(i.yamlpath)
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
