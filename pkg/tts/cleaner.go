package tts

import (
	"regexp"
	"strings"
)

// CleanTextForTTS cleans Bear markdown text for TTS processing
// It removes markdown syntax and adds natural pauses for better speech flow
func CleanTextForTTS(text string) string {
	// Remove code blocks (```code```) - indicate they were skipped
	text = regexp.MustCompile("(?s)```[^`]*?```").ReplaceAllString(text, "")

	// Remove inline code (`code`) - just remove the backticks
	text = regexp.MustCompile("`([^`]+)`").ReplaceAllString(text, "$1")

	// Remove markdown links but keep the text: [text](url) -> text
	text = regexp.MustCompile(`\[([^\]]+)\]\([^\)]+\)`).ReplaceAllString(text, "$1")

	// Remove Bear wiki links but keep text: [[link]] -> link
	text = regexp.MustCompile(`\[\[([^\]]+)\]\]`).ReplaceAllString(text, "$1")

	// Remove Bear highlights: ::text:: -> text
	text = regexp.MustCompile(`::([^:]+)::`).ReplaceAllString(text, "$1")

	// Handle headings with pauses for natural flow
	// H1: # Heading -> "Heading. [longer pause]"
	text = regexp.MustCompile(`(?m)^#\s+(.+)$`).ReplaceAllString(text, "$1.\n\n")
	// H2: ## Heading -> "Heading. [pause]"
	text = regexp.MustCompile(`(?m)^##\s+(.+)$`).ReplaceAllString(text, "$1.\n\n")
	// H3-H6: ### Heading -> "Heading. [short pause]"
	text = regexp.MustCompile(`(?m)^#{3,6}\s+(.+)$`).ReplaceAllString(text, "$1.\n")

	// Handle horizontal rules (---, ***, ___) - replace with paragraph break
	text = regexp.MustCompile(`(?m)^[-*_]{3,}$`).ReplaceAllString(text, "\n\n")

	// Handle blockquotes: > text -> text (with slight pause)
	text = regexp.MustCompile(`(?m)^>\s+(.+)$`).ReplaceAllString(text, "$1.\n")

	// Handle list items - remove bullets/numbers, keep items on separate lines
	// Unordered lists: - item, * item, + item
	text = regexp.MustCompile(`(?m)^[\s]*[-*+]\s+`).ReplaceAllString(text, "")
	// Ordered lists: 1. item, 2. item
	text = regexp.MustCompile(`(?m)^[\s]*\d+\.\s+`).ReplaceAllString(text, "")

	// Remove todo checkboxes but keep the item text
	text = regexp.MustCompile(`\[\s*\]\s*`).ReplaceAllString(text, "")
	text = regexp.MustCompile(`\[x\]\s*`).ReplaceAllString(text, "")

	// Remove Bear tags (#tag/subtag)
	text = regexp.MustCompile(`#[a-zA-Z0-9/_-]+`).ReplaceAllString(text, "")

	// Remove markdown formatting (* _ ~ for bold/italic/strikethrough)
	text = regexp.MustCompile(`[*_~]{1,3}`).ReplaceAllString(text, "")

	// Convert multiple dots into natural pauses
	// ... remains as a pause
	// Replace --em-dash-- with comma for natural pause
	text = strings.ReplaceAll(text, " -- ", ", ")
	text = strings.ReplaceAll(text, "—", ", ")

	// Clean up excessive whitespace
	// Collapse multiple spaces to single space
	text = regexp.MustCompile(`[ \t]+`).ReplaceAllString(text, " ")

	// Collapse multiple newlines (3+) to double newline (paragraph break)
	text = regexp.MustCompile("\n{3,}").ReplaceAllString(text, "\n\n")

	// Replace double newlines with period + comma for natural paragraph pauses
	text = strings.ReplaceAll(text, "\n\n", "., ")

	// Replace single newlines with comma for brief pauses (list items, etc)
	text = strings.ReplaceAll(text, "\n", ", ")

	// Clean up punctuation issues
	// Remove spaces before periods/commas
	text = regexp.MustCompile(`\s+([.,])`).ReplaceAllString(text, "$1")
	// Clean up multiple commas
	text = regexp.MustCompile(`,{2,}`).ReplaceAllString(text, ",")
	// Clean up comma before period
	text = regexp.MustCompile(`,\s*\.`).ReplaceAllString(text, ".")
	// Clean up multiple periods (but preserve ellipsis)
	text = regexp.MustCompile(`\.{4,}`).ReplaceAllString(text, "...")
	// Ensure space after periods and commas
	text = regexp.MustCompile(`([.,])([A-Za-z])`).ReplaceAllString(text, "$1 $2")

	// Trim leading/trailing whitespace and punctuation
	text = strings.TrimSpace(text)
	// Remove trailing comma at end
	text = strings.TrimSuffix(text, ",")
	text = strings.TrimSpace(text)

	return text
}

// ValidateTextLength checks if text is within acceptable length range
func ValidateTextLength(text string, minLength, maxLength int) bool {
	length := len(text)
	return length >= minLength && length <= maxLength
}

// TruncateText truncates text to maxLength while trying to preserve complete words/sentences
func TruncateText(text string, maxLength int) string {
	if len(text) <= maxLength {
		return text
	}

	// Truncate to maxLength
	truncated := text[:maxLength]

	// Try to find the last space to avoid cutting mid-word
	lastSpace := strings.LastIndex(truncated, " ")
	if lastSpace > maxLength-50 { // Only cut at space if it's close to the end
		truncated = truncated[:lastSpace]
	}

	return strings.TrimSpace(truncated) + "..."
}
