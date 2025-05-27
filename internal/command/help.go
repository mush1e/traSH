package command

import "fmt"

func printHelp() {
	help := `traSH - Built-in Commands:

Built-ins:
  cd <dir>     Change directory
  help/?       Show this help
  exit         Exit the shell

Features:
  • Arrow keys for cursor movement
  • Quoted arguments: "hello world"
  • Command options: -n, --verbose
  • External command support
  • Escape sequences in quotes: "line1\nline2"

Examples:
  cd "My Documents"
  echo -n "Hello, World!"
  grep "error message" /var/log/app.log
  ls -la
  mkdir "New Folder"
`
	fmt.Print(help)
}
