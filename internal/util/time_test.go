package util

import (
	"testing"
	"time"


)


func TestParseDate(t *testing.T) {
	today := time.Now().Format("2006-01-02")
	futureDate := time.Now().AddDate(0, 0, 1).Format("2006-01-02")
	pastDate := time.Now().AddDate(0, 0, -1).Format("2006-01-02")

	tests := []struct {
		name        string
		input       string
		expectError bool
		errorMsg    string
	}{
		{
			name:        "valid future date",
			input:       futureDate,
			expectError: false,
		},
		{
			name:        "valid today date",
			input:       today,
			expectError: false,
		},
		{
			name:        "past date",
			input:       pastDate,
			expectError: true,
			errorMsg:    "deadline cannot be in the past",
		},
		{
			name:        "invalid format",
			input:       "2025/08/30",
			expectError: true,
			errorMsg:    "invalid date format, expected YYYY-MM-DD",
		},
		{
			name:        "empty string",
			input:       "",
			expectError: true,
			errorMsg:    "invalid date format, expected YYYY-MM-DD",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parsed, err := ParseDate(tt.input)
			if tt.expectError {
				if err == nil {
					t.Fatalf("expected error but got nil")
				}
				if err.Error() != tt.errorMsg {
					t.Fatalf("expected error %q, got %q", tt.errorMsg, err.Error())
				}
			} else {
				if err != nil {
					t.Fatalf("expected no error, got %v", err)
				}
				// Check parsed date matches input string
				if parsed.Format("2006-01-02") != tt.input {
					t.Fatalf("expected date %v, got %v", tt.input, parsed.Format("2006-01-02"))
				}
			}
		})
	}
}
