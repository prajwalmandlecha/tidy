package watcher

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/prajwalmandlecha/tidy/config"
	"github.com/prajwalmandlecha/tidy/engine"
	"github.com/prajwalmandlecha/tidy/history"
)

const debounceDelay = 300 * time.Millisecond

func Watch(ctx context.Context, cfg *config.Config, dryRun bool) error {
	db, err := history.NewDB()
	if err != nil {
		return fmt.Errorf("watcher: failed to open history database: %w", err)
	}
	defer db.Close()

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return fmt.Errorf("watcher: failed to create: %w", err)
	}
	defer watcher.Close()

	for _, dir := range cfg.WatchDirs {
		err = watcher.Add(dir)
		if err != nil {
			return fmt.Errorf("watcher: failed to watch %q: %w", dir, err)
		}
		fmt.Printf("watching %s\n", dir)

	}

	timer := make(map[string]*time.Timer)
	mu := sync.Mutex{}

	for {
		select {
		case <-ctx.Done():

			return nil
		case event, ok := <-watcher.Events:
			if !ok {
				return nil
			}
			if !event.Has(fsnotify.Create) &&
				!event.Has(fsnotify.Rename) &&
				!event.Has(fsnotify.Write) {
				continue
			}
			filePath := event.Name
			if engine.IsIgnored(filePath, cfg.Ignore) {
				continue
			}

			mu.Lock()
			if t, exists := timer[filePath]; exists {
				t.Stop()
				t.Reset(debounceDelay)
			} else {
				timer[filePath] = time.AfterFunc(debounceDelay, func() {
					mu.Lock()
					delete(timer, filePath)
					mu.Unlock()

					res, err := engine.ProcessFile(filePath, cfg.Ignore, cfg.Rules, dryRun)
					if err != nil {
						fmt.Printf("error processing %q: %v\n", filePath, err)
						return
					}

					if res == nil || dryRun {
						return
					}

					_, err = db.Append(history.Move{
						Source:      res.Source,
						Destination: res.Destination,
						Rule:        res.Rule,
						Action:      string(res.Action),
						MovedAt:     time.Now(),
					})
					if err != nil {
						fmt.Printf("error writing history: %v\n", err)
					}
				})
			}

			mu.Unlock()
			fmt.Println("event:", filePath)
		case err, ok := <-watcher.Errors:
			if !ok {
				return nil
			}

			fmt.Printf("watcher error: %v\n", err)
		}

	}
}
