/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/prajwalmandlecha/tidy/history"
	"github.com/spf13/cobra"
)

var historyLimit int

// historyCmd represents the history command
var historyCmd = &cobra.Command{
	Use:   "history",
	Short: "Show the history of actions",

	RunE: func(cmd *cobra.Command, args []string) error {
		db, err := history.NewDB()
		if err != nil {
			return fmt.Errorf("history: could not open database: %w", err)
		}
		defer db.Close()

		moves, err := db.Latest(historyLimit)
		if err != nil {
			return fmt.Errorf("history: could not read history: %w", err)
		}
		if len(moves) == 0 {
			fmt.Println("No history entries found.")
			return nil
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "ID\tTIMESTAMP\tRULE\tACTION\tSTATUS\tMOVE")

		for _, m := range moves {
			status := "OK"
			if m.UndoneAt != nil {
				status = "UNDONE"
			}
			ts := m.MovedAt.Local().Format("2006-01-02 15:04:05")
			fmt.Fprintf(w, "%d\t%s\t%s\t%s\t%s\t%s -> %s\n",
				m.ID, ts, m.Rule, m.Action, status, m.Source, m.Destination)
		}
		w.Flush()
		return nil
	},
}

func init() {
	historyCmd.Flags().IntVar(&historyLimit, "limit", 20, "number of history entries to show")
	rootCmd.AddCommand(historyCmd)
}
