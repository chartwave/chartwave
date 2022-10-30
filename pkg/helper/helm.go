package helper

import (
	"fmt"

	log "github.com/sirupsen/logrus"
	"helm.sh/helm/v3/pkg/action"
	helm "helm.sh/helm/v3/pkg/cli"
	"helm.sh/helm/v3/pkg/registry"
	"k8s.io/cli-runtime/pkg/genericclioptions"
)

// NewHelm is a hack to create an instance of helm CLI and specifying namespace without environment variables.
func NewHelm() *helm.EnvSettings {
	return helm.New()
}

func NewCfg() (*action.Configuration, error) {
	cfg := new(action.Configuration)
	config := genericclioptions.NewConfigFlags(false)

	err := cfg.Init(config, "", "", log.Debugf)
	if err != nil {
		return nil, fmt.Errorf("failed to create helm configuration: %w", err)
	}

	helmRegistryClient, err := registry.NewClient(
		registry.ClientOptWriter(log.StandardLogger().Writer()),
		// registry.ClientOptCredentialsFile(Helm.RegistryConfig),
	)
	cfg.RegistryClient = helmRegistryClient

	return cfg, nil
}
