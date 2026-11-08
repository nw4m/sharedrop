package cmd

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	qrcode "github.com/skip2/go-qrcode"
	"github.com/nw4m/sharedrop/internal/config"
	"github.com/nw4m/sharedrop/internal/discovery"
	"github.com/nw4m/sharedrop/internal/server"
	"github.com/nw4m/sharedrop/internal/ui"
)

func startServer(cfg *config.Config, shareDir string) error {
	// build the HTTP handler
	srv := server.New(cfg, shareDir)

	addr := fmt.Sprintf(":%d", cfg.DefaultPort)

	// check if port is available
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		ui.Error(fmt.Sprintf("Port %d is already in use", cfg.DefaultPort))
		ui.Dim("Try a different port with: sharedrop -p 8081")
		return err
	}
	ln.Close()

	// get all URLs
	urls := discovery.GetAllShareURLs(cfg.DefaultPort)

	// show share info box
	ui.PrintShareInfo(urls, shareDir, cfg.PasswordSet)

	// print QR code for the primary URL
	primaryURL := discovery.GetShareURL(discovery.GetPrimaryIP(), cfg.DefaultPort)
	printQR(primaryURL)

	httpServer := &http.Server{
		Addr:         addr,
		Handler:      srv,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 0, // no timeout for large file downloads
		IdleTimeout:  120 * time.Second,
	}

	// handle graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	serverErr := make(chan error, 1)
	go func() {
		ui.Dim(fmt.Sprintf("Server listening on %s", addr))
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			serverErr <- err
		}
	}()

	// wait for interrupt or error
	select {
	case err := <-serverErr:
		ui.Error("Server error: " + err.Error())
		return err
	case sig := <-quit:
		_ = sig
		fmt.Println()
		ui.Info("Shutting down ShareDrop...")

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := httpServer.Shutdown(ctx); err != nil {
			ui.Error("Shutdown error: " + err.Error())
			return err
		}
		ui.Success("ShareDrop stopped. Bye!")
		fmt.Println()
	}

	return nil
}

// printQR renders a QR code for the given URL directly in the terminal
// uses Unicode half-block characters so it works on Windows PowerShell too
func printQR(url string) {
	qr, err := qrcode.New(url, qrcode.Medium)
	if err != nil {
		return
	}

	bmp := qr.Bitmap()

	fmt.Println()
	fmt.Println("  Scan to open on your phone:")
	fmt.Println()

	// render two rows at a time using half-block chars
	// ▀ = top filled, bottom empty
	// ▄ = top empty, bottom filled
	// █ = both filled
	// ' ' = both empty
	// add 2-cell quiet zone on each side
	quiet := 2

	// top quiet zone rows
	for i := 0; i < quiet; i++ {
		line := "  "
		for j := 0; j < len(bmp[0])+quiet*2; j++ {
			line += " "
		}
		fmt.Println(line)
	}

	// process bitmap rows in pairs
	for row := 0; row < len(bmp); row += 2 {
		line := "  " + strings.Repeat(" ", quiet) // left quiet zone

		for col := 0; col < len(bmp[row]); col++ {
			top := bmp[row][col]
			// handle odd number of rows
			bot := false
			if row+1 < len(bmp) {
				bot = bmp[row+1][col]
			}

			switch {
			case top && bot:
				line += "█"
			case top && !bot:
				line += "▀"
			case !top && bot:
				line += "▄"
			default:
				line += " "
			}
		}

		line += strings.Repeat(" ", quiet) // right quiet zone
		fmt.Println(line)
	}

	// bottom quiet zone
	for i := 0; i < quiet/2; i++ {
		line := "  "
		for j := 0; j < len(bmp[0])+quiet*2; j++ {
			line += " "
		}
		fmt.Println(line)
	}

	fmt.Println()
	fmt.Printf("  %s\n\n", url)
}
