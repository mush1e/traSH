package io

import (
	"fmt"
	"io"

	"github.com/mush1e/traSH/config"
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
	fmt.Fprintf(w, "\t\t\t\t\t%v\n", footer)

}

func WritePrompt(w io.Writer) {
	prompt := conf.Prompt + "> "
	fmt.Fprintf(w, "%v ", prompt)
}
