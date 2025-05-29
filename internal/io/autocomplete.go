package io

import (
	"os"
	"path/filepath"
	"strings"
)

func (ib *InputBuffer) CurrentWord() string {
	inputSoFar := string(ib.content[:ib.cursor])
	lastSpace := strings.LastIndex(inputSoFar, " ")
	if lastSpace == -1 {
		return inputSoFar
	}
	return inputSoFar[lastSpace+1:]
}

func (ib *InputBuffer) isFirstWord() bool {
	return strings.LastIndex(string(ib.content[:ib.cursor]), " ") == -1
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

func getFilePathSuggestions(prefix string) []string {
	homePath, err := os.UserHomeDir()
	if err != nil {
		return nil
	}

	originalPrefix := prefix

	// Handle tilde expansion
	if len(prefix) > 0 && prefix[0] == '~' {
		prefix = strings.Replace(prefix, "~", homePath, 1)
	}

	dir := filepath.Dir(prefix)
	filenamePrefix := filepath.Base(prefix)

	if strings.HasSuffix(prefix, string(filepath.Separator)) {
		dir = prefix
		filenamePrefix = ""
	}

	if prefix == "" {
		dir, err = os.Getwd()
		if err != nil {
			dir = homePath
		}
	}

	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil
	}

	suggestions := make([]string, 0)
	seen := make(map[string]bool)

	for _, entry := range entries {
		name := entry.Name()
		if seen[name] {
			continue
		}

		if filenamePrefix == "" || strings.HasPrefix(name, filenamePrefix) {
			var suggestion string

			// For cd completion, we want relative paths or just names
			if originalPrefix == "" {
				// No prefix - just return the name
				suggestion = name
			} else if strings.HasPrefix(originalPrefix, "~") {
				// Tilde prefix - return relative to home with tilde
				suggestion = filepath.Join("~", strings.TrimPrefix(filepath.Join(dir, name), homePath))
			} else if filepath.IsAbs(originalPrefix) {
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
			}

			suggestions = append(suggestions, suggestion)
			seen[name] = true
		}
	}

	return suggestions
}
