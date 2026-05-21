/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/prajwalmandlecha/tidy/config"
	"github.com/prajwalmandlecha/tidy/engine"
	"github.com/prajwalmandlecha/tidy/history"
	"github.com/spf13/cobra"
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Organize existing files once and exit",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load(cfgFile)
		if err != nil {
			return fmt.Errorf("failed to load config: %w", err)
		}

		for _, dir := range cfg.WatchDirs {
			entries, err := os.ReadDir(dir)
			if err != nil {
				fmt.Printf("could not read %s: %v\n", dir, err)
				continue
			}

			for _, entry := range entries {
				if entry.IsDir() {
					continue
				}

				path := filepath.Join(dir, entry.Name())
				res, err := engine.ProcessFile(path, cfg.Rules, dryRun)
				if err != nil {
					fmt.Printf("error processing %q: %v\n", path, err)
					continue
				}

				if res == nil || dryRun {
					continue
				}

				historyPath, err := history.DefaultPath()
				if err != nil {
					fmt.Printf("error getting history path: %v\n", err)
					continue
				}

				err = history.Append(historyPath, history.Record{
					Source:      res.Source,
					Destination: res.Destination,
					Rule:        res.Rule,
					Action:      string(res.Action),
					Timestamp:   time.Now(),
				})
				if err != nil {
					fmt.Printf("error writing history: %v\n", err)
				}
			}
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(runCmd)
}
