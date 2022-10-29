package action

import (
	"context"
	"os"

	"github.com/chartwave/chartwave/pkg/k8s"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

type Deploy struct {
	yamlpath string
}

func (i *Deploy) Run(ctx context.Context) error {
	f, err := os.Open(i.yamlpath)
	if err != nil {
		return err
	}
	defer f.Close()

	manifests, err := k8s.UnmarshalYAML(f)
	if err != nil {
		return err
	}

	for _, manifest := range manifests {
		log.WithField("manifest", manifest.Name).WithField("dependencies", manifest.Dependencies).Info("detected resource dependencies")
	}

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
