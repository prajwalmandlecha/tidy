package engine

import (
	"crypto/sha256"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/prajwalmandlecha/tidy/config"
)

type Action string

const (
	ActionMoved           Action = "moved"
	ActionRenamedConflict Action = "renamed_conflict"
)

type MoveResult struct {
	Source      string
	Destination string
	Rule        string
	Action      Action
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

func Apply(filePath string, rule config.Rule, dryRun bool) (*MoveResult, error) {
	fileName := filepath.Base(filePath)
	destPath := filepath.Join(rule.Dest, fileName)
	if filepath.Clean(filePath) == filepath.Clean(destPath) {
		return nil, nil
	}

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return nil, nil
	}

	if filepath.Clean(filepath.Dir(filePath)) == filepath.Clean(rule.Dest) {
		return nil, nil
	}

	action := ActionMoved

	if _, err := os.Stat(destPath); err == nil {
		duplicate, err := isDuplicate(filePath, destPath)
		if err != nil {
			return nil, fmt.Errorf("engine: error comparing files %q and %q: %w", filePath, destPath, err)
		}
		if duplicate {
			fmt.Printf("skipping duplicate: %s\n", filepath.Base(filePath))
			return nil, nil
		}
		destPath = resolveDestPath(rule.Dest, fileName)
		action = ActionRenamedConflict
	}

	if dryRun {
		fmt.Printf("[dry-run] %s  →  %s\n", filePath, destPath)
		return nil, nil
	}

	if err := os.MkdirAll(rule.Dest, 0755); err != nil {
		return nil, fmt.Errorf("engine: could not create destination dir %q: %w", rule.Dest, err)
	}
	if err := os.Rename(filePath, destPath); err != nil {
		return nil, fmt.Errorf("engine: could not move %q to %q: %w", filePath, destPath, err)
	}

	fmt.Printf("moved %s  →  %s\n", filePath, destPath)
	return &MoveResult{
		Source:      filePath,
		Destination: destPath,
		Rule:        rule.Name,
		Action:      action,
	}, nil
}

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

func IsIgnored(filePath string, ignorePatterns []string) bool {
	fileName := filepath.Base(filePath)
	for _, pattern := range ignorePatterns {
		pattern = strings.TrimSpace(pattern)
		if pattern == "" {
			continue
		}
		matched, err := filepath.Match(pattern, fileName)
		if err == nil && matched {
			return true
		}
	}
	return false
}

func ProcessFile(filePath string, ignorePatterns []string, rules []config.Rule, dryRun bool) (*MoveResult, error) {
	if IsIgnored(filePath, ignorePatterns) {
		return nil, nil
	}

	for _, rule := range rules {
		matches := Matches(filePath, rule)
		if matches {
			return Apply(filePath, rule, dryRun)
		}
	}
	return nil, nil
}

func hashFile(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()

	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", h.Sum(nil)), nil
}

func isDuplicate(src, dst string) (bool, error) {
	srcInfo, err := os.Stat(src)
	if err != nil {
		return false, err
	}
	dstInfo, err := os.Stat(dst)
	if err != nil {
		return false, err
	}
	if srcInfo.Size() != dstInfo.Size() {
		return false, nil
	}
	srcHash, err := hashFile(src)
	if err != nil {
		return false, err
	}

	dstHash, err := hashFile(dst)
	if err != nil {
		return false, err
	}

	return srcHash == dstHash, nil
}

func resolveDestPath(dir, fileName string) string {
	destPath := filepath.Join(dir, fileName)
	if _, err := os.Stat(destPath); os.IsNotExist(err) {
		return destPath
	}

	ext := filepath.Ext(fileName)
	name := strings.TrimSuffix(fileName, ext)

	for i := 1; ; i++ {
		newPath := filepath.Join(dir, fmt.Sprintf("%s_%d%s", name, i, ext))
		if _, err := os.Stat(newPath); os.IsNotExist(err) {
			return newPath
		}
	}
}
