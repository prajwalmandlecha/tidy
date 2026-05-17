/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/prajwalmandlecha/tidy/config"
	"github.com/spf13/cobra"
)

// validateCmd represents the validate command
var validateCmd = &cobra.Command{
	Use:   "validate",
	Short: "Validate the configuration",
	Long:  `Validate the configuration file for syntax and semantic errors.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		_, err := config.Load(cfgFile)
		if err != nil {
			return fmt.Errorf("invalid config: %w", err)
		}

		fmt.Println("config is valid")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(validateCmd)
}
