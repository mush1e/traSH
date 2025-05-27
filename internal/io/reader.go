package io

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"unicode"

	"github.com/mush1e/traSH/utils"
	"golang.org/x/term"
)

func ReadUserInput(prompt string) string {
	return readInteractiveLine(prompt)
}

func readInteractiveLine(prompt string) string {
	fd := int(os.Stdin.Fd())
	oldState, err := term.MakeRaw(fd)
	if err != nil {
		fmt.Println("failed to enter raw mode")
		return ""
	}
	defer term.Restore(fd, oldState)

	var input []rune
	reader := bufio.NewReader(os.Stdin)
	fmt.Print(prompt)

	for {
		char, _, err := reader.ReadRune()
		if err != nil {
			break
		}

		switch char {
		case '\r', '\n':
			fmt.Print("\r\n") // Ensure cursor moves to start of new line
			return string(input)
		case 127: // backspace
			if len(input) > 0 {
				input = input[:len(input)-1]
			}
		case '\x1b': // Escape sequence (arrow keys)
			next, _, err := reader.ReadRune()
			if err != nil {
				break
			}
			if next == '[' {
				dir, _, err := reader.ReadRune()
				if err != nil {
					break
				}
				if dir == 'A' || dir == 'B' || dir == 'C' || dir == 'D' {
					// Ignore arrow key input
					continue
				}
			}
			// Ignore other escape sequences
			continue
		default:
			if unicode.IsPrint(char) {
				input = append(input, char)
			}
		}

		// Clear entire line and redraw
		fmt.Print("\r\033[2K") // Clear entire line
		fmt.Print(prompt)

		words := strings.Fields(string(input))
		for i, word := range words {
			color := "reset"
			if i == 0 {
				color = "green" // Command
			} else if strings.HasPrefix(word, "-") {
				color = "white" // Flags
			}

			// Add space only between words, not after the last one
			if i > 0 {
				fmt.Print(" ")
			}
			fmt.Print(utils.Colorize(word, color))
		}
	}

	return string(input)
}
