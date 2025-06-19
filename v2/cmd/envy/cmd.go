package main

import (
	"github.com/dubbikins/envy/v2/cmd/envy/ls"
	"github.com/spf13/cobra"
)


var cmd = &cobra.Command{
		Use:  "envy",
		Short: "CLI tool to discover and show environment variables in Go code",
		Long:  "Shows the usage of a struct tag for a go module or package",
}

func main() {
	ls.Extends(cmd)
	// Add any other commands or flags here if needed
	cmd.Execute()
}


