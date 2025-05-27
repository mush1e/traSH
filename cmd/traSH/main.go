package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/mush1e/traSH/internal/command"
	"github.com/mush1e/traSH/internal/io"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	io.WriteHeader(os.Stdout)
	userExit := make(chan struct{})

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				io.WritePrompt(os.Stdout)
				cmd := command.ParseCommand(io.ReadUserInput(os.Stdin))

				if cmd.GetCommand() == "exit" {
					close(userExit)
					return
				}

				if err := command.HandleCommand(cmd); err != nil {
					fmt.Printf("error executing command - %v\n", err)
				}
			}
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	select {
	case sig := <-quit:
		fmt.Println()
		fmt.Println()
		log.Printf("Shutdown signal received (%v), initiating graceful shutdown...\n", sig)
		cancel()
	case <-userExit:
		fmt.Println()
		fmt.Println()
		log.Println("Exit command received, initiating graceful shutdown...")
		cancel()
	}

	log.Println("traSh has been killed (rightfully so)... Thanks for visiting :)")
}
