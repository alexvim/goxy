package main

import (
	"context"
	"goxy/internal/config"
	"goxy/internal/server"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

func main() {
	log.Println("main: start application")

	cfg, err := config.ReadFromArgs(os.Args[1:])
	if err != nil {
		log.Printf("failed to read config err=%s", err)
		return
	}

	ctx, cancel := context.WithCancel(context.Background())
	wg := &sync.WaitGroup{}

	wg.Add(1)
	go func() {
		defer wg.Done()
		monitorSyscall(ctx, cancel)
	}()

	server.Run(ctx, cfg)

	cancel()

	wg.Wait()

	log.Println("main: close application")
}

func monitorSyscall(ctx context.Context, doClose func()) {
	defer signal.Reset()

	done := make(chan os.Signal, 1)

	signal.Notify(done, syscall.SIGINT, syscall.SIGTERM)

	select {
	case <-done:
	case <-ctx.Done():
	}

	doClose()
}
