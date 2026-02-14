package validator

import "testing"

func TestIsValidCEP(t *testing.T) {
	tests := []struct {
		name string
		cep  string
		want bool
	}{
		{"valid", "29902555", true},
		{"too short", "123", false},
		{"too long", "123456789", false},
		{"non digit", "1234abcd", false},
		{"empty", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsValidCEP(tt.cep); got != tt.want {
				t.Fatalf("IsValidCEP(%q)=%v, want %v", tt.cep, got, tt.want)
			}
		})
	}
}
