package uuid

import (
	"bytes"
	"testing"
)

func TestUuidServiceParse(t *testing.T) {
	var tests = []struct {
		name      string
		input     string
		wantValid bool
	}{
		{"Valid UUID", "550e8400-e29b-41d4-a716-446655440000", true},
		{"Valid UUID", "00000000-0000-0000-0000-000000000000", true},
		{"Valid UUID", "FFFFFFFF-FFFF-FFFF-FFFF-FFFFFFFFFFFF", true},
		{"Valid UUID", "00000000000000000000000000000000", true},
		{"Invalid UUID", "FFFFFFFF-FFFF-FFFF-FFFF-FFFFFFFFFFF", false},
		{"Invalid UUID", "FFFFFFFF-FFFF-FFFF-FFFF-FFFFFFFFFFFF-FFFF", false},
		{"Invalid UUID", "FFFFFFFF-FFFF-FFFF-FFFF-FFFFFFFFFFFF123", false},
		{"Invalid UUID", "FFFFFFFF-FFFF-FFFF-FFFF-FFFFFFFFFFFV", false},
		{"Invalid UUID", "00000000 0000 0000 0000 000000000000", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			us := NewUUIDService()
			_, err := us.Parse(tt.input)
			isValid := err == nil

			if isValid != tt.wantValid {
				t.Errorf("got valid %v, want valid %v", isValid, tt.wantValid)
			}
		})
	}
}

func TestUuidServiceNew(t *testing.T) {
	us := NewUUIDService()
	uuid := us.New()

	if len(uuid) != 16 {
		t.Errorf("uuid length is not 16 bytes: %s", uuid)
	}
}
func TestUuidServiceFromBytes(t *testing.T) {
	us := NewUUIDService()

	// Test cases
	tests := []struct {
		name      string
		input     []byte
		wantValid bool
	}{
		{
			name:      "Valid 16 byte UUID",
			input:     make([]byte, 16),
			wantValid: true,
		},
		{
			name:      "Invalid byte length",
			input:     make([]byte, 15),
			wantValid: false,
		},
		{
			name:      "Nil bytes",
			input:     nil,
			wantValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uuid, err := us.FromBytes(tt.input)
			isValid := err == nil

			if isValid != tt.wantValid {
				t.Errorf("got valid %v, want valid %v", isValid, tt.wantValid)
			}
			uuidBytes := uuid[:]

			if tt.wantValid && !bytes.Equal(uuidBytes, tt.input) {
				t.Errorf("uuid value changed")
			}
		})
	}
}
