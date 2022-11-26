package k8s

import (
	"errors"
	"fmt"
	"gopkg.in/yaml.v3"
	"io"
)

type Manifest map[string]interface{}

var (
	UnnamedErr     = errors.New(fmt.Sprintf("manifest has no name (%q annotation)", AnnotationUniqueName))
	InvalidNameErr = errors.New("manifest name is invalid")
)

func (m Manifest) Name() (string, error) {
	meta, ok := m["metadata"]
	if !ok {
		return "", UnnamedErr
	}

	metadata, ok := meta.(map[string]interface{})
	if !ok {
		return "", UnnamedErr
	}

	name, ok := metadata[AnnotationUniqueName]
	if !ok {
		return "", UnnamedErr
	}

	nameString, ok := name.(string)
	if !ok {
		return "", InvalidNameErr
	}

	return nameString, nil
}

// UnmarshalYAML unmarshals multiple manifests from input stream and builds dependencies
func UnmarshalYAML(input io.Reader) ([]*Manifest, error) {
	result := make([]*Manifest, 0)
	decoder := yaml.NewDecoder(input)

	for {
		m := &Manifest{}
		err := decoder.Decode(m)
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return result, fmt.Errorf("failed to unmarshal manifest: %w", err)
		}

		result = append(result, m)
	}

	return result, nil
}
