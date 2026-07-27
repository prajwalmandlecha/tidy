package history

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"time"

	_ "modernc.org/sqlite"
)

type DB struct {
	sql *sql.DB
}

type Move struct {
	ID          int64
	Source      string
	Destination string
	Rule        string
	Action      string
	MovedAt     time.Time
	UndoneAt    *time.Time
}

func (db *DB) Append(move Move) (Move, error) {
	result, err := db.sql.Exec(
		`INSERT INTO moves (source, destination, rule, action, timestamp)
         VALUES (?, ?, ?, ?, ?)`,
		move.Source,
		move.Destination,
		move.Rule,
		move.Action,
		move.MovedAt.UTC().Format(time.RFC3339Nano),
	)
	if err != nil {
		return Move{}, fmt.Errorf("history: append move: %w", err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		return Move{}, fmt.Errorf("history: get inserted move ID: %w", err)
	}

	move.ID = id
	return move, nil
}

func scanMove(scan func(dest ...any) error) (Move, error) {
	var m Move
	var movedAtStr string
	var undoneAtStr sql.NullString

	err := scan(&m.ID, &m.Source, &m.Destination, &m.Rule, &m.Action, &movedAtStr, &undoneAtStr)
	if err != nil {
		return Move{}, err
	}

	t, err := time.Parse(time.RFC3339Nano, movedAtStr)
	if err != nil {
		t, _ = time.Parse(time.RFC3339, movedAtStr)
	}
	m.MovedAt = t

	if undoneAtStr.Valid && undoneAtStr.String != "" {
		ut, err := time.Parse(time.RFC3339Nano, undoneAtStr.String)
		if err != nil {
			ut, _ = time.Parse(time.RFC3339, undoneAtStr.String)
		}
		m.UndoneAt = &ut
	}

	return m, nil
}

func (db *DB) Latest(limit int) ([]Move, error) {
	q := `SELECT id, source, destination, rule, action, timestamp, undone_at
          FROM moves ORDER BY timestamp DESC`
	if limit > 0 {
		q += fmt.Sprintf(` LIMIT %d`, limit)
	}
	rows, err := db.sql.Query(q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var moves []Move
	for rows.Next() {
		m, err := scanMove(rows.Scan)
		if err != nil {
			return nil, err
		}
		moves = append(moves, m)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return moves, nil
}

func (db *DB) FindPending(id int64) (Move, bool, error) {
	q := `SELECT id, source, destination, rule, action, timestamp, undone_at FROM moves WHERE id = ? AND undone_at IS NULL`
	m, err := scanMove(db.sql.QueryRow(q, id).Scan)
	if err != nil {
		if err == sql.ErrNoRows {
			return Move{}, false, nil
		}
		return Move{}, false, err
	}

	return m, true, nil
}

func (db *DB) MarkUndone(id int64, at time.Time) error {
	q := `UPDATE moves SET undone_at = ? WHERE id = ?`
	_, err := db.sql.Exec(q, at.UTC().Format(time.RFC3339Nano), id)
	if err != nil {
		return err
	}
	return nil
}

func (db *DB) Pending(limit int) ([]Move, error) {
	q := `SELECT id, source, destination, rule, action, timestamp, undone_at FROM moves WHERE undone_at IS NULL ORDER BY timestamp DESC`
	if limit > 0 {
		q += fmt.Sprintf(` LIMIT %d`, limit)
	}
	rows, err := db.sql.Query(q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var moves []Move
	for rows.Next() {
		m, err := scanMove(rows.Scan)
		if err != nil {
			return nil, err
		}
		moves = append(moves, m)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return moves, nil
}

func Open(path string) (*DB, error) {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return nil, err
	}
	dbString := fmt.Sprintf("file:%s?_pragma=busy_timeout(5000)&_pragma=journal_mode(WAL)", path)
	s, err := sql.Open("sqlite", dbString)
	if err != nil {
		return nil, err
	}
	err = s.Ping()
	if err != nil {
		return nil, err
	}
	if _, err = s.Exec(`
		CREATE TABLE IF NOT EXISTS moves (
			id          INTEGER PRIMARY KEY AUTOINCREMENT,
			source      TEXT NOT NULL,
			destination TEXT NOT NULL,
			rule        TEXT NOT NULL,
			action      TEXT NOT NULL,
			timestamp   TEXT NOT NULL,
			undone_at   TEXT
		);
	`); err != nil {
		s.Close()
		return nil, err
	}
	return &DB{sql: s}, nil
}

func NewDB() (*DB, error) {
	dbPath, err := DefaultDBPath()
	if err != nil {
		return nil, err
	}
	return Open(dbPath)
}

func (db *DB) Close() error {
	return db.sql.Close()
}

func DefaultDBPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("history: could not get user home dir: %w", err)
	}
	return filepath.Join(home, ".local", "state", "tidy", "history.db"), nil
}

// old shizz
// type Record struct {
// 	Source      string
// 	Destination string
// 	Rule        string
// 	Action      string
// 	Timestamp   time.Time
// }

// type Entry struct {
// 	ID     int
// 	Record Record
// }

// func ReadEntries(path string) ([]Entry, error) {
// 	file, err := os.Open(path)
// 	if err != nil {
// 		if os.IsNotExist(err) {
// 			return nil, nil
// 		}
// 		return nil, fmt.Errorf("history: could not open history file: %w", err)
// 	}
// 	defer file.Close()

// 	entries := []Entry{}
// 	decoder := json.NewDecoder(file)
// 	id := 1
// 	for {
// 		record := Record{}
// 		err := decoder.Decode(&record)
// 		if err != nil {
// 			if errors.Is(err, io.EOF) {
// 				break
// 			}

// 			return nil, fmt.Errorf("history: could not decode record: %w", err)
// 		}
// 		entries = append(entries, Entry{ID: id, Record: record})
// 		id++
// 	}
// 	return entries, nil
// }

// func Last(path string) (Record, bool, error) {
// 	entries, err := ReadEntries(path)
// 	if err != nil {
// 		return Record{}, false, fmt.Errorf("history: could not get record: %w", err)
// 	}
// 	if len(entries) == 0 {
// 		return Record{}, false, nil
// 	}

// 	return entries[len(entries)-1].Record, true, nil
// }

// func Latest(path string, limit int) ([]Entry, error) {
// 	entries, err := ReadEntries(path)
// 	if err != nil {
// 		return nil, fmt.Errorf("history: could not get records: %w", err)
// 	}
// 	if len(entries) == 0 {
// 		return nil, nil
// 	}
// 	if limit <= 0 {
// 		return nil, nil
// 	}

// 	if limit > len(entries) {
// 		limit = len(entries)
// 	}

// 	return entries[len(entries)-limit:], nil
// }

// func Find(path string, id int) (Entry, bool, error) {
// 	entries, err := ReadEntries(path)
// 	if err != nil {
// 		return Entry{}, false, fmt.Errorf("history: could not get records: %w", err)
// 	}
// 	for _, entry := range entries {
// 		if entry.ID == id {
// 			return entry, true, nil
// 		}
// 	}
// 	return Entry{}, false, nil
// }

// func Append(path string, record Record) error {
// 	err := os.MkdirAll(filepath.Dir(path), 0755)
// 	if err != nil {
// 		return fmt.Errorf("history: could not create history dir: %w", err)
// 	}
// 	file, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
// 	if err != nil {
// 		return fmt.Errorf("history: could not open history file: %w", err)
// 	}
// 	defer file.Close()

// 	err = json.NewEncoder(file).Encode(record)
// 	if err != nil {
// 		return fmt.Errorf("history: could not write record: %w", err)
// 	}

// 	return nil
// }

// func DefaultPath() (string, error) {
// 	home, err := os.UserHomeDir()
// 	if err != nil {
// 		return "", fmt.Errorf("history: could not get user home dir: %w", err)
// 	}

// 	return filepath.Join(home, ".local/state/tidy/history.jsonl"), nil
// }
