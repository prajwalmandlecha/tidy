/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/prajwalmandlecha/tidy/history"
	"github.com/spf13/cobra"
)

var undoID int

// undoCmd represents the undo command
var undoCmd = &cobra.Command{
	Use:   "undo",
	Short: "Undo a recorded move",
	RunE: func(cmd *cobra.Command, args []string) error {

		historyPath, err := history.DefaultPath()
		if err != nil {
			return fmt.Errorf("could not get history path: %w", err)
		}

		var record history.Record

		if undoID > 0 {
			entry, found, err := history.Find(historyPath, undoID)
			if err != nil {
				return fmt.Errorf("could not get the record: %w", err)
			}
			if !found {
				fmt.Printf("No history entry found with ID %d.\n", undoID)
				return nil
			}
			record = entry.Record
		} else {
			lastRecord, found, err := history.Last(historyPath)
			if err != nil {
				return fmt.Errorf("could not get last record: %w", err)
			}
			if !found {
				fmt.Println("No actions to undo.")
				return nil
			}
			record = lastRecord
		}

		_, err = os.Stat(record.Source)
		if err == nil {
			return fmt.Errorf("source file %q already exists, cannot undo last action", record.Source)
		} else if !os.IsNotExist(err) {
			return fmt.Errorf("could not check source path %q: %w", record.Source, err)
		}

		_, err = os.Stat(record.Destination)
		if err != nil {
			return fmt.Errorf("could not find file to restore %q: %w", record.Destination, err)
		}

		if dryRun {
			fmt.Printf("[dry-run] would move %q back to %q\n", record.Destination, record.Source)
			return nil
		}

		err = os.Rename(record.Destination, record.Source)
		if err != nil {
			return fmt.Errorf("could not undo last action: %w", err)
		}

		fmt.Printf("Undid last action: moved %q back to %q\n", record.Destination, record.Source)
		return nil
	},
}

func init() {
	undoCmd.Flags().IntVar(&undoID, "id", 0, "history entry ID to undo")
	rootCmd.AddCommand(undoCmd)
}
