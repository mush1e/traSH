package config

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/mush1e/traSH/utils"
)

type Config struct {
	Prompt       string
	PromptColor  string
	PromptSymbol string
}

var conf *Config
var once sync.Once

var defaultConfig = &Config{
	Prompt:       "🗑️ traSH",
	PromptColor:  "yellow",
	PromptSymbol: " $_",
}

func loadConfig() *Config {
	home, err := os.UserHomeDir()
	if err != nil {
		return defaultConfig
	}

	filePath := filepath.Join(home, ".trashrc")
	trashRC := ParseTrashRC(filePath)

	if trashRC == nil {
		return defaultConfig
	}

	return &Config{
		Prompt:       utils.Coalesce(trashRC["prompt"], defaultConfig.Prompt),
		PromptColor:  utils.Coalesce(trashRC["color"], defaultConfig.PromptColor),
		PromptSymbol: utils.Coalesce(trashRC["symbol"], defaultConfig.PromptSymbol),
	}

}

func GetConfig() *Config {
	once.Do(func() {
		conf = loadConfig()
	})
	return conf
}

func ParseTrashRC(filepath string) map[string]string {
	file, err := os.Open(filepath)

	if err != nil {
		return nil
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)
	trashRC := make(map[string]string)
	lineNum := 0
	for scanner.Scan() {
		lineNum++
		line := scanner.Text()
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		pair := strings.SplitN(line, "=", 3)
		if len(pair) == 2 {
			trashRC[strings.TrimSpace(pair[0])] = strings.TrimSpace(pair[1])
		} else {
			errStr := fmt.Sprintf("invalid syntax in .trashrc file in line : %d\nSWITCHING TO DEFAULT CONFIG\n", lineNum-1)
			fmt.Println(utils.Colorize(errStr, "red"))
			continue
		}
	}

	return trashRC
}
