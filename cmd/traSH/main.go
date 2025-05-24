package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/mush1e/traSH/internal/command"
	"github.com/mush1e/traSH/internal/io"
)

func main() {
	go func() {
		io.WriteHeader(os.Stdout)
		io.WritePrompt(os.Stdout)
		command.ParseCommand(io.ReadUserInput(os.Stdin))
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	sig := <-quit
	log.Printf("Shutdown signal received (%v), initiating graceful shutdown...\n", sig)
	log.Println("traSh has been killed (rightfully so)... Thanks for visiting :)")
}
