/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

const starterConfig = `watch_dirs:
  - ~/Downloads
  - ~/Pictures

rules:
  - name: Documents
    extensions: [".pdf", ".doc", ".docx", ".txt", ".rtf", ".odt", ".md", ".csv", ".xls", ".xlsx", ".ppt", ".pptx"]
    dest: ~/Downloads/Documents

  - name: Images
    extensions: [".png", ".jpg", ".jpeg", ".gif", ".webp", ".svg"]
    dest: ~/Downloads/Images

  - name: Video
    extensions: [".mp4", ".mkv", ".mov", ".avi", ".webm", ".m4v"]
    dest: ~/Downloads/Video

  - name: Audio
    extensions: [".mp3", ".wav", ".flac", ".aac", ".ogg", ".m4a"]
    dest: ~/Downloads/Music

  - name: Archives
    extensions: [".zip", ".rar", ".7z", ".tar", ".tar.gz", ".tgz", ".gz", ".bz2", ".xz"]
    dest: ~/Downloads/Compressed

  - name: Programs
    extensions: [".exe", ".msi", ".deb", ".rpm", ".dmg", ".pkg", ".apk", ".appimage", ".jar"]
    dest: ~/Downloads/Programs

  - name: Code
    extensions: [".ipynb", ".html", ".css", ".js", ".ts", ".tsx", ".jsx", ".go", ".py", ".java", ".c", ".cpp", ".h", ".hpp", ".yaml", ".yml", ".xml", ".skill"]
    dest: ~/Downloads/Code

  - name: Configs
    extensions: [".json", ".pem", ".key", ".ovpn", ".toml", ".ini", ".cfg"]
    dest: ~/Downloads/Configs

  - name: Screenshots
    pattern: "Screenshot*"
    dest: ~/Pictures/Screenshots

  - name: Torrents
    extensions: [".torrent"]
    dest: ~/Downloads/Torrents
`

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize a new tidy configuration",
	RunE: func(cmd *cobra.Command, args []string) error {
		path := cfgFile

		if len(path) > 1 && path[:2] == "~/" {
			home, err := os.UserHomeDir()
			if err != nil {
				return err
			}
			path = filepath.Join(home, path[2:])
		}

		if _, err := os.Stat(path); err == nil {
			fmt.Printf("config already exists at %s — skipping\n", path)
			return nil

		}

		if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
			return fmt.Errorf("could not create config dir: %w", err)

		}
		if err := os.WriteFile(path, []byte(starterConfig), 0644); err != nil {
			return fmt.Errorf("could not write config: %w", err)

		}

		fmt.Printf("config written to %s\n", path)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(initCmd)

}
