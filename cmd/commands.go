package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/yourusername/things3-cli/pkg/formatter"
	"github.com/yourusername/things3-cli/pkg/things"
	"github.com/yourusername/things3-cli/pkg/util"
)

func addStringParam(cmd *cobra.Command, params map[string]string, flagName, paramName string) {
	if cmd.Flags().Changed(flagName) {
		value, _ := cmd.Flags().GetString(flagName)
		params[paramName] = value
	}
}

func addBoolParam(cmd *cobra.Command, params map[string]string, flagName, paramName string) {
	if cmd.Flags().Changed(flagName) {
		value, _ := cmd.Flags().GetBool(flagName)
		if value {
			params[paramName] = "true"
		} else {
			params[paramName] = "false"
		}
	}
}

func addStringArrayParam(cmd *cobra.Command, params map[string]string, flagName, paramName string) {
	if cmd.Flags().Changed(flagName) {
		values, _ := cmd.Flags().GetStringArray(flagName)
		if len(values) == 0 {
			params[paramName] = ""
			return
		}
		params[paramName] = strings.Join(values, "\n")
	}
}

func runAction(action string, params map[string]string, opts things.ExecuteOptions) error {
	client, err := things.NewClient()
	if err != nil {
		formatter.PrintError("Failed to initialize Things client", "CLIENT_ERROR", err.Error())
		return nil
	}

	callback, err := client.Execute(action, params, opts)
	if err != nil {
		if cbErr, ok := err.(*things.CallbackError); ok {
			code := cbErr.Code
			if code == "" {
				code = "THINGS_ERROR"
			}
			formatter.PrintError(cbErr.Message, code, "")
			return nil
		}
		formatter.PrintError(fmt.Sprintf("Failed to execute Things action: %v", err), "THINGS_ERROR", err.Error())
		return nil
	}

	result := things.NormalizeResponse(action, callback)
	formatter.PrintSuccess(result)
	return nil
}

// addCmd creates a new to-do in Things
var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a new to-do in Things",
	Long: `Add a new to-do with title, notes, tags, and scheduling options.

Examples:
  things add --title "Buy milk" --when today --tags "errands"
  things add --titles "Buy milk" --titles "Send invoices" --when anytime
  things add --title "Review PR" --checklist-items "Read diff" --checklist-items "Run tests"`,
	RunE: func(cmd *cobra.Command, args []string) error {
		params := make(map[string]string)

		if cmd.Flags().Changed("titles") {
			addStringArrayParam(cmd, params, "titles", "titles")
		} else {
			addStringParam(cmd, params, "title", "title")
		}

		addStringParam(cmd, params, "notes", "notes")
		addStringParam(cmd, params, "when", "when")
		addStringParam(cmd, params, "deadline", "deadline")
		addStringParam(cmd, params, "tags", "tags")
		addStringParam(cmd, params, "list", "list")
		addStringParam(cmd, params, "list-id", "list-id")
		addStringParam(cmd, params, "heading", "heading")
		addStringParam(cmd, params, "heading-id", "heading-id")
		addStringParam(cmd, params, "use-clipboard", "use-clipboard")
		addStringParam(cmd, params, "creation-date", "creation-date")
		addStringParam(cmd, params, "completion-date", "completion-date")
		addStringArrayParam(cmd, params, "checklist-items", "checklist-items")
		addBoolParam(cmd, params, "completed", "completed")
		addBoolParam(cmd, params, "canceled", "canceled")
		addBoolParam(cmd, params, "show-quick-entry", "show-quick-entry")
		addBoolParam(cmd, params, "reveal", "reveal")

		return runAction("add", params, things.ExecuteOptions{})
	},
}

