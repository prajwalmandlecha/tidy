# tidy

`tidy` is a small Go CLI that keeps folders organized with simple YAML rules.

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

Validate the config:

```sh
./tidy validate
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
| `tidy validate`        | Checks config syntax, watched folders, and rule paths             |
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

  - name: Screenshots
    pattern: "Screenshot*"
    dest: ~/Downloads/Images/Screenshots

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
- It is automatically renamed with a number suffix (e.g., `file_1.txt`) if the
  contents differ, preventing accidental overwrites.

## Linux Systemd User Service

Linux users who want `tidy` to run in the background can create a user service
at `~/.config/systemd/user/tidy.service`:

Build and install the binary first:

```sh
go build -o tidy .
sudo install -m 0755 tidy /usr/local/bin/tidy
```

```ini
[Unit]
Description=Tidy - automatic folder organizer
After=graphical-session.target

[Service]
Type=simple
ExecStart=/usr/local/bin/tidy watch
Restart=on-failure
RestartSec=5s

[Install]
WantedBy=default.target
```

This uses the default config path, `~/.config/tidy/config.yaml`. For a custom
config file, pass `--config` before `watch`:

```ini
ExecStart=/usr/local/bin/tidy --config /path/to/config.yaml watch
```

Reload systemd and enable the service:

```sh
systemctl --user daemon-reload
systemctl --user enable --now tidy.service
```

`enable --now` starts the service immediately and enables it to start
automatically in future user sessions.

Restart or disable the service when needed:

```sh
systemctl --user restart tidy.service
systemctl --user disable --now tidy.service
```

Remove the service file completely:

```sh
rm ~/.config/systemd/user/tidy.service
systemctl --user daemon-reload
```

## Development

Run tests:

```sh
go test ./...
```

Build:

```sh
go build -o tidy .
```

Code layout:

```text
.
|-- cmd/       # Cobra commands: init, validate, run, watch
|-- config/    # YAML loading, path expansion, validation
|-- engine/    # Rule matching and file moves
|-- watcher/   # fsnotify integration with debounce
|-- config.yaml
|-- go.mod
`-- main.go
```
