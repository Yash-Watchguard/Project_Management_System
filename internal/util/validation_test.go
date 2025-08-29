package util

import (
	"testing"
)


func TestValidateEmail(t *testing.T) {
	tests := []struct {
		name      string
		email     string
		expectErr bool
	}{
		{"valid email", "test@example.com", false},
		{"uppercase email", "Test@Example.COM", false},
		{"missing @", "testexample.com", true},
		{"invalid domain", "test@.com", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err :=ValidateEmail(tt.email)
			if (err != nil) != tt.expectErr {
				t.Errorf("ValidateEmail(%q) error = %v, wantErr %v", tt.email, err, tt.expectErr)
			}
		})
	}
}

func TestValidateMobileNumber(t *testing.T) {
	tests := []struct {
		name      string
		phone     string
		expectErr bool
	}{
		{"valid phone", "9876543210", false},
		{"starts with 5", "5876543210", true},
		{"less digits", "98765432", true},
		{"more digits", "987654321012", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateMobileNumber(tt.phone)
			if (err != nil) != tt.expectErr {
				t.Errorf("ValidateMobileNumber(%q) error = %v, wantErr %v", tt.phone, err, tt.expectErr)
			}
		})
	}
}

func TestValidatePassword(t *testing.T) {
	tests := []struct {
		name      string
		password  string
		expectErr bool
	}{
		{"valid password", "Abcdef123!@#", false},
		{"too short", "Abc123!@#", true},
		{"missing uppercase", "abcdef123!@#", true},
		{"missing lowercase", "ABCDEF123!@#", true},
		{"missing digit", "Abcdefghijk!@", true},
		{"missing special", "Abcdef123456", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err :=ValidatePassword(tt.password)
			if (err != nil) != tt.expectErr {
				t.Errorf("ValidatePassword(%q) error = %v, wantErr %v", tt.password, err, tt.expectErr)
			}
		})
	}
}
