package config

import (
	"io/ioutil"

	"gopkg.in/yaml.v3"

	"github.com/kalogs-c/resizenator/image"
)

type Config struct {
	Prefix            string            `yaml:"prefix"`
	Sizes             []image.ImageSize `yaml:"sizes"`
	Algorithm         image.Algorithm   `yaml:"algorithm"`
	MaxConcurrency    int               `yaml:"max_concurrency"`
	DeleteAfterResize bool              `yaml:"delete_after_resize"`
	TargetFormat      image.ImageFormat `yaml:"target_format"`
}

func Load(path string) (*Config, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	return &config, nil
}
