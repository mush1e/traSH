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

// ANSI escape codes for terminal control
const (
	// Cursor movement
	CursorUp    = "\033[A"
	CursorDown  = "\033[B"
	CursorRight = "\033[C"
	CursorLeft  = "\033[D"
	CursorHome  = "\033[H"

	// Line operations
	ClearLine     = "\033[2K"
	ClearToEnd    = "\033[0K"
	SaveCursor    = "\033[s"
	RestoreCursor = "\033[u"

	// Keys
	KeyEscape    = 27
	KeyBackspace = 127
	KeyDelete    = 126
	KeyTab       = 9
	KeyEnter     = 13
	KeyCtrlC     = 3
	KeyCtrlD     = 4
)

type InputBuffer struct {
	content   []rune
	cursor    int
	prompt    string
	promptLen int
}

func NewInputBuffer(prompt string) *InputBuffer {
	// Calculate actual prompt length (excluding ANSI codes)
	promptLen := len(stripAnsiCodes(prompt))

	return &InputBuffer{
		content:   make([]rune, 0),
		cursor:    0,
		prompt:    prompt,
		promptLen: promptLen,
	}
}

// Strip ANSI escape codes to get actual display length
func stripAnsiCodes(s string) string {
	var result strings.Builder
	inEscape := false

	for _, r := range s {
		if r == '\033' {
			inEscape = true
			continue
		}
		if inEscape {
			if r == 'm' || r == 'K' || r == 'J' {
				inEscape = false
			}
			continue
		}
		result.WriteRune(r)
	}
	return result.String()
}

func (ib *InputBuffer) insertRune(r rune) {
	if ib.cursor == len(ib.content) {
		// Insert at end
		ib.content = append(ib.content, r)
	} else {
		// Insert in middle
		ib.content = append(ib.content[:ib.cursor+1], ib.content[ib.cursor:]...)
		ib.content[ib.cursor] = r
	}
	ib.cursor++
}

func (ib *InputBuffer) deleteBackward() bool {
	if ib.cursor == 0 {
		return false
	}

	if ib.cursor == len(ib.content) {
		// Delete from end
		ib.content = ib.content[:len(ib.content)-1]
	} else {
		// Delete from middle
		ib.content = append(ib.content[:ib.cursor-1], ib.content[ib.cursor:]...)
	}
	ib.cursor--
	return true
}

func (ib *InputBuffer) deleteForward() bool {
	if ib.cursor >= len(ib.content) {
		return false
	}

	ib.content = append(ib.content[:ib.cursor], ib.content[ib.cursor+1:]...)
	return true
}

func (ib *InputBuffer) moveCursorLeft() bool {
	if ib.cursor > 0 {
		ib.cursor--
		return true
	}
	return false
}

func (ib *InputBuffer) moveCursorRight() bool {
	if ib.cursor < len(ib.content) {
		ib.cursor++
		return true
	}
	return false
}

func (ib *InputBuffer) moveCursorHome() {
	ib.cursor = 0
}

func (ib *InputBuffer) moveCursorEnd() {
	ib.cursor = len(ib.content)
}

func (ib *InputBuffer) getText() string {
	return string(ib.content)
}

func (ib *InputBuffer) render() {
	// Move to beginning of line and clear it
	fmt.Print("\r" + ClearLine)

	// Print prompt
	fmt.Print(ib.prompt + " ")

	// Colorize and print content
	words := strings.Fields(string(ib.content))
	currentPos := 0

	for i, word := range words {
		// Add spaces between words (except first)
		if i > 0 {
			fmt.Print(" ")
			currentPos++
		}

		// Colorize based on position
		color := "reset"
		if i == 0 {
			color = "green" // Command
		} else if strings.HasPrefix(word, "-") {
			color = "yellow" // Flags
		} else if strings.Contains(word, "/") || strings.Contains(word, "~") {
			color = "blue" // Paths
		}

		fmt.Print(utils.Colorize(word, color))
		currentPos += len(word)
	}

	// Position cursor correctly
	actualCursorPos := ib.promptLen + 1 + ib.cursor
	if len(ib.content) > 0 {
		// Account for spaces between words
		spaceCount := len(strings.Fields(string(ib.content[:ib.cursor]))) - 1
		if spaceCount > 0 {
			actualCursorPos += spaceCount
		}
	}

	// Move cursor to correct position
	fmt.Print("\r")
	for i := 0; i < actualCursorPos; i++ {
		fmt.Print(CursorRight)
	}
}

func ReadUserInput(prompt string) string {
	return readEnhancedInput(prompt)
}

func readEnhancedInput(prompt string) string {
	fd := int(os.Stdin.Fd())
	oldState, err := term.MakeRaw(fd)
	if err != nil {
		fmt.Println("Failed to enter raw mode, falling back to basic input")
		return readBasicInput(prompt)
	}
	defer term.Restore(fd, oldState)

	buffer := NewInputBuffer(prompt)
	reader := bufio.NewReader(os.Stdin)

	// Initial render
	buffer.render()

	for {
		char, _, err := reader.ReadRune()
		if err != nil {
			break
		}

		switch char {
		case KeyEnter:
			fmt.Print("\n")
			return buffer.getText()

		case KeyCtrlC:
			fmt.Print("^C\n")
			return ""

		case KeyCtrlD:
			if len(buffer.content) == 0 {
				fmt.Print("\n")
				return "exit"
			}
			buffer.deleteForward()

		case KeyBackspace:
			buffer.deleteBackward()

		case KeyEscape:
			// Handle escape sequences (arrow keys, etc.)
			if handleEscapeSequence(reader, buffer) {
				// Escape sequence handled
			} else {
				// Regular escape key pressed
				continue
			}

		case KeyTab:
			// TODO: Add tab completion here
			buffer.insertRune(' ')
			buffer.insertRune(' ')

		default:
			if unicode.IsPrint(char) && char != 0 {
				buffer.insertRune(char)
			}
		}

		buffer.render()
	}

	return buffer.getText()
}

func handleEscapeSequence(reader *bufio.Reader, buffer *InputBuffer) bool {
	// Read the next character
	char2, _, err := reader.ReadRune()
	if err != nil {
		return false
	}

	if char2 != '[' {
		return false
	}

	// Read the final character
	char3, _, err := reader.ReadRune()
	if err != nil {
		return false
	}

	switch char3 {
	case 'A': // Up arrow
		// TODO: Add history support
		return true
	case 'B': // Down arrow
		// TODO: Add history support
		return true
	case 'C': // Right arrow
		return buffer.moveCursorRight()
	case 'D': // Left arrow
		return buffer.moveCursorLeft()
	case 'H': // Home
		buffer.moveCursorHome()
		return true
	case 'F': // End
		buffer.moveCursorEnd()
		return true
	case '3': // Delete key (followed by ~)
		delChar, _, err := reader.ReadRune()
		if err == nil && delChar == '~' {
			return buffer.deleteForward()
		}
		return false
	}

	return false
}

// Fallback for when raw mode fails
func readBasicInput(prompt string) string {
	fmt.Print(prompt + " ")
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	return strings.TrimSpace(input)
}
