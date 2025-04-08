package config

import (
	"os"
	"sort"

	"github.com/goccy/go-yaml"
)

type TokenConfig struct {
	Name     string `yaml:"name"`
	Pattern  string `yaml:"pattern"`
	Priority int    `yaml:"priority"`
}

type Config struct {
	Tokens []TokenConfig `yaml:"tokens"`
}

// Returns tokens' configuraion in sorted by priority order
func ParseConfig(configPath string) ([]TokenConfig, error) {
	configBytes, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	var config Config
	if err := yaml.Unmarshal(configBytes, &config); err != nil {
		return nil, err
	}

	sort.Slice(config.Tokens, func(i, j int) bool { return config.Tokens[i].Priority < config.Tokens[j].Priority })
	return config.Tokens, nil
}
