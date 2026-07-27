/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/prajwalmandlecha/tidy/history"
	"github.com/spf13/cobra"
)

var undoID int

// undoCmd represents the undo command
var undoCmd = &cobra.Command{
	Use:   "undo",
	Short: "Undo a recorded move",
	RunE: func(cmd *cobra.Command, args []string) error {
		db, err := history.NewDB()
		if err != nil {
			return fmt.Errorf("could not open history database: %w", err)
		}
		defer db.Close()

		var targetMove history.Move

		if undoID > 0 {
			move, found, err := db.FindPending(int64(undoID))
			if err != nil {
				return fmt.Errorf("could not get move record: %w", err)
			}
			if !found {
				fmt.Printf("No pending move entry found with ID %d.\n", undoID)
				return nil
			}
			targetMove = move
		} else {
			pending, err := db.Pending(1)
			if err != nil {
				return fmt.Errorf("could not get pending moves: %w", err)
			}
			if len(pending) == 0 {
				fmt.Println("No actions to undo.")
				return nil
			}
			targetMove = pending[0]
		}

		_, err = os.Stat(targetMove.Source)
		if err == nil {
			return fmt.Errorf("source file %q already exists, cannot undo action", targetMove.Source)
		} else if !os.IsNotExist(err) {
			return fmt.Errorf("could not check source path %q: %w", targetMove.Source, err)
		}

		_, err = os.Stat(targetMove.Destination)
		if err != nil {
			return fmt.Errorf("could not find file to restore %q: %w", targetMove.Destination, err)
		}

		if dryRun {
			fmt.Printf("[dry-run] would move %q back to %q\n", targetMove.Destination, targetMove.Source)
			return nil
		}

		err = os.Rename(targetMove.Destination, targetMove.Source)
		if err != nil {
			return fmt.Errorf("could not undo action: %w", err)
		}

		if err := db.MarkUndone(targetMove.ID, time.Now()); err != nil {
			return fmt.Errorf("could not mark move as undone in database: %w", err)
		}

		fmt.Printf("Undid action (ID %d): moved %q back to %q\n", targetMove.ID, targetMove.Destination, targetMove.Source)
		return nil
	},
}

func init() {
	undoCmd.Flags().IntVar(&undoID, "id", 0, "history entry ID to undo")
	rootCmd.AddCommand(undoCmd)
}
