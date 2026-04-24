/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"fmt"
	"os/signal"
	"syscall"

	"github.com/prajwalmandlecha/tidy/config"
	"github.com/prajwalmandlecha/tidy/watcher"
	"github.com/spf13/cobra"
)

// watchCmd represents the watch command
var watchCmd = &cobra.Command{
	Use:   "watch",
	Short: "Watch folders and organize files as they arrive",

	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load(cfgFile)
		if err != nil {
			return fmt.Errorf("failed to load config: %w", err)
		}

		ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
		defer stop()

		fmt.Println("tidy is watching... (ctrl+c to stop)")
		return watcher.Watch(ctx, cfg, dryRun)
	},
}

func init() {
	rootCmd.AddCommand(watchCmd)
}
