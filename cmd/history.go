/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/prajwalmandlecha/tidy/history"
	"github.com/spf13/cobra"
)

// historyCmd represents the history command
var historyCmd = &cobra.Command{
	Use:   "history",
	Short: "Show the history of actions",

	RunE: func(cmd *cobra.Command, args []string) error {
		historyPath, err := history.DefaultPath()
		if err != nil {
			return fmt.Errorf("history: could not get default path: %w", err)
		}
		err = history.PrintHistory(historyPath)
		if err != nil {
			return fmt.Errorf("history: could not print history: %w", err)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(historyCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// historyCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// historyCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
