package commands

import (
	"context"
	"fmt"

	"github.com/mohanpothukuchi/dfs/server/apiserver"
	"github.com/spf13/cobra"
)

//NewServerCommand is used to start the master server instance for the file-system.
func NewServerCommand() *cobra.Command {
	var (
		port        string
		servAddress string
	)

	var command = cobra.Command{
		Use:   "server",
		Short: "Start the dfs Server",
		Run: func(cmd *cobra.Command, args []string) {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			server, err := apiserver.NewMasterServer(port, servAddress)
			if err != nil {
				fmt.Println(err)
			}
			server.Run(ctx)
		},
	}

	command.Flags().StringVarP(&port, "port", "p", "2746", "Port to listen on")
	command.Flags().StringVarP(&servAddress, "address", "s", "localhost", "Default to localhost")
	return &command
}
