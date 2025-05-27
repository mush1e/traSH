package command

import "fmt"

func printHelp() {
	fmt.Println("\nAvailable commands:")
	fmt.Println("  cd <dir>        change directory")
	fmt.Println("  exit            exit the shell")
	fmt.Println("  help            show this help message")
	fmt.Println()
	fmt.Println("\nYou can also run any system command like 'ls', 'cat', 'git', etc.")
}
