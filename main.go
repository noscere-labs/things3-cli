package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/yourusername/bear-cli/cmd"
)

// rootCmd is the main command that all subcommands attach to
var rootCmd = &cobra.Command{
	Use:   "bear",
	Short: "A command-line interface for Bear note app",
	Long: `bear - A powerful CLI for programmatic interaction with Bear notes

This tool allows you to create, read, update, and manage notes in Bear
from the command line. Perfect for integration with scripts and tools like Claude Code.

For more information, visit: https://github.com/yourusername/bear-cli

Examples:
  bear create --title "My Note" --content "Note content" --tags "work"
  bear read --id "7E4B681B-..."
  bear update --id "7E4B681B-..." --content "Updated content"
  bear list --tag "work"
  bear config set-token --token "YOUR_API_TOKEN"`,
	// Silently ignore if no command is provided (Cobra default behavior)
	// User will see help text when they run 'bear' without arguments
}

// versionCmd displays the version of bear-cli
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show bear CLI version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("bear version 1.0.0")
	},
}

// helpCmd provides help information
var helpCmd = &cobra.Command{
	Use:   "help [command]",
	Short: "Help about any command",
	Long: `Help provides help for any command in the application.
Simply type bear help [path to command] for full details.`,
	Run: func(cmd *cobra.Command, args []string) {
		// This is handled automatically by Cobra
		rootCmd.Help()
	},
}

// init initializes the command structure
// This function sets up all subcommands and configures global options
func init() {
	// Disable automatic help flag to manage it ourselves if needed
	// rootCmd.DisableFlagParsing = false

	// Add all subcommands from cmd package
	for _, c := range cmd.GetCommands() {
		rootCmd.AddCommand(c)
	}

	// Add version command
	rootCmd.AddCommand(versionCmd)

	// Configure output behavior
	// Disable sorting of commands in help (we'll use our own order)
	// rootCmd.SortCommandsByString = true

	// Handle completion for bash/zsh (optional)
	// This would allow tab completion if generated properly
}

// main is the entry point for the bear CLI
// It parses command-line arguments and executes the appropriate command
func main() {
	// Execute the root command
	// This will parse args, run the appropriate subcommand, or show help/errors
	if err := rootCmd.Execute(); err != nil {
		// Cobra handles most error printing automatically
		// This catch-all is for unexpected errors
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
