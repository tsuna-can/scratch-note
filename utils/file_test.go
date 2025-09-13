package utils

import (
	"strings"
	"testing"
	"time"
)

func TestGenerateFileName(t *testing.T) {
	tests := []struct {
		name     string
		title    string
		time     time.Time
		expected string
	}{
		{
			name:     "without title",
			title:    "",
			time:     time.Date(2025, 8, 16, 14, 30, 45, 0, time.UTC),
			expected: "2025-08-16_143045.md",
		},
		{
			name:     "with title",
			title:    "shopping-list",
			time:     time.Date(2025, 8, 16, 14, 30, 45, 0, time.UTC),
			expected: "2025-08-16_143045_shopping-list.md",
		},
		{
			name:     "with title containing spaces",
			title:    "my daily notes",
			time:     time.Date(2025, 12, 25, 9, 5, 30, 0, time.UTC),
			expected: "2025-12-25_090530_my-daily-notes.md",
		},
		{
			name:     "with special characters in title",
			title:    "notes/with\\special:chars",
			time:     time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
			expected: "2025-01-01_000000_notes-with-special-chars.md",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GenerateFileName(tt.title, tt.time)
			if result != tt.expected {
				t.Errorf("GenerateFileName(%q, %v) = %q, want %q", tt.title, tt.time, result, tt.expected)
			}
		})
	}
}

func TestGenerateFileNameCurrentTime(t *testing.T) {
	// Test that the function works with current time
	testTime := time.Now()
	result := GenerateFileName("", testTime)
	
	t.Logf("Generated filename: %s (length: %d)", result, len(result))

	// Check format
	if !strings.HasSuffix(result, ".md") {
		t.Errorf("Generated filename should end with .md, got: %s", result)
	}

	// Check length (YYYY-MM-DD_HHMMSS.md = 20 characters)
	if len(result) != 20 {
		t.Errorf("Generated filename should be 20 characters long, got: %d", len(result))
	}

	// Verify time format by parsing
	timeStr := result[:17] // Extract YYYY-MM-DD_HHMMSS part (17 characters)
	_, err := time.Parse("2006-01-02_150405", timeStr)
	if err != nil {
		t.Errorf("Failed to parse time from filename: %v", err)
	}

	// Check that parsed time matches input time (within same second)
	expectedStr := testTime.Format("2006-01-02_150405")
	if timeStr != expectedStr {
		t.Errorf("Time in filename %s does not match expected %s", timeStr, expectedStr)
	}
}