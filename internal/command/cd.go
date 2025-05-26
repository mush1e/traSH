package command

import "os"

func HandleCD(cmd *Command) error {
	if len(cmd.args) == 0 {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return err
		}
		if err := os.Chdir(homeDir); err != nil {
			return err
		}
	} else {
		if err := os.Chdir(cmd.args[0]); err != nil {
			return err
		}
	}

	return nil
}
