package engine

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/prajwalmandlecha/tidy/config"
)

func matchesExtension(filePath string, extensions []string) bool {
	fileName := strings.ToLower(filepath.Base(filePath))

	for _, extension := range extensions {
		normalized := strings.ToLower(strings.TrimSpace(extension))
		if normalized == "" {
			continue
		}
		if !strings.HasPrefix(normalized, ".") {
			normalized = "." + normalized
		}
		if strings.HasSuffix(fileName, normalized) {
			return true
		}
	}

	return false
}

func Matches(filePath string, rule config.Rule) bool {
	if matchesExtension(filePath, rule.Extensions) {
		return true
	}

	if rule.Pattern == "" {
		return false
	}

	matched, err := filepath.Match(rule.Pattern, filepath.Base(filePath))
	if err != nil {
		return false
	}
	return matched
}

func Apply(filePath string, rule config.Rule, dryRun bool) error {
	fileName := filepath.Base(filePath)
	destPath := filepath.Join(rule.Dest, fileName)
	if filepath.Clean(filePath) == filepath.Clean(destPath) {
		return nil
	}

	if filepath.Dir(filePath) == rule.Dest {
		return nil
	}

	if dryRun {
		fmt.Printf("[dry-run] %s  →  %s\n", filePath, destPath)
		return nil
	}

	if err := os.MkdirAll(rule.Dest, 0755); err != nil {
		return fmt.Errorf("engine: could not create destination dir %q: %w", rule.Dest, err)
	}
	if err := os.Rename(filePath, destPath); err != nil {
		return fmt.Errorf("engine: could not move %q to %q: %w", filePath, destPath, err)
	}

	fmt.Printf("moved %s  →  %s\n", filePath, destPath)
	return nil
}

func ProcessFile(filePath string, rules []config.Rule, dryRun bool) error {
	for _, rule := range rules {
		matches := Matches(filePath, rule)
		if matches {
			return Apply(filePath, rule, dryRun)
		}
	}
	return nil
}
