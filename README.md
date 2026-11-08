# 📦 ShareDrop

> Local network file sharing, done right.

ShareDrop is a fast, beautiful CLI tool written in Go that lets you share files between devices on your local network — no cloud, no accounts, no setup hassle.

## Features

- 🌐 **Web UI** — browse, download, and upload files from any browser on your network
- 🔒 **Password protection** — optional password stored locally with bcrypt hashing
- ⬆️ **Drag & drop upload** — upload files directly from the browser
- 🗜️ **ZIP download** — download entire folders as a single ZIP
- 🔍 **Live search** — filter files instantly in the browser
- 📱 **Mobile friendly** — responsive layout works on phones too
- 🎨 **Beautiful CLI** — colorful terminal output with ASCII art banner
- 🚀 **First-run wizard** — guided setup on first launch
- ⚙️ **Per-device config** — settings saved in `~/.sharedrop/config.json`
- 🏃 **Single binary** — no runtime deps, works on Windows, Linux, and macOS

## Install

Download the binary for your platform from [Releases](../../releases).

```bash
# Linux / macOS
chmod +x sharedrop-linux-amd64
mv sharedrop-linux-amd64 /usr/local/bin/sharedrop
```

## Usage

```bash
# Share the current folder
sharedrop

# Share a specific folder
sharedrop -d /path/to/folder
sharedrop /path/to/folder

# Custom port
sharedrop -p 9000

# Skip password even if one is saved
sharedrop --no-password
```

## Commands

| Command | Description |
|---|---|
| `sharedrop` | Start sharing (first run triggers setup wizard) |
| `sharedrop setup` | Re-run the setup wizard |
| `sharedrop config` | Show current configuration |
| `sharedrop passwd` | Set or remove the share password |
| `sharedrop version` | Show version info |

## Configuration

Config lives at `~/.sharedrop/config.json`. You can edit it manually or use `sharedrop setup`.

```json
{
  "password_hash": "...",
  "password_set": true,
  "default_port": 8080,
  "allow_upload": true,
  "allow_delete": false,
  "show_hidden": false,
  "max_upload_mb": 500
}
```

## Build from source

```bash
git clone https://github.com/user/sharedrop
cd sharedrop
go build -o sharedrop ./
```

## License

MIT
