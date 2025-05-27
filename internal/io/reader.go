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

var history = NewHistory()

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

	// Print prompt with space
	fmt.Print(ib.prompt + " ")

	// Build the display text and track cursor position
	var displayText strings.Builder
	visualCursorPos := ib.promptLen + 1 // Start after prompt + space

	contentStr := string(ib.content)

	// If we have content, colorize it
	if len(ib.content) > 0 {
		words := strings.Fields(contentStr)
		wordStart := 0

		for i, word := range words {
			// Find where this word starts in the original content
			wordPos := strings.Index(contentStr[wordStart:], word)
			if wordPos != -1 {
				wordStart += wordPos
			}

			// Add spaces before word (except first)
			if i > 0 {
				displayText.WriteString(" ")
				if ib.cursor > wordStart {
					visualCursorPos++
				}
			}

			// Determine color
			color := "reset"
			if i == 0 {
				color = "green" // Command
			} else if strings.HasPrefix(word, "-") {
				color = "yellow" // Flags
			} else if strings.Contains(word, "/") || strings.Contains(word, "~") {
				color = "blue" // Paths
			}

			// Add colored word
			displayText.WriteString(utils.Colorize(word, color))

			// Update cursor position if cursor is within or after this word
			if ib.cursor > wordStart && ib.cursor <= wordStart+len(word) {
				visualCursorPos += ib.cursor - wordStart
			} else if ib.cursor > wordStart+len(word) {
				visualCursorPos += len(word)
			}

			wordStart += len(word)
		}
	}

	// Print the colorized content
	fmt.Print(displayText.String())

	// Position cursor at the correct location
	// Simple approach: calculate visual position based on logical cursor
	cursorCol := ib.promptLen + 1 + ib.cursor
	fmt.Printf("\r\033[%dC", cursorCol)
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

	buffer.render()

	for {
		char, _, err := reader.ReadRune()
		if err != nil {
			break
		}

		switch char {
		case KeyEnter:
			fmt.Print("\r\n")
			history.Add(buffer.getText())
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
		prev := history.GetPrevious()
		if prev != "" {
			buffer.content = []rune(prev)
			buffer.cursor = len(buffer.content)
		}

		return true
	case 'B': // Down arrow
		next := history.GetNext()
		buffer.content = []rune(next)
		buffer.cursor = len(buffer.content)
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
