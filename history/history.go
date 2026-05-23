package history

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
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

func PrintHistory(path string) error {
	records, err := ReadAll(path)
	if err != nil {
		return fmt.Errorf("history: could not read history: %w", err)
	}

	for _, record := range records {
		fmt.Printf("%s: %s -> %s (rule: %s, action: %s)\n",
			record.Timestamp.Format("2006-01-02 15:04:05"), record.Source, record.Destination, record.Rule, record.Action)
	}
	return nil
}

func ReadAll(path string) ([]Record, error) {
	file, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("history: could not open history file: %w", err)
	}
	defer file.Close()

	records := []Record{}
	decoder := json.NewDecoder(file)

	for {
		record := Record{}
		err := decoder.Decode(&record)
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}

			return nil, fmt.Errorf("history: could not decode record: %w", err)
		}
		records = append(records, record)

	}
	return records, nil
}

func Last(path string) (Record, bool, error) {
	records, err := ReadAll(path)
	if err != nil {
		return Record{}, false, fmt.Errorf("history: could not get record: %w", err)
	}
	if len(records) == 0 {
		return Record{}, false, nil
	}

	return records[len(records)-1], true, nil
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
