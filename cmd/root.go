package cmd

import (
	"fmt"
	"os"
	"strconv"

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

func init() {
	rootCmd.PersistentFlags().IntP("port", "p", 0, "Port to listen on (overrides config)")
	rootCmd.PersistentFlags().StringP("dir", "d", "", "Directory to share (current dir if not set)")
	rootCmd.PersistentFlags().BoolP("no-password", "n", false, "Skip password even if one is set")

	rootCmd.AddCommand(setupCmd)
	rootCmd.AddCommand(configCmd)
	rootCmd.AddCommand(passwdCmd)
	rootCmd.AddCommand(versionCmd)
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

	// first run check
	dir, _ := config.ConfigDir()
	if _, statErr := os.Stat(dir + "/config.json"); os.IsNotExist(statErr) {
		ui.Info("Looks like your first time! Let's do a quick setup.")
		fmt.Println()
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
