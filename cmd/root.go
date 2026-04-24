/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var cfgFile string
var dryRun bool

var rootCmd = &cobra.Command{
	Use:   "tidy",
	Short: "Automatically organize your folders based on rules",
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "~/.config/tidy/config.yaml", "config path")
	rootCmd.PersistentFlags().BoolVar(&dryRun, "dry-run", false, "preview without moving files")
}
