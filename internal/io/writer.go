package io

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/mush1e/traSH/config"
	"github.com/mush1e/traSH/utils"
)

var conf = config.GetConfig()

func WriteHeader(w io.Writer) {
	header := `████████╗██████╗  █████╗ ███████╗██╗  ██╗
╚══██╔══╝██╔══██╗██╔══██╗██╔════╝██║  ██║
   ██║   ██████╔╝███████║███████╗███████║
   ██║   ██╔══██╗██╔══██║╚════██║██╔══██║
   ██║   ██║  ██║██║  ██║███████║██║  ██║
   ╚═╝   ╚═╝  ╚═╝╚═╝  ╚═╝╚══════╝╚═╝  ╚═`
	subHeader := "it's in the name, I have no idea why you're using this shell.... enjoy your stay though :)"
	footer := "- Mush1e"
	fmt.Fprintf(w, "\n\n%v\n\n\n", header)
	fmt.Fprintf(w, "%v\n\n", subHeader)
	fmt.Fprintf(w, "\t\t\t\t\t%v\n\n\n", footer)

}

func WritePrompt(w io.Writer) {
	prompt := conf.Prompt + conf.PromptSymbol
	currPath, err := os.Getwd()
	if err != nil {
		fmt.Fprintf(w, "error getting current path - %q\n", err)
	}
	coloredPrompt := utils.Colorize(prompt, conf.PromptColor)
	fmt.Fprintf(w, "%v | %v ", currPath, coloredPrompt)
}

func BuildPrompt() string {
	conf := config.GetConfig()
	cwd, err := os.Getwd()
	if err != nil {
		cwd = "~"
	}

	// Shorten long paths
	if len(cwd) > 50 {
		parts := strings.Split(cwd, "/")
		if len(parts) > 3 {
			cwd = ".../" + strings.Join(parts[len(parts)-2:], "/")
		}
	}

	return utils.Colorize(conf.Prompt+" "+cwd+conf.PromptSymbol, conf.PromptColor)
}
