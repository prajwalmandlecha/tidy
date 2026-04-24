package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

type Config struct {
	WatchDirs []string `yaml:"watch_dirs"`
	Rules     []Rule   `yaml:"rules"`
}

type Rule struct {
	Name       string   `yaml:"name"`
	Extensions []string `yaml:"extensions,omitempty"`
	Dest       string   `yaml:"dest"`
	Pattern    string   `yaml:"pattern,omitempty"`
}

func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg Config
	err = yaml.Unmarshal(data, &cfg)
	if err != nil {
		return nil, err
	}

	err = cfg.expand()
	if err != nil {
		return nil, err
	}

	err = cfg.validate()
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}

func expandPath(path string) (string, error) {
	if !strings.HasPrefix(path, "~") {
		return path, nil
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, path[1:]), nil
}

func (cfg *Config) expand() error {
	for i, dir := range cfg.WatchDirs {
		expanded, err := expandPath(dir)
		if err != nil {
			return err
		}
		cfg.WatchDirs[i] = expanded
	}

	for i, rule := range cfg.Rules {
		expanded, err := expandPath(rule.Dest)
		if err != nil {
			return err
		} 
		cfg.Rules[i].Dest = expanded
	}

	return nil
}

func (c *Config) validate() error {
	if len(c.WatchDirs) == 0 {
		return errors.New("config: no watch_dirs specified")
	}
	for _, rule := range c.Rules {
		if rule.Dest == "" {
			return fmt.Errorf("config: rule %q has no dest", rule.Name)
		}
		if len(rule.Extensions) == 0 && rule.Pattern == "" {
			return fmt.Errorf("config: rule %q has no extensions or pattern", rule.Name)
		}
	}
	return nil
}
