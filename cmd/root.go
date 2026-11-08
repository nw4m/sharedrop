package cmd

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
	"github.com/user/sharedrop/internal/auth"
	"github.com/user/sharedrop/internal/config"
	"github.com/user/sharedrop/internal/ui"
)

var rootCmd = &cobra.Command{
	Use:   "sharedrop",
	Short: "ShareDrop - Local network file sharing",
	Long:  "ShareDrop lets you share files on your local network with a slick web UI.",
	RunE: func(cmd *cobra.Command, args []string) error {
		return runShare(cmd, args)
	},
}

var setupCmd = &cobra.Command{
	Use:   "setup",
	Short: "Run the setup wizard again",
	RunE: func(cmd *cobra.Command, args []string) error {
		ui.PrintBanner()
		return runSetup()
	},
}

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Show current configuration",
	RunE: func(cmd *cobra.Command, args []string) error {
		ui.PrintBanner()
		cfg, err := config.Load()
		if err != nil {
			ui.Error("Failed to load config: " + err.Error())
			return err
		}
		ui.PrintConfig(cfg.DefaultPort, cfg.AllowUpload, cfg.AllowDelete, cfg.ShowHidden, cfg.MaxUploadMB, cfg.PasswordSet)
		return nil
	},
}

var passwdCmd = &cobra.Command{
	Use:   "passwd",
	Short: "Set or remove share password",
	RunE: func(cmd *cobra.Command, args []string) error {
		ui.PrintBanner()
		return runPasswd()
	},
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version info",
	Run: func(cmd *cobra.Command, args []string) {
		ui.PrintBanner()
		ui.Info("ShareDrop v1.0.0")
		ui.Dim("Built with Go · MIT License")
		ui.Dim("github.com/user/sharedrop")
	},
}

var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Install ShareDrop system-wide (adds to PATH)",
	RunE: func(cmd *cobra.Command, args []string) error {
		ui.PrintBanner()
		return runInstall()
	},
}

var uninstallCmd = &cobra.Command{
	Use:   "uninstall",
	Short: "Remove ShareDrop from PATH and installed location",
	RunE: func(cmd *cobra.Command, args []string) error {
		ui.PrintBanner()
		return runUninstall()
	},
}

