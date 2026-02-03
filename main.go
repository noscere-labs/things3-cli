package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/yourusername/things3-cli/cmd"
)

// rootCmd is the main command that all subcommands attach to
var rootCmd = &cobra.Command{
	Use:   "things",
	Short: "A command-line interface for Things 3",
	Long: `things - A CLI for programmatic interaction with Things using its URL scheme.

This tool lets you add and update to-dos or projects, open lists, and send
JSON payloads to Things from the command line.

For more information, visit: https://culturedcode.com/things/`,
}

// helpCmd provides help information
var helpCmd = &cobra.Command{
	Use:   "help [command]",
	Short: "Help about any command",
	Run: func(cmd *cobra.Command, args []string) {
		rootCmd.Help()
	},
}

func init() {
	for _, c := range cmd.GetCommands() {
		rootCmd.AddCommand(c)
	}
	rootCmd.AddCommand(helpCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
