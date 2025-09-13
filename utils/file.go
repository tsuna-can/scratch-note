package utils

import (
	"fmt"
	"regexp"
	"strings"
	"time"
)

// GenerateFileName generates a timestamped filename with optional title
func GenerateFileName(title string, t time.Time) string {
	// Format: YYYY-MM-DD_HHMMSS
	timestamp := t.Format("2006-01-02_150405")
	
	if title == "" {
		return timestamp + ".md"
	}
	
	// Clean title: replace spaces and special characters with hyphens
	cleanTitle := cleanTitle(title)
	return fmt.Sprintf("%s_%s.md", timestamp, cleanTitle)
}

// cleanTitle removes special characters and replaces spaces with hyphens
func cleanTitle(title string) string {
	// Replace spaces with hyphens
	cleaned := strings.ReplaceAll(title, " ", "-")
	
	// Remove or replace special characters with hyphens
	reg := regexp.MustCompile(`[/\\:*?"<>|]`)
	cleaned = reg.ReplaceAllString(cleaned, "-")
	
	// Remove multiple consecutive hyphens
	reg = regexp.MustCompile(`-+`)
	cleaned = reg.ReplaceAllString(cleaned, "-")
	
	// Trim hyphens from start and end
	cleaned = strings.Trim(cleaned, "-")
	
	return cleaned
}