package grunter

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"

	"github.com/mitchellh/mapstructure"
	"gopkg.in/yaml.v2"

	"github.com/romainframe/grunter/pkg/grunter/block"
	"github.com/romainframe/grunter/pkg/grunter/system"
)

const (
	ObjectKindBlock  = "Block"
	ObjectKindSystem = "System"
)

// Object represents a Grunter object with an API version, kind, metadata, and spec.
type Object struct {
	ApiVersion string            `yaml:"apiVersion"`
	Kind       string            `yaml:"kind"`
	Metadata   map[string]string `yaml:"metadata"`
	Spec       interface{}       `yaml:"spec"`

	block  block.Block
	system system.System
}

func NewObjectFromFile(objectPath string) (Object, error) {
	// Read the entire file into memory.
	fileContents, err := os.ReadFile(objectPath)
	if err != nil {
		return Object{}, err // Return an empty Object and the error.
	}

	// Determine the file extension to decide on the unmarshalling method.
	ext := filepath.Ext(objectPath)
	var object Object

	switch ext {
	case ".json":
		if err := json.Unmarshal(fileContents, &object); err != nil {
			return Object{}, err // Return an error if the JSON is invalid.
		}
	case ".yaml", ".yml":
		if err := yaml.Unmarshal(fileContents, &object); err != nil {
			return Object{}, err // Return an error if the YAML is invalid.
		}
	default:
		return Object{}, errors.New("unsupported file type")
	}

	if object.ApiVersion == "" {
		object.ApiVersion = "v1" // Set the default API version.
	}

	if err := object.isValidKind(); err != nil {
		return Object{}, err // Return an error if the kind is invalid.
	}

	if object.Spec == nil {
		return Object{}, errors.New("spec is required") // Return an error if the spec is missing.
	}

	return object, nil // Return the fully initialized Object.
}

func (o Object) isValidKind() error {
	switch o.Kind {
	case ObjectKindBlock, ObjectKindSystem:
		return nil
	default:
		return errors.New("invalid kind")
	}
}

func (o Object) Build() (Object, error) {
	switch o.Kind {
	case ObjectKindBlock:
		return o.buildBlock()
	case ObjectKindSystem:
		return o.buildSystem()
	default:
		return Object{}, errors.New("invalid kind")
	}
}

func (o Object) buildBlock() (Object, error) {
	// Unmarshal the spec into a block.
	var block block.Block
	if err := mapstructure.Decode(o.Spec, &block); err != nil {
		return Object{}, err
	}

	// Perform any additional setup or validation.
	b, err := block.Build("")
	if err != nil {
		return Object{}, err
	}

	// Return the updated Object with the block spec.
	o.block = b
	return o, nil
}

func (o Object) buildSystem() (Object, error) {
	// Unmarshal the spec into a block.
	var sys system.System
	if err := mapstructure.Decode(o.Spec, &sys); err != nil {
		return Object{}, err
	}

	// Perform any additional setup or validation.
	b, err := sys.Build()
	if err != nil {
		return Object{}, err
	}

	// Return the updated Object with the block spec.
	o.system = b
	return o, nil
}