// addProjectCmd creates a new project in Things
var addProjectCmd = &cobra.Command{
	Use:   "add-project",
	Short: "Add a new project in Things",
	Long: `Add a new project with notes, tags, and optional area placement.

Examples:
  things add-project --title "Launch" --when someday
  things add-project --title "Website" --area "Work" --to-dos "Design" --to-dos "Build"`,
	RunE: func(cmd *cobra.Command, args []string) error {
		params := make(map[string]string)

		addStringParam(cmd, params, "title", "title")
		addStringParam(cmd, params, "notes", "notes")
		addStringParam(cmd, params, "when", "when")
		addStringParam(cmd, params, "deadline", "deadline")
		addStringParam(cmd, params, "tags", "tags")
		addStringParam(cmd, params, "area", "area")
		addStringParam(cmd, params, "area-id", "area-id")
		addStringArrayParam(cmd, params, "to-dos", "to-dos")
		addStringParam(cmd, params, "creation-date", "creation-date")
		addStringParam(cmd, params, "completion-date", "completion-date")
		addBoolParam(cmd, params, "completed", "completed")
		addBoolParam(cmd, params, "canceled", "canceled")
		addBoolParam(cmd, params, "reveal", "reveal")

		return runAction("add-project", params, things.ExecuteOptions{})
	},
}

// updateCmd modifies an existing to-do in Things
var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update an existing to-do in Things",
	Long: `Update a to-do by ID. Requires an auth token.

Examples:
  things update --id "THINGS-ID" --title "Updated title"
  things update --id "THINGS-ID" --prepend-notes "Urgent" --reveal`,
	RunE: func(cmd *cobra.Command, args []string) error {
		id, _ := cmd.Flags().GetString("id")
		if id == "" {
			formatter.PrintError("To-do ID (--id) is required", "INVALID_ARGUMENTS", "")
			return nil
		}

		params := map[string]string{"id": id}
		addStringParam(cmd, params, "title", "title")
		addStringParam(cmd, params, "notes", "notes")
		addStringParam(cmd, params, "prepend-notes", "prepend-notes")
		addStringParam(cmd, params, "append-notes", "append-notes")
		addStringParam(cmd, params, "when", "when")
		addStringParam(cmd, params, "deadline", "deadline")
		addStringParam(cmd, params, "tags", "tags")
		addStringParam(cmd, params, "add-tags", "add-tags")
		addStringArrayParam(cmd, params, "checklist-items", "checklist-items")
		addStringArrayParam(cmd, params, "prepend-checklist-items", "prepend-checklist-items")
		addStringArrayParam(cmd, params, "append-checklist-items", "append-checklist-items")
		addStringParam(cmd, params, "list", "list")
		addStringParam(cmd, params, "list-id", "list-id")
		addStringParam(cmd, params, "heading", "heading")
		addStringParam(cmd, params, "heading-id", "heading-id")
		addStringParam(cmd, params, "creation-date", "creation-date")
		addStringParam(cmd, params, "completion-date", "completion-date")
		addStringParam(cmd, params, "use-clipboard", "use-clipboard")
		addBoolParam(cmd, params, "completed", "completed")
		addBoolParam(cmd, params, "canceled", "canceled")
		addBoolParam(cmd, params, "reveal", "reveal")
		addBoolParam(cmd, params, "duplicate", "duplicate")
		addStringParam(cmd, params, "auth-token", "auth-token")

		return runAction("update", params, things.ExecuteOptions{RequiresAuth: true, UseAuthIfAvailable: true})
	},
}

// updateProjectCmd modifies an existing project in Things
var updateProjectCmd = &cobra.Command{
	Use:   "update-project",
	Short: "Update an existing project in Things",
	Long: `Update a project by ID. Requires an auth token.

Examples:
  things update-project --id "THINGS-ID" --title "Updated project" --reveal`,
	RunE: func(cmd *cobra.Command, args []string) error {
		id, _ := cmd.Flags().GetString("id")
		if id == "" {
			formatter.PrintError("Project ID (--id) is required", "INVALID_ARGUMENTS", "")
			return nil
		}

		params := map[string]string{"id": id}
		addStringParam(cmd, params, "title", "title")
		addStringParam(cmd, params, "notes", "notes")
		addStringParam(cmd, params, "prepend-notes", "prepend-notes")
		addStringParam(cmd, params, "append-notes", "append-notes")
		addStringParam(cmd, params, "when", "when")
		addStringParam(cmd, params, "deadline", "deadline")
		addStringParam(cmd, params, "tags", "tags")
		addStringParam(cmd, params, "add-tags", "add-tags")
		addStringParam(cmd, params, "area", "area")
		addStringParam(cmd, params, "area-id", "area-id")
		addStringParam(cmd, params, "creation-date", "creation-date")
		addStringParam(cmd, params, "completion-date", "completion-date")
		addBoolParam(cmd, params, "completed", "completed")
		addBoolParam(cmd, params, "canceled", "canceled")
		addBoolParam(cmd, params, "reveal", "reveal")
		addBoolParam(cmd, params, "duplicate", "duplicate")
		addStringParam(cmd, params, "auth-token", "auth-token")

		return runAction("update-project", params, things.ExecuteOptions{RequiresAuth: true, UseAuthIfAvailable: true})
	},
}

