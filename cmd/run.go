/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/prajwalmandlecha/tidy/config"
	"github.com/prajwalmandlecha/tidy/engine"
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
				if err := engine.ProcessFile(path, cfg.Rules, dryRun); err != nil {
					fmt.Printf("error processing %q: %v\n", path, err)
				}
			}
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(runCmd)
}
