package formatter

import (
	"encoding/json"
	"fmt"
)

// FormatSuccess formats a successful operation response as JSON
// data: The data to include in the response (can be any type)
func FormatSuccess(data interface{}) string {
	response := map[string]interface{}{
		"success": true,
		"data":    data,
	}
	return formatAsJSON(response)
}

// FormatError formats an error response as JSON
// errorMsg: Human-readable error message
// code: Machine-readable error code (e.g., "NOTE_NOT_FOUND")
// details: Optional additional details about the error
func FormatError(errorMsg string, code string, details string) string {
	response := map[string]interface{}{
		"success":    false,
		"error":      errorMsg,
		"error_code": code,
	}

	if details != "" {
		response["details"] = details
	}

	return formatAsJSON(response)
}

// FormatNoteList formats a list of notes as JSON
// This is a convenience function for consistent list formatting
func FormatNoteList(notes interface{}, count int) string {
	response := map[string]interface{}{
		"success": true,
		"count":   count,
		"data":    notes,
	}
	return formatAsJSON(response)
}

// FormatTagList formats a list of tags as JSON
// This is a convenience function for consistent tag list formatting
func FormatTagList(tags interface{}, count int) string {
	response := map[string]interface{}{
		"success": true,
		"count":   count,
		"data":    tags,
	}
	return formatAsJSON(response)
}

// formatAsJSON converts any Go value to pretty-printed JSON
// This ensures consistent formatting across all output
func formatAsJSON(v interface{}) string {
	// Marshal with indentation for readability
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		// If marshaling fails, return an error response
		fallback := map[string]interface{}{
			"success":    false,
			"error":      "Failed to format response",
			"error_code": "FORMAT_ERROR",
		}
		if data, err := json.MarshalIndent(fallback, "", "  "); err == nil {
			return string(data)
		}
		return `{"success": false, "error": "Critical formatting error"}`
	}

	return string(data)
}

// PrintJSON prints a JSON response to stdout
// This is the main function called by CLI commands to output results
func PrintJSON(v interface{}) {
	output := formatAsJSON(v)
	fmt.Println(output)
}

// PrintSuccess prints a success response to stdout
func PrintSuccess(data interface{}) {
	PrintJSON(map[string]interface{}{
		"success": true,
		"data":    data,
	})
}

// PrintError prints an error response to stdout
func PrintError(errorMsg string, code string, details string) {
	response := map[string]interface{}{
		"success":    false,
		"error":      errorMsg,
		"error_code": code,
	}

	if details != "" {
		response["details"] = details
	}

	PrintJSON(response)
}
