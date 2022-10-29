package k8s

import (
	"errors"
	"fmt"
	"io"
	"reflect"
	"strings"

	"github.com/chartwave/chartwave/pkg/yamlpath"
	"github.com/helmwave/helmwave/pkg/release/dependency"
	"gopkg.in/yaml.v3"
)

type Manifest struct {
	Name string
	Data map[string]interface{}

	dependencies []Dependency
}

type Dependency struct {
	Path *yamlpath.YamlPath
	Name string
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

		m.buildDependencies()
		result = append(result, m)
	}

	return result, nil
}

func (m *Manifest) UnmarshalYAML(value *yaml.Node) error {
	if len(value.Content) > 0 {
		m.Name = value.Content[0].HeadComment
		m.Name = strings.TrimPrefix(m.Name, "#")
		m.Name = strings.TrimSpace(m.Name)
	}

	return value.Decode(&m.Data)
}

func (m *Manifest) buildDependencies() {
	m.dependencies = m.getDependencies()
}

func (m *Manifest) getDependencies() []Dependency {
	return GetDependencies(&yamlpath.YamlPath{
		Root: &yamlpath.SubPath{
			Key: m.Name,
		},
	}, m.Data)
}

func GetDependencies(currentPath *yamlpath.YamlPath, v interface{}) []Dependency {
	value := reflect.ValueOf(v)
	switch value.Kind() {
	case reflect.Map:
		result := make([]Dependency, 0)

		iter := value.MapRange()
		for iter.Next() {
			i := interface{}(iter.Value().Interface())
			path := currentPath.Copy()
			path.AddSub(&yamlpath.SubPath{
				Key: iter.Key().String(),
			})
			result = append(result, GetDependencies(path, i)...)
		}

		return result

	case reflect.Slice:
		result := make([]Dependency, 0)

		for i := 0; i < value.Len(); i++ {
			k := interface{}(value.Index(i).Interface())
			path := currentPath.Copy()
			path.AddSub(&yamlpath.SubPath{
				Index: &i,
			})
			result = append(result, GetDependencies(path, k)...)
		}

		return result
	case reflect.String:
		s := v.(string)

		y, err := yamlpath.ParsePath(s)
		if err != nil {
			return []Dependency{}
		}

		return []Dependency{
			{Path: currentPath, Name: y.Name()},
		}
	default:
		return []Dependency{}
	}
}

func GetManifestsGraph() *dependency.Graph[string, Manifest] {
	return dependency.NewGraph[string, Manifest]()
}
