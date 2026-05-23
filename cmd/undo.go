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

// undoCmd represents the undo command
var undoCmd = &cobra.Command{
	Use:   "undo",
	Short: "Undo the last action",
	RunE: func(cmd *cobra.Command, args []string) error {

		historyPath, err := history.DefaultPath()
		if err != nil {
			return fmt.Errorf("could not get history path: %w", err)
		}

		lastRecord, found, err := history.Last(historyPath)
		if err != nil {
			return fmt.Errorf("could not get last record: %w", err)
		}
		if !found {
			fmt.Println("No actions to undo.")
			return nil
		}

		_, err = os.Stat(lastRecord.Source)
		if err == nil {
			return fmt.Errorf("source file %q already exists, cannot undo last action", lastRecord.Source)
		} else if !os.IsNotExist(err) {
			return fmt.Errorf("could not check source path %q: %w", lastRecord.Source, err)
		}

		if dryRun {
			fmt.Printf("[dry-run] would move %q back to %q\n", lastRecord.Destination, lastRecord.Source)
			return nil
		}

		_, err = os.Stat(lastRecord.Destination)
		if err != nil {
			return fmt.Errorf("could not find file to restore %q: %w", lastRecord.Destination, err)
		}

		err = os.Rename(lastRecord.Destination, lastRecord.Source)
		if err != nil {
			return fmt.Errorf("could not undo last action: %w", err)
		}

		fmt.Printf("Undid last action: moved %q back to %q\n", lastRecord.Destination, lastRecord.Source)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(undoCmd)
}
