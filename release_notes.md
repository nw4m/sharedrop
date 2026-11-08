## ShareDrop v1.1.0

Local network file sharing — single binary, no setup.

### New in this release

- Dark / Light mode toggle (saved in browser)
- QR code in terminal — scan with phone to open instantly
- Upload progress bar with live transfer speed
- Fixed upload drag & drop on all browsers
- Filesystem root browsing with `--root` flag
- Self-install to PATH with `install` command

### Install

**Windows** (PowerShell):
```powershell
.\sharedrop-windows-amd64.exe install
sharedrop
```

**Linux / macOS**:
```bash
chmod +x sharedrop-linux-amd64
sudo mv sharedrop-linux-amd64 /usr/local/bin/sharedrop
sharedrop install
sharedrop
```

### All platforms

| Platform | File |
|---|---|
| Windows x64 | sharedrop-windows-amd64.exe |
| Windows ARM | sharedrop-windows-arm64.exe |
| Linux x64 | sharedrop-linux-amd64 |
| Linux ARM | sharedrop-linux-arm64 |
| macOS Intel | sharedrop-darwin-amd64 |
| macOS Apple Silicon | sharedrop-darwin-arm64 |
