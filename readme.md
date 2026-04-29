# tidy

`tidy` is a small Go CLI that keeps folders organized with simple YAML rules.
It is currently tuned for an IDM-style `~/Downloads` workflow: documents,
images, videos, archives, programs, code, configs, torrents, and similar files
can be moved into predictable buckets inside `~/Downloads`.

The tool can run once over existing files, or watch a folder and organize new
files as they arrive.

## Features

- One-shot organization with `tidy run`
- Live folder watching with `tidy watch`
- Safe preview mode with `--dry-run`
- Safely skips identical duplicate files
- Automatically renames files with a number suffix on filename collisions
- YAML config stored by default at `~/.config/tidy/config.yaml`
- Case-insensitive extension matching
- Multi-part extension matching such as `.tar.gz`
- Optional glob-style filename patterns
- Top-level-only scanning so already-organized subfolders are left alone

## Quick Start

Build the binary:

```sh
go build -o tidy .
```

Create the starter config:

```sh
./tidy init
```

Preview what would move:

```sh
./tidy run --dry-run
```

Apply the rules once:

```sh
./tidy run
```

Watch for new files:

```sh
./tidy watch
```

Use a custom config when needed:

```sh
./tidy --config ./config.yaml run --dry-run
```

## Commands

| Command                | What it does                                                      |
| ---------------------- | ----------------------------------------------------------------- |
| `tidy init`            | Writes a starter config unless one already exists                 |
| `tidy run`             | Organizes existing top-level files in each watched directory      |
| `tidy run --dry-run`   | Prints planned moves without changing files                       |
| `tidy watch`           | Watches configured directories and processes new or renamed files |
| `tidy watch --dry-run` | Logs what would move while watching                               |

Global flags:

| Flag        | Default                      | Purpose                              |
| ----------- | ---------------------------- | ------------------------------------ |
| `--config`  | `~/.config/tidy/config.yaml` | Path to the YAML config              |
| `--dry-run` | `false`                      | Preview moves without touching files |

## Configuration

The default config watches `~/Downloads` and moves files into buckets inside
that same directory.

```yaml
watch_dirs:
  - ~/Downloads

rules:
  - name: Documents
    extensions: [".pdf", ".doc", ".docx", ".txt", ".md", ".csv"]
    dest: ~/Downloads/Documents

  - name: Images
    extensions: [".png", ".jpg", ".jpeg", ".gif", ".webp", ".svg"]
    dest: ~/Downloads/Images

  - name: Archives
    extensions: [".zip", ".rar", ".7z", ".tar", ".tar.gz", ".tgz"]
    dest: ~/Downloads/Compressed

  - name: Torrents
    extensions: [".torrent"]
    dest: ~/Downloads/Torrents
```

Each rule supports:

| Field        | Required                        | Description                                             |
| ------------ | ------------------------------- | ------------------------------------------------------- |
| `name`       | Recommended                     | Human-readable rule name used in validation messages    |
| `extensions` | Yes, unless `pattern` is set    | File extensions to match, with or without a leading dot |
| `pattern`    | Yes, unless `extensions` is set | Glob pattern matched against the base filename          |
| `dest`       | Yes                             | Destination directory                                   |

Rules are checked in order. The first matching rule wins.

## Matching Behavior

`tidy` matches extensions case-insensitively, so `.PDF` and `.pdf` are treated
the same. It also checks full filename suffixes, which means rules like
`.tar.gz` work as expected.

Pattern rules use Go's `filepath.Match` against the filename only:

```yaml
rules:
  - name: Screenshots
    pattern: "Screenshot*"
    dest: ~/Downloads/Images/Screenshots
```

## Safety Notes

Start with a dry run:

```sh
./tidy run --dry-run
```

`tidy run` scans only the direct files inside each `watch_dirs` entry and skips
subdirectories. This keeps existing organized folders such as
`~/Downloads/Documents` or extracted archives from being flattened.

Destination directories are created automatically. If a source file is already
in its destination folder, it is skipped. If a file with the same name already
exists in the destination:

- It is skipped completely if the file content hashes are identical.
- It is automatically renamed with a number suffix (e.g., `file_1.txt`) if the contents differ preventing accidental overwrites.

## Development

Run tests:

```sh
go test ./...
```

Build:

```sh
make build
```

Code layout:

```text
.
|-- cmd/       # Cobra commands: init, run, watch
|-- config/    # YAML loading, path expansion, validation
|-- engine/    # Rule matching and file moves
|-- watcher/   # fsnotify integration with debounce
|-- config.yaml
|-- go.mod
`-- main.go
```

General Make targets:

| Target  | Command                  |
| ------- | ------------------------ |
| `build` | Builds the `tidy` binary |
| `test`  | Runs `go test -v ./...`  |

### Linux User Service Helpers

The Makefile also includes convenience targets for running `tidy` as a Linux
user service through systemd. These are optional; use the normal CLI commands
above if you just want to run `tidy` manually.

These targets assume:

- Linux with systemd
- `sudo` access for installing/removing `/usr/local/bin/tidy`
- A matching user service named `tidy` is already installed and enabled, for
  example `~/.config/systemd/user/tidy.service`

| Target      | Command                                                                                             |
| ----------- | --------------------------------------------------------------------------------------------------- |
| `install`   | Builds `tidy`, moves it to `/usr/local/bin/tidy`, reloads user systemd, and restarts `tidy.service` |
| `logs`      | Follows logs for the user service with `journalctl --user -u tidy -f`                               |
| `status`    | Shows `systemctl --user status tidy`                                                                |
| `uninstall` | Stops/disables the user service and removes `/usr/local/bin/tidy`                                   |
