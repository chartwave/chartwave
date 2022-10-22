package version

import (
	log "github.com/sirupsen/logrus"
)

// Version is chartwave binary version.
// It should be var not const.
// It will override by goreleaser during release.
// -X github.com/chartwave/chartwave/pkg/version.Version={{ .Version }}.
//
//nolint:gochecknoglobals
var Version = "dev"

// Check compares chartwave versions and logs difference.
func Check(a, b string) {
	if a != b {
		log.Warn("⚠️ Unsupported version ", b)
		log.Debug("🌊 chartwave version ", a)
	}
}