// showCmd shows a list or item in Things
var showCmd = &cobra.Command{
	Use:   "show",
	Short: "Show a list or item in Things",
	Long: `Show a list (by query) or a specific item by ID.

Examples:
  things show --query Today
  things show --id "THINGS-ID"`,
	RunE: func(cmd *cobra.Command, args []string) error {
		params := make(map[string]string)
		addStringParam(cmd, params, "id", "id")
		addStringParam(cmd, params, "query", "query")

		if len(params) == 0 {
			formatter.PrintError("Provide --id or --query", "INVALID_ARGUMENTS", "")
			return nil
		}

		return runAction("show", params, things.ExecuteOptions{})
	},
}

// searchCmd searches Things
var searchCmd = &cobra.Command{
	Use:   "search",
	Short: "Search in Things",
	Long: `Search Things using a query string.

Example:
  things search --query "project"`,
	RunE: func(cmd *cobra.Command, args []string) error {
		params := make(map[string]string)
		addStringParam(cmd, params, "query", "query")
		if len(params) == 0 {
			formatter.PrintError("Provide --query", "INVALID_ARGUMENTS", "")
			return nil
		}

		return runAction("search", params, things.ExecuteOptions{})
	},
}

// jsonCmd sends JSON payloads to Things
var jsonCmd = &cobra.Command{
	Use:   "json",
	Short: "Send a JSON payload to Things",
	Long: `Send JSON data to Things for batch creation or updates.

Examples:
  things json --file payload.json
  things json --data '{"items":[]}'`,
	RunE: func(cmd *cobra.Command, args []string) error {
		data, _ := cmd.Flags().GetString("data")
		filePath, _ := cmd.Flags().GetString("file")

		if filePath != "" {
			expanded, err := util.ExpandHomePath(filePath)
			if err != nil {
				formatter.PrintError("Invalid file path", "INVALID_ARGUMENTS", err.Error())
				return nil
			}
			payload, err := os.ReadFile(expanded)
			if err != nil {
				formatter.PrintError("Failed to read JSON file", "FILE_ERROR", err.Error())
				return nil
			}
			data = string(payload)
		}

		if strings.TrimSpace(data) == "" {
			formatter.PrintError("Provide --data or --file", "INVALID_ARGUMENTS", "")
			return nil
		}

		params := make(map[string]string)
		params["data"] = data
		addBoolParam(cmd, params, "reveal", "reveal")
		addStringParam(cmd, params, "auth-token", "auth-token")

		return runAction("json", params, things.ExecuteOptions{UseAuthIfAvailable: true})
	},
}

// versionCmd displays the Things URL scheme version
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show Things URL scheme version",
	RunE: func(cmd *cobra.Command, args []string) error {
		return runAction("version", map[string]string{}, things.ExecuteOptions{})
	},
}

// configCmd manages CLI configuration
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage Things CLI configuration",
}

var configSetTokenCmd = &cobra.Command{
	Use:   "set-token",
	Short: "Store Things auth token",
	RunE: func(cmd *cobra.Command, args []string) error {
		token, _ := cmd.Flags().GetString("auth-token")
		if token == "" {
			formatter.PrintError("Auth token (--auth-token) is required", "INVALID_ARGUMENTS", "")
			return nil
		}

		if err := util.SetAuthToken(token); err != nil {
			formatter.PrintError("Failed to save auth token", "CONFIG_ERROR", err.Error())
			return nil
		}

		formatter.PrintSuccess(map[string]interface{}{
			"status": "auth token saved",
		})
		return nil
	},
}

