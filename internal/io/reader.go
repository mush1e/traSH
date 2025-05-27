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

const (
	CursorUp    = "\033[A"
	CursorDown  = "\033[B"
	CursorRight = "\033[C"
	CursorLeft  = "\033[D"
	CursorHome  = "\033[H"
	ClearLine   = "\033[2K"

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
	content      []rune
	cursor       int
	prompt       string
	promptLen    int
	suggestions  []string
	suggestIndex int
	lastPrefix   string
}

func NewInputBuffer(prompt string) *InputBuffer {
	promptLen := len(stripAnsiCodes(prompt))
	return &InputBuffer{
		content:   make([]rune, 0),
		cursor:    0,
		prompt:    prompt,
		promptLen: promptLen,
	}
}

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
	ib.content = append(ib.content[:ib.cursor], append([]rune{r}, ib.content[ib.cursor:]...)...)
	ib.cursor++
}

func (ib *InputBuffer) deleteBackward() bool {
	if ib.cursor == 0 {
		return false
	}
	ib.content = append(ib.content[:ib.cursor-1], ib.content[ib.cursor:]...)
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
	fmt.Print("\r" + ClearLine)
	fmt.Print(ib.prompt + " ")
	contentStr := string(ib.content)
	words := strings.Fields(contentStr)
	wordStart := 0
	for i, word := range words {
		wordPos := strings.Index(contentStr[wordStart:], word)
		if wordPos != -1 {
			wordStart += wordPos
		}
		if i > 0 {
			fmt.Print(" ")
		}
		color := "reset"
		if i == 0 {
			color = "green"
		} else if strings.HasPrefix(word, "-") {
			color = "yellow"
		} else if strings.Contains(word, "/") || strings.Contains(word, "~") {
			color = "blue"
		}
		fmt.Print(utils.Colorize(word, color))
		wordStart += len(word)
	}
	fmt.Printf("\r\033[%dC", ib.promptLen+1+ib.cursor)
}

func ReadUserInput(prompt string) string {
	fd := int(os.Stdin.Fd())
	oldState, err := term.MakeRaw(fd)
	if err != nil {
		fmt.Println("Failed to enter raw mode, fallback")
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
			handleEscapeSequence(reader, buffer)
		case KeyTab:
			// Capture current word only on first Tab
			currWord := buffer.CurrentWord()
			if len(buffer.suggestions) == 0 {
				// Generate suggestions based on original prefix
				buffer.lastPrefix = currWord
				buffer.suggestIndex = 0
				if buffer.isFirstWord() {
					buffer.suggestions = getCommandSuggestions(currWord)
				} else {
					buffer.suggestions = getFilePathSuggestions(currWord)
				}
			}
			// Cycle through suggestions
			if len(buffer.suggestions) > 0 {
				suggestion := buffer.suggestions[buffer.suggestIndex%len(buffer.suggestions)]
				// Remove the current word (could be original prefix or previous suggestion)
				start := buffer.cursor - len([]rune(currWord))
				buffer.content = append(buffer.content[:start], append([]rune(suggestion), buffer.content[buffer.cursor:]...)...)
				buffer.cursor = start + len([]rune(suggestion))
				buffer.suggestIndex++
			}
		default:
			if unicode.IsPrint(char) && char != 0 {
				buffer.insertRune(char)
				buffer.suggestions = nil
				buffer.suggestIndex = 0
				buffer.lastPrefix = ""
			}
		}
		buffer.render()
	}
	return buffer.getText()
}

func handleEscapeSequence(reader *bufio.Reader, buffer *InputBuffer) bool {
	char2, _, err := reader.ReadRune()
	if err != nil || char2 != '[' {
		return false
	}
	char3, _, err := reader.ReadRune()
	if err != nil {
		return false
	}
	switch char3 {
	case 'A':
		prev := history.GetPrevious()
		if prev != "" {
			buffer.content = []rune(prev)
			buffer.cursor = len(buffer.content)
		}
		return true
	case 'B':
		next := history.GetNext()
		buffer.content = []rune(next)
		buffer.cursor = len(buffer.content)
		return true
	case 'C':
		return buffer.moveCursorRight()
	case 'D':
		return buffer.moveCursorLeft()
	case 'H':
		buffer.moveCursorHome()
		return true
	case 'F':
		buffer.moveCursorEnd()
		return true
	case '3':
		delChar, _, err := reader.ReadRune()
		if err == nil && delChar == '~' {
			return buffer.deleteForward()
		}
		return false
	}
	return false
}

func readBasicInput(prompt string) string {
	fmt.Print(prompt + " ")
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	return strings.TrimSpace(input)
}
