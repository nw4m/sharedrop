package cmd

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/user/sharedrop/internal/config"
	"github.com/user/sharedrop/internal/discovery"
	"github.com/user/sharedrop/internal/server"
	"github.com/user/sharedrop/internal/ui"
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
