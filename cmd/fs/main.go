package main

import (
	"fmt"
	"os"

	"github.com/mohanpothukuchi/dfs/cmd/fs/commands"
)

func main() {
	if err := commands.NewCommand().Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