func init() {
	rootCmd.PersistentFlags().IntP("port", "p", 0, "Port to listen on (overrides config)")
	rootCmd.PersistentFlags().StringP("dir", "d", "", "Directory to share (current dir if not set)")
	rootCmd.PersistentFlags().BoolP("no-password", "n", false, "Skip password even if one is set")

	rootCmd.AddCommand(setupCmd)
	rootCmd.AddCommand(configCmd)
	rootCmd.AddCommand(passwdCmd)
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(installCmd)
	rootCmd.AddCommand(uninstallCmd)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func runShare(cmd *cobra.Command, args []string) error {
	ui.PrintBanner()

	cfg, err := config.Load()
	if err != nil {
		ui.Error("Config error: " + err.Error())
		return err
	}

	// first run: setup + offer to install
	dir, _ := config.ConfigDir()
	if _, statErr := os.Stat(dir + "/config.json"); os.IsNotExist(statErr) {
		ui.Info("Looks like your first time! Let's do a quick setup.")
		fmt.Println()

		// offer to install before setup so next time it works from anywhere
		if !isInstalledInPath() {
			ui.Title("Install ShareDrop")
			fmt.Println()
			ui.Dim("Install ShareDrop so you can run it from any terminal?")
			fmt.Println()
			if ui.PromptDefault("Install to PATH now?", true) {
				if err := runInstall(); err != nil {
					ui.Warn("Install failed, continuing anyway...")
				}
				fmt.Println()
			}
		}

		if err := runSetup(); err != nil {
			return err
		}
		// reload after setup
		cfg, err = config.Load()
		if err != nil {
			return err
		}
	}

	// apply CLI flag overrides
	if portFlag, _ := cmd.Flags().GetInt("port"); portFlag != 0 {
		cfg.DefaultPort = portFlag
	}

	shareDir, _ := cmd.Flags().GetString("dir")
	if shareDir == "" {
		if len(args) > 0 {
			shareDir = args[0]
		} else if cfg.DefaultFolder != "" {
			shareDir = cfg.DefaultFolder
		} else {
			// default to current directory
			shareDir, _ = os.Getwd()
		}
	}

	// validate directory
	info, err := os.Stat(shareDir)
	if err != nil || !info.IsDir() {
		ui.Error("Directory not found: " + shareDir)
		return fmt.Errorf("invalid directory: %s", shareDir)
	}

	noPass, _ := cmd.Flags().GetBool("no-password")
	if noPass {
		cfg.PasswordSet = false
	}

	return startServer(cfg, shareDir)
}

func runSetup() error {
	ui.PrintSetupHeader()

	cfg, err := config.Load()
	if err != nil {
		cfg = config.DefaultConfig()
	}

	// Step 1: Password
	ui.PrintStep(1, "Password Protection")
	fmt.Println()
	ui.Dim("Do you want to require a password when someone accesses ShareDrop?")
	fmt.Println()

	if ui.PromptDefault("Enable password protection?", false) {
		pass := ui.PasswordInput("Enter password")
		if pass == "" {
			ui.Warn("Empty password — skipping password protection")
		} else {
			confirm := ui.PasswordInput("Confirm password")
			if pass != confirm {
				ui.Error("Passwords don't match — password protection disabled")
			} else {
				hash, err := auth.HashPassword(pass)
				if err != nil {
					ui.Error("Failed to hash password")
				} else {
					cfg.PasswordHash = hash
					cfg.PasswordSet = true
					ui.Success("Password set!")
				}
			}
		}
	} else {
		cfg.PasswordSet = false
		cfg.PasswordHash = ""
		ui.Dim("No password — anyone on your network can access files")
	}

	// Step 2: Port
	fmt.Println()
	ui.PrintStep(2, "Network Port")
	fmt.Println()
	ui.Dim("Which port should ShareDrop listen on?")
	fmt.Println()

	portStr := ui.Input("Port number", strconv.Itoa(cfg.DefaultPort))
	if p, err := strconv.Atoi(portStr); err == nil && p > 0 && p < 65536 {
		cfg.DefaultPort = p
	} else {
		ui.Warn("Invalid port, keeping default: " + strconv.Itoa(cfg.DefaultPort))
	}

	// Step 3: Permissions
	fmt.Println()
	ui.PrintStep(3, "Permissions")
	fmt.Println()

	cfg.AllowUpload = ui.PromptDefault("Allow others to upload files to your device?", true)
	cfg.AllowDelete = ui.Prompt("Allow others to delete files? (not recommended)")
	cfg.ShowHidden = ui.Prompt("Show hidden files (dotfiles)?")

	// Step 4: Upload limit
	fmt.Println()
	ui.PrintStep(4, "Upload Limit")
	fmt.Println()
	ui.Dim("Max size per upload in MB (e.g. 500 = 500MB)")
	fmt.Println()

	limitStr := ui.Input("Max upload MB", strconv.Itoa(cfg.MaxUploadMB))
	if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
		cfg.MaxUploadMB = l
	}

	// save
	if err := cfg.Save(); err != nil {
		ui.Error("Failed to save config: " + err.Error())
		return err
	}

	fmt.Println()
	ui.Success("Setup complete! Config saved to ~/.sharedrop/config.json")
	fmt.Println()
	ui.PrintConfig(cfg.DefaultPort, cfg.AllowUpload, cfg.AllowDelete, cfg.ShowHidden, cfg.MaxUploadMB, cfg.PasswordSet)

	return nil
}

func runPasswd() error {
	cfg, err := config.Load()
	if err != nil {
		ui.Error("Failed to load config")
		return err
	}

	ui.Title("Password Management")
	fmt.Println()

	if cfg.PasswordSet {
		ui.Info("A password is currently set.")
		if ui.Prompt("Remove existing password?") {
			cfg.PasswordSet = false
			cfg.PasswordHash = ""
			cfg.Save()
			ui.Success("Password removed — share is now open")
			return nil
		}
		ui.Info("Setting a new password...")
	}

	pass := ui.PasswordInput("New password")
	if pass == "" {
		ui.Warn("Empty input — no changes made")
		return nil
	}

	confirm := ui.PasswordInput("Confirm password")
	if pass != confirm {
		ui.Error("Passwords don't match — no changes made")
		return nil
	}

	hash, err := auth.HashPassword(pass)
	if err != nil {
		ui.Error("Failed to hash password")
		return err
	}

	cfg.PasswordHash = hash
	cfg.PasswordSet = true

	if err := cfg.Save(); err != nil {
		ui.Error("Failed to save config")
		return err
	}

	ui.Success("Password updated!")
	return nil
}

