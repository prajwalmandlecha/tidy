/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/prajwalmandlecha/tidy/history"
	"github.com/spf13/cobra"
)

var historyLimit int

// historyCmd represents the history command
var historyCmd = &cobra.Command{
	Use:   "history",
	Short: "Show the history of actions",

	RunE: func(cmd *cobra.Command, args []string) error {
		historyPath, err := history.DefaultPath()
		if err != nil {
			return fmt.Errorf("history: could not get default path: %w", err)
		}
		entries, err := history.Latest(historyPath, historyLimit)
		if err != nil {
			return fmt.Errorf("history: could not read history: %w", err)
		}
		if len(entries) == 0 {
			fmt.Println("No history entries found.")
			return nil
		}

		for _, entry := range entries {
			id := entry.ID
			record := entry.Record
			fmt.Printf("%d: %s -> %s (rule: %s, action: %s)\n",
				id, record.Source, record.Destination, record.Rule, record.Action)
		}
		return nil
	},
}

func init() {
	historyCmd.Flags().IntVar(&historyLimit, "limit", 20, "number of history entries to show")
	rootCmd.AddCommand(historyCmd)
}
