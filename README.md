# ShareDrop

Local network file sharing, done right. Single binary, no setup, works on every device on your WiFi.

## Features

- Web UI — browse, download, upload, delete, ZIP folders
- Password protection with bcrypt
- Dark / Light mode toggle
- QR code in terminal — scan with your phone instantly
- Upload progress bar
- Drag & drop upload
- Live file search
- Sort by name, size, date
- JSON API at `/api/files`
- Run from filesystem root with `--root`
- `install` command adds itself to PATH
- Cross-platform: Windows, Linux, macOS

## Install

Download the binary for your platform from [Releases](https://github.com/nw4m/sharedrop/releases).

```bash
# Linux / macOS
chmod +x sharedrop-linux-amd64
sudo mv sharedrop-linux-amd64 /usr/local/bin/sharedrop
sharedrop install
```

```powershell
# Windows — run in PowerShell
.\sharedrop-windows-amd64.exe install
```

## Usage

```bash
# Share current folder (first run triggers setup wizard)
sharedrop

# Share a specific folder
sharedrop -d /path/to/folder

# Share entire filesystem
sharedrop --root

# Custom port
sharedrop -p 9000

# Skip password
sharedrop --no-password
```

## Commands

| Command | Description |
|---|---|
| `sharedrop` | Start sharing |
| `sharedrop install` | Add to PATH permanently |
| `sharedrop setup` | Re-run setup wizard |
| `sharedrop config` | Show current config |
| `sharedrop passwd` | Change / remove password |
| `sharedrop version` | Version info |

## Build from source

```bash
git clone https://github.com/nw4m/sharedrop
cd sharedrop
go build -o sharedrop ./
```

## License

MIT
