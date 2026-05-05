package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Config represents the top-level patchwork configuration.
type Config struct {
	Version string `yaml:"version"`
	Repos   []Repo `yaml:"repos"`
}

// Repo represents a single tracked git repository.
type Repo struct {
	Name   string `yaml:"name"`
	Path   string `yaml:"path"`
	Remote string `yaml:"remote,omitempty"`
}

// Load reads and parses a patchwork config file from the given path.
func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading config %q: %w", path, err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parsing config %q: %w", path, err)
	}

	if err := cfg.validate(); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}

	return &cfg, nil
}

// Save writes the config to the given path in YAML format.
func (c *Config) Save(path string) error {
	data, err := yaml.Marshal(c)
	if err != nil {
		return fmt.Errorf("marshalling config: %w", err)
	}
	return os.WriteFile(path, data, 0o644)
}

func (c *Config) validate() error {
	seen := make(map[string]struct{}, len(c.Repos))
	for i, r := range c.Repos {
		if r.Name == "" {
			return fmt.Errorf("repo[%d]: name is required", i)
		}
		if r.Path == "" {
			return fmt.Errorf("repo %q: path is required", r.Name)
		}
		if _, dup := seen[r.Name]; dup {
			return fmt.Errorf("duplicate repo name %q", r.Name)
		}
		seen[r.Name] = struct{}{}
	}
	return nil
}
