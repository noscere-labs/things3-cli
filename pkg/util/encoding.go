package util

import (
	"net/url"
	"os"
	"strings"
	"time"
)

// EncodeParam URL-encodes a string parameter for use in Things URLs
// This ensures special characters are properly escaped
func EncodeParam(value string) string {
	return url.QueryEscape(value)
}

// EncodeParams takes a map of string parameters and returns a URL-encoded query string
// Useful for building the parameter section of x-callback-urls
func EncodeParams(params map[string]string) string {
	v := url.Values{}
	for key, value := range params {
		v.Add(key, value)
	}
	// Encode and replace + with %20 for proper space encoding in x-callback-urls
	// Things' callback mechanism doesn't properly decode + as space, so we use %20 instead
	encoded := v.Encode()
	return strings.ReplaceAll(encoded, "+", "%20")
}

// GetTimestamp returns current date/time formatted for prepending to notes
// Format: "2024-01-15 10:30:00"
func GetTimestamp() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

// ParseTags takes a comma-separated string and returns a slice of trimmed tag strings
// Example: "work,urgent,projects" → []string{"work", "urgent", "projects"}
func ParseTags(tagString string) []string {
	if tagString == "" {
		return []string{}
	}

	tags := strings.Split(tagString, ",")
	var result []string
	for _, tag := range tags {
		trimmed := strings.TrimSpace(tag)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}

// JoinTags takes a slice of tags and returns them as a comma-separated string
// Inverse of ParseTags
func JoinTags(tags []string) string {
	return strings.Join(tags, ",")
}

// ExpandHomePath expands ~ to the user's home directory
// Example: "~/Documents/note.txt" → "/Users/username/Documents/note.txt"
func ExpandHomePath(path string) (string, error) {
	if strings.HasPrefix(path, "~") {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		return strings.Replace(path, "~", home, 1), nil
	}
	return path, nil
}
