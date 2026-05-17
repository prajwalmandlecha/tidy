package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func checkExistingDir(path string) error {
	info, err := os.Stat(path)
	if err != nil {
		return err
	}
	if !info.IsDir() {
		return fmt.Errorf("path is not a directory: %s", path)
	}
	return nil
}

func (c *Config) Validate() error {
	var errs []error

	if len(c.WatchDirs) == 0 {
		errs = append(errs, errors.New("config: no watch_dirs specified"))
	}

	for _, dir := range c.WatchDirs {
		err := checkExistingDir(dir)
		if err != nil {
			if os.IsNotExist(err) {
				errs = append(errs, fmt.Errorf("config: watch dir does not exist: %s", dir))
			} else {
				errs = append(errs, fmt.Errorf("config: invalid watch path: %w", err))
			}
			continue
		}
	}

	seenRuleNames := make(map[string]bool)

	for _, rule := range c.Rules {
		dest := strings.TrimSpace(rule.Dest)
		if dest == "" {
			errs = append(errs, fmt.Errorf("config: rule %q has no dest", rule.Name))
		}
		if len(rule.Extensions) == 0 && rule.Pattern == "" {
			errs = append(errs, fmt.Errorf("config: rule %q has no extensions or pattern", rule.Name))
		}

		if rule.Pattern != "" {
			_, err := filepath.Match(rule.Pattern, "test")
			if err != nil {
				errs = append(errs, fmt.Errorf("config: rule %q has invalid pattern %q: %w", rule.Name, rule.Pattern, err))
			}
		}

		for i, ext := range rule.Extensions {
			trimmed := strings.TrimSpace(ext)
			if trimmed == "" {
				errs = append(errs, fmt.Errorf("config: rule %q has empty extension at index %d", rule.Name, i))
				continue
			}

			if trimmed == "." || trimmed == ".." || strings.ContainsAny(trimmed, `/\`) {
				errs = append(errs, fmt.Errorf("config: rule %q has invalid extension %q", rule.Name, trimmed))
			}
		}

		trimmed := strings.TrimSpace(rule.Name)
		if trimmed != "" {
			if seenRuleNames[trimmed] {
				errs = append(errs, fmt.Errorf("config: duplicate rule name %q", trimmed))
			}
			seenRuleNames[trimmed] = true
		}

		if dest == "" {
			continue
		}

		if err := checkExistingDir(rule.Dest); err != nil {
			if os.IsNotExist(err) {
				continue
			}
			errs = append(errs, fmt.Errorf("config: rule %q has invalid dest path: %w", rule.Name, err))

		}
	}
	return errors.Join(errs...)
}
