package history

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

type Record struct {
	Source      string    `json:"source"`
	Destination string    `json:"destination"`
	Rule        string    `json:"rule"`
	Action      string    `json:"action"`
	Timestamp   time.Time `json:"timestamp"`
}

func Append(path string, record Record) error {
	err := os.MkdirAll(filepath.Dir(path), 0755)
	if err != nil {
		return fmt.Errorf("history: could not create history dir: %w", err)
	}
	file, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("history: could not open history file: %w", err)
	}
	defer file.Close()

	err = json.NewEncoder(file).Encode(record)
	if err != nil {
		return fmt.Errorf("history: could not write record: %w", err)
	}

	return nil
}

func DefaultPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("history: could not get user home dir: %w", err)
	}

	return filepath.Join(home, ".local/state/tidy/history.jsonl"), nil
}
