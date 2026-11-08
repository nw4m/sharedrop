package ui

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"syscall"

	"github.com/charmbracelet/lipgloss"
	"golang.org/x/term"
)

// color palette
var (
	styleBanner = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#a78bfa")).
			Bold(true)

	styleTitle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#7c6af7")).
			Bold(true)

	styleSuccess = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#4ade80")).
			Bold(true)

	styleError = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#f87171")).
			Bold(true)

	styleWarn = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#fb923c"))

	styleInfo = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#60a5fa"))

	styleDim = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#6b7280"))

	styleHighlight = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#a78bfa")).
			Bold(true)

	styleBox = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#2e2e4a")).
			Padding(1, 2)

	styleURLBox = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#7c6af7")).
			Padding(0, 2).
			Foreground(lipgloss.Color("#a78bfa"))
)

// PrintBanner prints the ASCII art banner
func PrintBanner() {
	banner := `
 ____  _                    ____                
/ ___|| |__   __ _ _ __ ___| __ ) _ __ ___  _ __
\___ \| '_ \ / _' | '__/ _ \  _ \| '__/ _ \| '_ \
 ___) | | | | (_| | | |  __/ |_) | | | (_) | |_) |
|____/|_| |_|\__,_|_|  \___|____/|_|  \___/| .__/ 
                                            |_|    `

	fmt.Println(styleBanner.Render(banner))
	fmt.Println(styleDim.Render("  v1.0.0 · Local Network File Sharing · Built with ❤️ in Go"))
	fmt.Println()
}

// PrintSeparator prints a styled separator line
func PrintSeparator() {
	fmt.Println(styleDim.Render(strings.Repeat("─", 56)))
}

// Success prints a success message
func Success(msg string) {
	fmt.Println(styleSuccess.Render("  ✓ " + msg))
}

// Error prints an error message
func Error(msg string) {
	fmt.Println(styleError.Render("  ✗ " + msg))
}

// Warn prints a warning
func Warn(msg string) {
	fmt.Println(styleWarn.Render("  ⚠ " + msg))
}

// Info prints an info message
func Info(msg string) {
	fmt.Println(styleInfo.Render("  → " + msg))
}

// Dim prints dimmed text
func Dim(msg string) {
	fmt.Println(styleDim.Render("    " + msg))
}

// Title prints a section title
func Title(msg string) {
	fmt.Println()
	fmt.Println(styleTitle.Render("  " + msg))
	fmt.Println(styleDim.Render("  " + strings.Repeat("─", len(msg)+2)))
}

// Prompt asks a yes/no question, returns true if yes
func Prompt(question string) bool {
	fmt.Print(styleHighlight.Render("  ? ") + question + styleDim.Render(" [y/N] ") + ": ")
	reader := bufio.NewReader(os.Stdin)
	answer, _ := reader.ReadString('\n')
	answer = strings.TrimSpace(strings.ToLower(answer))
	return answer == "y" || answer == "yes"
}

// PromptDefault asks y/n with a default of yes
func PromptDefault(question string, defaultYes bool) bool {
	hint := "[Y/n]"
	if !defaultYes {
		hint = "[y/N]"
	}
	fmt.Print(styleHighlight.Render("  ? ") + question + " " + styleDim.Render(hint) + ": ")
	reader := bufio.NewReader(os.Stdin)
	answer, _ := reader.ReadString('\n')
	answer = strings.TrimSpace(strings.ToLower(answer))
	if answer == "" {
		return defaultYes
	}
	return answer == "y" || answer == "yes"
}

// Input reads a text input from the user
func Input(prompt string, defaultVal string) string {
	if defaultVal != "" {
		fmt.Print(styleHighlight.Render("  > ") + prompt + styleDim.Render(" ["+defaultVal+"]") + ": ")
	} else {
		fmt.Print(styleHighlight.Render("  > ") + prompt + ": ")
	}
	reader := bufio.NewReader(os.Stdin)
	val, _ := reader.ReadString('\n')
	val = strings.TrimSpace(val)
	if val == "" {
		return defaultVal
	}
	return val
}

// PasswordInput reads a password without echoing it
func PasswordInput(prompt string) string {
	fmt.Print(styleHighlight.Render("  🔑 ") + prompt + ": ")
	pass, err := term.ReadPassword(int(syscall.Stdin))
	fmt.Println() // newline after hidden input
	if err != nil {
		// fallback to regular input if terminal doesn't support it
		reader := bufio.NewReader(os.Stdin)
		val, _ := reader.ReadString('\n')
		return strings.TrimSpace(val)
	}
	return string(pass)
}

// PrintShareInfo prints the sharing box with URLs and QR code hint
func PrintShareInfo(urls []string, folder string, passwordSet bool) {
	fmt.Println()
	fmt.Println(styleSuccess.Render("  ✓ ShareDrop is running!"))
	fmt.Println()

	// folder info
	fmt.Println(styleInfo.Render("  📁 Sharing:") + " " + styleHighlight.Render(folder))
	fmt.Println()

	// URLs
	fmt.Println(styleInfo.Render("  🌐 Access URLs:"))
	for _, url := range urls {
		fmt.Println("     " + styleURLBox.Render(url))
	}

	if passwordSet {
		fmt.Println()
		fmt.Println(styleWarn.Render("  🔒 Password protection is enabled"))
	}

	fmt.Println()
	PrintSeparator()
	fmt.Println(styleDim.Render("  Press Ctrl+C to stop sharing"))
	PrintSeparator()
	fmt.Println()
}

// PrintSetupHeader prints the first-run setup header
func PrintSetupHeader() {
	fmt.Println(styleBox.Render(
		styleTitle.Render("First Run Setup") + "\n\n" +
			styleDim.Render("Let's get ShareDrop configured.\nThis only runs once — settings are saved locally."),
	))
	fmt.Println()
}

// PrintStep prints a numbered setup step
func PrintStep(n int, msg string) {
	num := styleDim.Render(fmt.Sprintf("  [%d]", n))
	fmt.Println(num + " " + styleTitle.Render(msg))
}

// PrintConfig prints a summary of current config
func PrintConfig(port int, allowUpload, allowDelete, showHidden bool, maxMB int, passwordSet bool) {
	fmt.Println()
	fmt.Println(styleBox.Render(
		styleTitle.Render("Configuration Summary") + "\n\n" +
			configRow("Port", fmt.Sprintf("%d", port)) + "\n" +
			configRow("Allow Upload", boolStr(allowUpload)) + "\n" +
			configRow("Allow Delete", boolStr(allowDelete)) + "\n" +
			configRow("Show Hidden Files", boolStr(showHidden)) + "\n" +
			configRow("Max Upload Size", fmt.Sprintf("%d MB", maxMB)) + "\n" +
			configRow("Password Protected", boolStr(passwordSet)),
	))
	fmt.Println()
}

func configRow(key, val string) string {
	padded := key + strings.Repeat(" ", 22-len(key))
	return styleDim.Render("  "+padded) + styleHighlight.Render(val)
}

func boolStr(b bool) string {
	if b {
		return "Yes ✓"
	}
	return "No"
}