var configGetTokenCmd = &cobra.Command{
	Use:   "get-token",
	Short: "Display stored auth token (masked)",
	RunE: func(cmd *cobra.Command, args []string) error {
		token, err := util.GetAuthToken()
		if err != nil || token == "" {
			formatter.PrintError("No auth token configured", "CONFIG_ERROR", "Set one with 'things config set-token --auth-token YOUR_TOKEN'")
			return nil
		}

		formatter.PrintSuccess(map[string]interface{}{
			"auth_token": util.MaskToken(token),
		})
		return nil
	},
}

var configShowCmd = &cobra.Command{
	Use:   "show",
	Short: "Show current configuration",
	RunE: func(cmd *cobra.Command, args []string) error {
		config, err := util.LoadConfig()
		if err != nil {
			formatter.PrintError("Failed to load config", "CONFIG_ERROR", err.Error())
			return nil
		}

		configPath, _ := util.ConfigPath()
		tokenDisplay := "not set"
		if config.AuthToken != "" {
			tokenDisplay = util.MaskToken(config.AuthToken)
		}

		response := map[string]interface{}{
			"auth_token_set":        config.AuthToken != "",
			"auth_token":            tokenDisplay,
			"callback_port":         config.CallbackPort,
			"timeout_sec":           config.CallbackTimeoutSeconds,
			"output_format":         config.OutputFormat,
			"config_path":           configPath,
			"last_updated":          config.LastUpdated,
		}

		formatter.PrintSuccess(response)
		return nil
	},
}

