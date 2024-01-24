package main

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Title string `yaml:"title"`
	Usage string `yaml:"usage"`

	Commands CliCommandEntries `yaml:"commands"`
}

func loadConfig() (*Config, error) {
	path := filepath.Join(filepath.Dir(os.Args[0]), ".console.yml")

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var c Config
	if err := yaml.Unmarshal(data, &c); err != nil {
		return nil, err
	}

	return &c, nil
}
