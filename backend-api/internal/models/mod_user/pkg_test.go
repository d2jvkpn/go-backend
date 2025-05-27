package mod_user

import (
	// "fmt"
	"regexp"
	"testing"
)

func TestPassword(t *testing.T) {
	// pattern := `^(?=.*[a-z])(?=.*[A-Z])(?=.*\d)(?=.*[!@#$%^&*])[a-zA-Z\d!@#$%^&*]{8,32}$`
	// reg := regexp.MustCompile(pattern) // error parsing regexp: invalid or unsupported Perl syntax: `(?=` [recovered]

	legal := regexp.MustCompile(`^[a-zA-Z\d!@#$%^&*]{8,32}$`)
	legals := []*regexp.Regexp{
		regexp.MustCompile(`[a-z]+`),
		regexp.MustCompile(`[A-Z]+`),
		regexp.MustCompile(`[\d]+`),
		regexp.MustCompile(`[!@#$%^&*]+`),
	}

	validate := func(input string) bool {
		if !legal.MatchString(input) {
			return false
		}
		for i := range legals {
			if !legals[i].MatchString(input) {
				return false
			}
		}

		return true
	}

	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{"Valid password", "Abc12345!", true},
		{"Valid password", "Aa1@#$%^&*", true},
		{"Missing uppercase", "abc12345", false},
		{"Missing lowercase", "ABC12345!", false},
		{"Missing digit", "Abcdefgh!", false},
		{"Contains invalid symbol", "Abc12345~", false},
		{"Contains Chinese", "Abc12345!中文", false},
		{"Too short", "Ab1!", false},
		{"Too long", "Abc123456789012345678901234567890!", false},
	}

	for _, item := range tests {
		t.Run(item.name, func(t *testing.T) {
			actual := validate(item.input)
			if actual != item.expected {
				t.Errorf("expected %v, got %v (input: %s)", item.expected, actual, item.input)
			}
		})
	}
}