func init() {
	addCmd.Flags().String("title", "", "To-do title")
	addCmd.Flags().StringArray("titles", []string{}, "Multiple to-do titles (repeat flag)")
	addCmd.Flags().String("notes", "", "Notes for the to-do")
	addCmd.Flags().String("when", "", "When to schedule (today, tonight, anytime, someday, or date)")
	addCmd.Flags().String("deadline", "", "Deadline date (YYYY-MM-DD)")
	addCmd.Flags().String("tags", "", "Comma-separated tags")
	addCmd.Flags().String("list", "", "List name or project title")
	addCmd.Flags().String("list-id", "", "List or project ID")
	addCmd.Flags().String("heading", "", "Heading title")
	addCmd.Flags().String("heading-id", "", "Heading ID")
	addCmd.Flags().StringArray("checklist-items", []string{}, "Checklist items (repeat flag)")
	addCmd.Flags().Bool("completed", false, "Mark as completed")
	addCmd.Flags().Bool("canceled", false, "Mark as canceled")
	addCmd.Flags().Bool("show-quick-entry", false, "Show quick entry after adding")
	addCmd.Flags().Bool("reveal", false, "Reveal the created to-do in Things")
	addCmd.Flags().String("creation-date", "", "Creation date (ISO 8601)")
	addCmd.Flags().String("completion-date", "", "Completion date (ISO 8601)")
	addCmd.Flags().String("use-clipboard", "", "Use clipboard content (replace-title|replace-notes|replace-checklist-items)")

	addProjectCmd.Flags().String("title", "", "Project title")
	addProjectCmd.Flags().String("notes", "", "Project notes")
	addProjectCmd.Flags().String("when", "", "When to schedule (today, tonight, anytime, someday, or date)")
	addProjectCmd.Flags().String("deadline", "", "Deadline date (YYYY-MM-DD)")
	addProjectCmd.Flags().String("tags", "", "Comma-separated tags")
	addProjectCmd.Flags().String("area", "", "Area name")
	addProjectCmd.Flags().String("area-id", "", "Area ID")
	addProjectCmd.Flags().StringArray("to-dos", []string{}, "Project to-dos (repeat flag)")
	addProjectCmd.Flags().Bool("completed", false, "Mark as completed")
	addProjectCmd.Flags().Bool("canceled", false, "Mark as canceled")
	addProjectCmd.Flags().Bool("reveal", false, "Reveal the created project in Things")
	addProjectCmd.Flags().String("creation-date", "", "Creation date (ISO 8601)")
	addProjectCmd.Flags().String("completion-date", "", "Completion date (ISO 8601)")

	updateCmd.Flags().String("id", "", "To-do ID (required)")
	updateCmd.Flags().String("title", "", "Updated title")
	updateCmd.Flags().String("notes", "", "Replace notes")
	updateCmd.Flags().String("prepend-notes", "", "Prepend notes")
	updateCmd.Flags().String("append-notes", "", "Append notes")
	updateCmd.Flags().String("when", "", "Update schedule")
	updateCmd.Flags().String("deadline", "", "Update deadline")
	updateCmd.Flags().String("tags", "", "Replace tags")
	updateCmd.Flags().String("add-tags", "", "Add tags")
	updateCmd.Flags().StringArray("checklist-items", []string{}, "Replace checklist items (repeat flag)")
	updateCmd.Flags().StringArray("prepend-checklist-items", []string{}, "Prepend checklist items (repeat flag)")
	updateCmd.Flags().StringArray("append-checklist-items", []string{}, "Append checklist items (repeat flag)")
	updateCmd.Flags().String("list", "", "Move to list by name")
	updateCmd.Flags().String("list-id", "", "Move to list by ID")
	updateCmd.Flags().String("heading", "", "Move to heading by name")
	updateCmd.Flags().String("heading-id", "", "Move to heading by ID")
	updateCmd.Flags().Bool("completed", false, "Mark as completed")
	updateCmd.Flags().Bool("canceled", false, "Mark as canceled")
	updateCmd.Flags().Bool("reveal", false, "Reveal the updated to-do")
	updateCmd.Flags().Bool("duplicate", false, "Duplicate the to-do")
	updateCmd.Flags().String("creation-date", "", "Set creation date (ISO 8601)")
	updateCmd.Flags().String("completion-date", "", "Set completion date (ISO 8601)")
	updateCmd.Flags().String("use-clipboard", "", "Use clipboard content (replace-title|replace-notes|replace-checklist-items)")
	updateCmd.Flags().String("auth-token", "", "Things auth token (overrides config/ENV)")

	updateProjectCmd.Flags().String("id", "", "Project ID (required)")
	updateProjectCmd.Flags().String("title", "", "Updated title")
	updateProjectCmd.Flags().String("notes", "", "Replace notes")
	updateProjectCmd.Flags().String("prepend-notes", "", "Prepend notes")
	updateProjectCmd.Flags().String("append-notes", "", "Append notes")
	updateProjectCmd.Flags().String("when", "", "Update schedule")
	updateProjectCmd.Flags().String("deadline", "", "Update deadline")
	updateProjectCmd.Flags().String("tags", "", "Replace tags")
	updateProjectCmd.Flags().String("add-tags", "", "Add tags")
	updateProjectCmd.Flags().String("area", "", "Move to area by name")
	updateProjectCmd.Flags().String("area-id", "", "Move to area by ID")
	updateProjectCmd.Flags().Bool("completed", false, "Mark as completed")
	updateProjectCmd.Flags().Bool("canceled", false, "Mark as canceled")
	updateProjectCmd.Flags().Bool("reveal", false, "Reveal the updated project")
	updateProjectCmd.Flags().Bool("duplicate", false, "Duplicate the project")
	updateProjectCmd.Flags().String("creation-date", "", "Set creation date (ISO 8601)")
	updateProjectCmd.Flags().String("completion-date", "", "Set completion date (ISO 8601)")
	updateProjectCmd.Flags().String("auth-token", "", "Things auth token (overrides config/ENV)")

	showCmd.Flags().String("id", "", "Item ID to show")
	showCmd.Flags().String("query", "", "List query (Inbox, Today, Upcoming, etc)")

	searchCmd.Flags().String("query", "", "Search query")

	jsonCmd.Flags().String("data", "", "JSON payload string")
	jsonCmd.Flags().String("file", "", "Path to JSON payload file")
	jsonCmd.Flags().Bool("reveal", false, "Reveal created items")
	jsonCmd.Flags().String("auth-token", "", "Things auth token (overrides config/ENV)")

	configSetTokenCmd.Flags().String("auth-token", "", "Things auth token")

	configCmd.AddCommand(configSetTokenCmd)
	configCmd.AddCommand(configGetTokenCmd)
	configCmd.AddCommand(configShowCmd)
}

// GetCommands returns all available commands for the root command
func GetCommands() []*cobra.Command {
	return []*cobra.Command{
		addCmd,
		addProjectCmd,
		updateCmd,
		updateProjectCmd,
		showCmd,
		searchCmd,
		jsonCmd,
		versionCmd,
		configCmd,
	}
}