// installDir returns the directory where we install the binary
func installDir() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".local", "bin")
}

// installDirWindows returns the Windows-specific install path
func installDirWindows() string {
	// prefer AppData/Local/Programs — no admin needed
	appData := os.Getenv("APPDATA")
	if appData == "" {
		home, _ := os.UserHomeDir()
		appData = filepath.Join(home, "AppData", "Roaming")
	}
	// use LocalAppData instead (better practice)
	localAppData := strings.Replace(appData, "Roaming", "Local", 1)
	return filepath.Join(localAppData, "Programs", "sharedrop")
}

func runInstall() error {
	// figure out where this exe currently lives
	selfPath, err := os.Executable()
	if err != nil {
		ui.Error("Cannot find current executable: " + err.Error())
		return err
	}
	selfPath, _ = filepath.EvalSymlinks(selfPath)

	ui.Title("Installing ShareDrop")
	fmt.Println()

	// pick install directory based on OS
	var destDir string
	if runtime.GOOS == "windows" {
		destDir = installDirWindows()
	} else {
		destDir = installDir()
	}

	// create the install dir
	if err := os.MkdirAll(destDir, 0755); err != nil {
		ui.Error("Cannot create install directory: " + err.Error())
		return err
	}

	// destination binary path
	binaryName := "sharedrop"
	if runtime.GOOS == "windows" {
		binaryName = "sharedrop.exe"
	}
	destPath := filepath.Join(destDir, binaryName)

	// don't copy over ourselves if we're already installed there
	selfAbs, _ := filepath.Abs(selfPath)
	destAbs, _ := filepath.Abs(destPath)
	if strings.EqualFold(selfAbs, destAbs) {
		ui.Success("Already installed at: " + destPath)
	} else {
		// copy the binary
		ui.Info("Copying binary to: " + destPath)
		if err := copyFile(selfPath, destPath); err != nil {
			ui.Error("Copy failed: " + err.Error())
			return err
		}
		// make executable on unix
		if runtime.GOOS != "windows" {
			os.Chmod(destPath, 0755)
		}
		ui.Success("Binary installed!")
	}

	// add to PATH
	fmt.Println()
	added, err := addToPath(destDir)
	if err != nil {
		ui.Warn("Could not update PATH automatically: " + err.Error())
		ui.Dim("Add this to your PATH manually: " + destDir)
	} else if added {
		ui.Success("Added to PATH: " + destDir)
		fmt.Println()
		if runtime.GOOS == "windows" {
			ui.Dim("Open a new terminal window and type: sharedrop")
		} else {
			ui.Dim("Run: source ~/.bashrc  (or open a new terminal)")
			ui.Dim("Then type: sharedrop")
		}
	} else {
		ui.Info("PATH already contains: " + destDir)
	}

	fmt.Println()
	ui.Success("Installation complete!")
	ui.Dim("Run 'sharedrop' from any terminal to start sharing")
	fmt.Println()
	return nil
}

func runUninstall() error {
	ui.Title("Uninstalling ShareDrop")
	fmt.Println()

	var destDir string
	if runtime.GOOS == "windows" {
		destDir = installDirWindows()
	} else {
		destDir = installDir()
	}

	binaryName := "sharedrop"
	if runtime.GOOS == "windows" {
		binaryName = "sharedrop.exe"
	}
	destPath := filepath.Join(destDir, binaryName)

	// remove the binary
	if _, err := os.Stat(destPath); err == nil {
		if err := os.Remove(destPath); err != nil {
			ui.Error("Failed to remove binary: " + err.Error())
			return err
		}
		ui.Success("Removed binary: " + destPath)
	} else {
		ui.Warn("Binary not found at: " + destPath)
	}

	// try to remove the dir if empty
	os.Remove(destDir)

	// remove from PATH
	removed, err := removeFromPath(destDir)
	if err != nil {
		ui.Warn("Could not update PATH: " + err.Error())
	} else if removed {
		ui.Success("Removed from PATH")
	}

	fmt.Println()
	ui.Info("Config file kept at: ~/.sharedrop/config.json")
	ui.Dim("Delete it manually if you want a full clean uninstall")
	fmt.Println()
	ui.Success("Uninstall complete!")
	return nil
}

// copyFile copies src to dst, creating/overwriting dst
func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	return err
}

// addToPath adds dir to the user's PATH persistently
// returns (true, nil) if it was added, (false, nil) if already there
func addToPath(dir string) (bool, error) {
	if runtime.GOOS == "windows" {
		return addToPathWindows(dir)
	}
	return addToPathUnix(dir)
}

