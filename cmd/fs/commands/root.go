package commands

import "github.com/spf13/cobra"

const (
	// CLIName is the name of the CLI
	CLIName = "fs"
)

// NewCommand returns a new instance of an argo command
func NewCommand() *cobra.Command {
	var command = &cobra.Command{
		Use:   CLIName,
		Short: "fs is the command line interface for distributed file-system",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			cmd.HelpFunc()(cmd, args)
		},
	}

	command.AddCommand(NewServerCommand())
	return command
}
