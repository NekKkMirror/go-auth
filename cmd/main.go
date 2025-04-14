package main

import (
	"context"
	"github.com/NekKkMirror/go-auth/internal/app"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	errChan := make(chan error, 1)

	application := app.NewApp()

	go func() {
		log.Println("Starting application...")
		if err := application.Run(); err != nil {
			errChan <- err
		}
	}()

	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

		sig := <-sigChan
		log.Printf("Received signal: %s", sig)
		cancel()
	}()

	select {
	case <-ctx.Done():
		log.Println("Shutting down gracefully...")
	case err := <-errChan:
		log.Printf("Application error: %v", err)
	}

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	application.Shutdown(shutdownCtx)
	log.Println("Application stopped")
}