func removeFromPath(dir string) (bool, error) {
	if runtime.GOOS == "windows" {
		return removeFromPathWindows(dir)
	}
	return removeFromPathUnix(dir)
}

// addToPathWindows uses the registry (user-level, no admin needed)
func addToPathWindows(dir string) (bool, error) {
	// use PowerShell to read/write user PATH from registry
	// [Environment]::GetEnvironmentVariable reads from registry, not just current process
	checkScript := fmt.Sprintf(
		`[System.Environment]::GetEnvironmentVariable("PATH","User")`,
	)
	out, err := runPowerShell(checkScript)
	if err != nil {
		return false, err
	}

	currentPath := strings.TrimSpace(out)
	parts := strings.Split(currentPath, ";")
	for _, p := range parts {
		if strings.EqualFold(strings.TrimSpace(p), dir) {
			return false, nil // already there
		}
	}

	newPath := currentPath
	if newPath != "" && !strings.HasSuffix(newPath, ";") {
		newPath += ";"
	}
	newPath += dir

	setScript := fmt.Sprintf(
		`[System.Environment]::SetEnvironmentVariable("PATH","%s","User")`,
		strings.ReplaceAll(newPath, `"`, `\"`),
	)
	_, err = runPowerShell(setScript)
	return err == nil, err
}

func removeFromPathWindows(dir string) (bool, error) {
	checkScript := `[System.Environment]::GetEnvironmentVariable("PATH","User")`
	out, err := runPowerShell(checkScript)
	if err != nil {
		return false, err
	}

	currentPath := strings.TrimSpace(out)
	parts := strings.Split(currentPath, ";")
	newParts := make([]string, 0, len(parts))
	found := false
	for _, p := range parts {
		if strings.EqualFold(strings.TrimSpace(p), dir) {
			found = true
			continue
		}
		if strings.TrimSpace(p) != "" {
			newParts = append(newParts, p)
		}
	}

	if !found {
		return false, nil
	}

	newPath := strings.Join(newParts, ";")
	setScript := fmt.Sprintf(
		`[System.Environment]::SetEnvironmentVariable("PATH","%s","User")`,
		strings.ReplaceAll(newPath, `"`, `\"`),
	)
	_, err = runPowerShell(setScript)
	return err == nil, err
}

// runPowerShell runs a powershell command and returns stdout
func runPowerShell(script string) (string, error) {
	// use os/exec without importing it at the top — add it to imports
	return execCommand("powershell", "-NoProfile", "-Command", script)
}

// addToPathUnix appends export to ~/.bashrc and ~/.zshrc
func addToPathUnix(dir string) (bool, error) {
	line := fmt.Sprintf(`export PATH="$PATH:%s"`, dir)

	home, _ := os.UserHomeDir()
	added := false

	for _, rc := range []string{".bashrc", ".zshrc", ".profile"} {
		rcPath := filepath.Join(home, rc)
		if _, err := os.Stat(rcPath); err != nil {
			continue // skip if doesn't exist
		}
		content, _ := os.ReadFile(rcPath)
		if strings.Contains(string(content), dir) {
			continue // already there
		}
		f, err := os.OpenFile(rcPath, os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			continue
		}
		fmt.Fprintf(f, "\n# Added by ShareDrop installer\n%s\n", line)
		f.Close()
		added = true
	}

	return added, nil
}

func removeFromPathUnix(dir string) (bool, error) {
	home, _ := os.UserHomeDir()
	removed := false

	for _, rc := range []string{".bashrc", ".zshrc", ".profile"} {
		rcPath := filepath.Join(home, rc)
		content, err := os.ReadFile(rcPath)
		if err != nil {
			continue
		}

		lines := strings.Split(string(content), "\n")
		newLines := make([]string, 0, len(lines))
		skip := false
		for _, line := range lines {
			if strings.Contains(line, "Added by ShareDrop installer") {
				skip = true
				removed = true
				continue
			}
			if skip && strings.Contains(line, dir) {
				skip = false
				continue
			}
			skip = false
			newLines = append(newLines, line)
		}

		os.WriteFile(rcPath, []byte(strings.Join(newLines, "\n")), 0644)
	}

	return removed, nil
}

// isInstalledInPath checks if sharedrop is already reachable via PATH
func isInstalledInPath() bool {
	_, err := execCommand("sharedrop", "version")
	return err == nil
}
