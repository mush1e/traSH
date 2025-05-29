package io

import (
	"os"
	"path/filepath"
	"strings"
)

func (ib *InputBuffer) CurrentWord() string {
	inputSoFar := string(ib.content[:ib.cursor])

	// Find the start of the current argument, considering quotes and escapes
	wordStart := ib.findWordStart(inputSoFar)
	return inputSoFar[wordStart:]
}

func (ib *InputBuffer) findWordStart(input string) int {
	if len(input) == 0 {
		return 0
	}

	inQuotes := false
	var quoteChar rune
	escaped := false
	wordStart := 0

	for i, r := range input {
		if escaped {
			escaped = false
			continue
		}

		if r == '\\' {
			escaped = true
			continue
		}

		if !inQuotes && (r == '"' || r == '\'') {
			inQuotes = true
			quoteChar = r
			wordStart = i
			continue
		}

		if inQuotes && r == quoteChar {
			inQuotes = false
			continue
		}

		if !inQuotes && r == ' ' {
			wordStart = i + 1
		}
	}

	return wordStart
}

func (ib *InputBuffer) isFirstWord() bool {
	input := string(ib.content[:ib.cursor])
	wordStart := ib.findWordStart(input)

	// Check if there are any non-whitespace characters before wordStart
	for i := 0; i < wordStart; i++ {
		if input[i] != ' ' {
			return false
		}
	}
	return true
}

func getCommandSuggestions(prefix string) []string {
	suggestions := make([]string, 0)
	seen := make(map[string]bool)
	pathEnv := os.Getenv("PATH")
	dirs := strings.Split(pathEnv, ":")

	for _, dir := range dirs {
		files, err := os.ReadDir(dir)
		if err != nil {
			continue
		}
		for _, file := range files {
			name := file.Name()
			if strings.HasPrefix(name, prefix) && !seen[name] {
				fullPath := filepath.Join(dir, name)

				if info, err := os.Stat(fullPath); err == nil && info.Mode().IsRegular() && info.Mode().Perm()&0111 != 0 {
					suggestions = append(suggestions, name)
					seen[name] = true
				}
			}
		}
	}
	return suggestions
}

// Helper function to escape spaces and special characters for shell
func escapeForShell(s string) string {
	// Characters that need escaping in shell
	specialChars := []string{" ", "(", ")", "[", "]", "{", "}", "'", "\"", "\\", "$", "`", "!", "&", ";", "<", ">", "|", "*", "?"}

	result := s
	for _, char := range specialChars {
		if strings.Contains(result, char) {
			// Use backslash escaping for most characters
			result = strings.ReplaceAll(result, char, "\\"+char)
		}
	}
	return result
}

// Helper function to check if a string needs escaping
func needsEscaping(s string) bool {
	specialChars := []string{" ", "(", ")", "[", "]", "{", "}", "'", "\"", "\\", "$", "`", "!", "&", ";", "<", ">", "|", "*", "?"}
	for _, char := range specialChars {
		if strings.Contains(s, char) {
			return true
		}
	}
	return false
}

// Helper function to unescape a string for matching
func unescapeForMatching(s string) string {
	// Remove backslashes before special characters
	result := s
	result = strings.ReplaceAll(result, "\\ ", " ")
	result = strings.ReplaceAll(result, "\\(", "(")
	result = strings.ReplaceAll(result, "\\)", ")")
	result = strings.ReplaceAll(result, "\\[", "[")
	result = strings.ReplaceAll(result, "\\]", "]")
	result = strings.ReplaceAll(result, "\\{", "{")
	result = strings.ReplaceAll(result, "\\}", "}")
	result = strings.ReplaceAll(result, "\\'", "'")
	result = strings.ReplaceAll(result, "\\\"", "\"")
	result = strings.ReplaceAll(result, "\\\\", "\\")
	result = strings.ReplaceAll(result, "\\$", "$")
	result = strings.ReplaceAll(result, "\\`", "`")
	result = strings.ReplaceAll(result, "\\!", "!")
	result = strings.ReplaceAll(result, "\\&", "&")
	result = strings.ReplaceAll(result, "\\;", ";")
	result = strings.ReplaceAll(result, "\\<", "<")
	result = strings.ReplaceAll(result, "\\>", ">")
	result = strings.ReplaceAll(result, "\\|", "|")
	result = strings.ReplaceAll(result, "\\*", "*")
	result = strings.ReplaceAll(result, "\\?", "?")
	return result
}

func getFilePathSuggestions(prefix string) []string {
	homePath, err := os.UserHomeDir()
	if err != nil {
		return nil
	}

	originalPrefix := prefix

	// Clean the prefix - remove quotes and handle escapes
	cleanPrefix := cleanPath(prefix)

	// Handle tilde expansion
	if len(cleanPrefix) > 0 && cleanPrefix[0] == '~' {
		cleanPrefix = strings.Replace(cleanPrefix, "~", homePath, 1)
	}

	var dir string
	var filenamePrefix string

	if strings.HasSuffix(cleanPrefix, string(filepath.Separator)) {
		dir = cleanPrefix
		filenamePrefix = ""
	} else if cleanPrefix == "" {
		// Handle empty prefix specially
		dir, err = os.Getwd()
		if err != nil {
			dir = homePath
		}
		filenamePrefix = ""
	} else {
		dir = filepath.Dir(cleanPrefix)
		filenamePrefix = filepath.Base(cleanPrefix)
	}

	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil
	}

	suggestions := make([]string, 0)
	seen := make(map[string]bool)

	// Separate directories and files, prioritize directories for cd
	var dirSuggestions []string
	var fileSuggestions []string

	for _, entry := range entries {
		name := entry.Name()
		if seen[name] {
			continue
		}

		// Skip hidden files unless specifically requested
		if len(filenamePrefix) == 0 && strings.HasPrefix(name, ".") {
			continue
		}

		if filenamePrefix == "" || strings.HasPrefix(name, filenamePrefix) {
			var suggestion string

			// For cd completion, we want relative paths or just names
			if originalPrefix == "" || cleanPrefix == "" {
				// No prefix - just return the name
				suggestion = name
			} else if strings.HasPrefix(originalPrefix, "~") || strings.HasPrefix(cleanPrefix, "~") {
				// Tilde prefix - return relative to home with tilde
				suggestion = filepath.Join("~", strings.TrimPrefix(filepath.Join(dir, name), homePath))
			} else if filepath.IsAbs(cleanPrefix) {
				// Absolute path - return full path
				suggestion = filepath.Join(dir, name)
			} else {
				// Relative path - return relative
				if dir == "." {
					suggestion = name
				} else {
					suggestion = filepath.Join(dir, name)
				}
			}

			if entry.IsDir() {
				suggestion += string(filepath.Separator)
				dirSuggestions = append(dirSuggestions, suggestion)
			} else {
				fileSuggestions = append(fileSuggestions, suggestion)
			}
			seen[name] = true
		}
	}

	// For cd command, prioritize directories
	suggestions = append(suggestions, dirSuggestions...)
	suggestions = append(suggestions, fileSuggestions...)
	return suggestions
}

func cleanPath(path string) string {
	if len(path) == 0 {
		return path
	}

	// Remove surrounding quotes
	if (path[0] == '"' && path[len(path)-1] == '"') ||
		(path[0] == '\'' && path[len(path)-1] == '\'') {
		path = path[1 : len(path)-1]
	}

	// Handle escaped spaces (convert "\ " to " ")
	path = strings.ReplaceAll(path, "\\ ", " ")

	return path
}
