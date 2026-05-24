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

type Entry struct {
	ID     int
	Record Record
}

type Record struct {
	Source      string    `json:"source"`
	Destination string    `json:"destination"`
	Rule        string    `json:"rule"`
	Action      string    `json:"action"`
	Timestamp   time.Time `json:"timestamp"`
}

func ReadEntries(path string) ([]Entry, error) {
	file, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("history: could not open history file: %w", err)
	}
	defer file.Close()

	entries := []Entry{}
	decoder := json.NewDecoder(file)
	id := 1
	for {
		record := Record{}
		err := decoder.Decode(&record)
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}

			return nil, fmt.Errorf("history: could not decode record: %w", err)
		}
		entries = append(entries, Entry{ID: id, Record: record})
		id++
	}
	return entries, nil
}

func Last(path string) (Record, bool, error) {
	entries, err := ReadEntries(path)
	if err != nil {
		return Record{}, false, fmt.Errorf("history: could not get record: %w", err)
	}
	if len(entries) == 0 {
		return Record{}, false, nil
	}

	return entries[len(entries)-1].Record, true, nil
}

func Latest(path string, limit int) ([]Entry, error) {
	entries, err := ReadEntries(path)
	if err != nil {
		return nil, fmt.Errorf("history: could not get records: %w", err)
	}
	if len(entries) == 0 {
		return nil, nil
	}
	if limit <= 0 {
		return nil, nil
	}

	if limit > len(entries) {
		limit = len(entries)
	}

	return entries[len(entries)-limit:], nil
}

func Find(path string, id int) (Entry, bool, error) {
	entries, err := ReadEntries(path)
	if err != nil {
		return Entry{}, false, fmt.Errorf("history: could not get records: %w", err)
	}
	for _, entry := range entries {
		if entry.ID == id {
			return entry, true, nil
		}
	}
	return Entry{}, false, nil
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
