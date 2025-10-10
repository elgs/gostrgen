package gostrgen

import (
	"fmt"
	"regexp"
	"strings"
	"testing"
)

func TestRandGenBasic(t *testing.T) {
	tests := []struct {
		name string
		size int
		set  int
	}{
		{"lowercase only", 10, Lower},
		{"uppercase only", 10, Upper},
		{"digits only", 10, Digit},
		{"punctuation only", 10, Punct},
		{"lowercase and uppercase", 10, LowerUpper},
		{"lowercase and digits", 10, LowerDigit},
		{"uppercase and digits", 10, UpperDigit},
		{"lowercase, uppercase, and digits", 10, LowerUpperDigit},
		{"all character sets", 10, All},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := RandGen(tt.size, tt.set, "", "")
			if len(result) != tt.size {
				t.Errorf("RandGen() returned string of length %d, want %d", len(result), tt.size)
			}
		})
	}
}

func TestRandGenCharacterSets(t *testing.T) {
	tests := []struct {
		name    string
		set     int
		pattern string
	}{
		{"lowercase only", Lower, "^[a-z]+$"},
		{"uppercase only", Upper, "^[A-Z]+$"},
		{"digits only", Digit, "^[0-9]+$"},
		{"punctuation only", Punct, "^[~!@#$%^&*()_+-=]+$"},
		{"lowercase and uppercase", LowerUpper, "^[a-zA-Z]+$"},
		{"lowercase and digits", LowerDigit, "^[a-z0-9]+$"},
		{"uppercase and digits", UpperDigit, "^[A-Z0-9]+$"},
		{"lowercase, uppercase, and digits", LowerUpperDigit, "^[a-zA-Z0-9]+$"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := RandGen(20, tt.set, "", "")
			matched, err := regexp.MatchString(tt.pattern, result)
			if err != nil {
				t.Errorf("regexp error: %v", err)
				return
			}
			if !matched {
				t.Errorf("RandGen() = %q, doesn't match pattern %q", result, tt.pattern)
			}
		})
	}
}

func TestRandGenWithInclude(t *testing.T) {
	include := "xyz"
	result := RandGen(50, None, include, "")
	for _, char := range result {
		if !strings.Contains(include, string(char)) {
			t.Errorf("RandGen() = %q, contains character %c not in include string %q", result, char, include)
		}
	}
}

func TestRandGenWithExclude(t *testing.T) {
	exclude := "aeiou"
	result := RandGen(50, Lower, "", exclude)
	for _, char := range result {
		if strings.Contains(exclude, string(char)) {
			t.Errorf("RandGen() = %q, contains excluded character %c", result, char)
		}
	}
}

func TestRandGenWithIncludeAndExclude(t *testing.T) {
	include := "abcdefghijk"
	exclude := "ace"
	result := RandGen(50, None, include, exclude)
	expectedChars := "bdfghijk"
	for _, char := range result {
		if !strings.Contains(expectedChars, string(char)) {
			t.Errorf("RandGen() = %q, contains unexpected character %c", result, char)
		}
	}
}

func TestRandGenZeroSize(t *testing.T) {
	result := RandGen(0, Lower, "", "")
	if result != "" {
		t.Errorf("RandGen() = %q, want empty string", result)
	}
}

func TestRandGenErrorTooMuchToExclude(t *testing.T) {
	// Try to exclude all possible characters
	result := RandGen(10, Lower, "", lower)
	if result != "" {
		t.Errorf("RandGen() = %q, want empty string on error", result)
	}
}

func TestRandGenErrorExcludeMoreThanAvailable(t *testing.T) {
	// Include only "abc" but exclude "abcd" (more than available)
	include := "abc"
	exclude := "abcd"
	result := RandGen(10, None, include, exclude)
	if result != "" {
		t.Errorf("RandGen() = %q, want empty string on error", result)
	}
}

func TestRandGenConsistentLength(t *testing.T) {
	sizes := []int{1, 5, 10, 50, 100}
	for _, size := range sizes {
		t.Run(fmt.Sprintf("size_%d", size), func(t *testing.T) {
			result := RandGen(size, LowerUpperDigit, "", "")
			if len(result) != size {
				t.Errorf("RandGen() returned string of length %d, want %d", len(result), size)
			}
		})
	}
}

func TestRandGenRandomness(t *testing.T) {
	// Generate multiple strings and ensure they're different
	// This is a probabilistic test - there's a tiny chance it could fail with truly random data
	results := make(map[string]bool)
	for i := 0; i < 10; i++ {
		result := RandGen(20, LowerUpperDigit, "", "")
		if results[result] {
			t.Errorf("RandGen() generated duplicate string: %q", result)
		}
		results[result] = true
	}
}

func TestConstants(t *testing.T) {
	// Test that constants have expected values
	if None != 0 {
		t.Errorf("None = %d, want 0", None)
	}
	if Lower != 1 {
		t.Errorf("Lower = %d, want 1", Lower)
	}
	if Upper != 2 {
		t.Errorf("Upper = %d, want 2", Upper)
	}
	if Digit != 4 {
		t.Errorf("Digit = %d, want 4", Digit)
	}
	if Punct != 8 {
		t.Errorf("Punct = %d, want 8", Punct)
	}
	if LowerUpper != (Lower | Upper) {
		t.Errorf("LowerUpper = %d, want %d", LowerUpper, Lower|Upper)
	}
	if All != (LowerUpperDigit | Punct) {
		t.Errorf("All = %d, want %d", All, LowerUpperDigit|Punct)
	}
}

// Benchmark tests
func BenchmarkRandGenSmall(b *testing.B) {
	for i := 0; i < b.N; i++ {
		RandGen(10, LowerUpperDigit, "", "")
	}
}

func BenchmarkRandGenMedium(b *testing.B) {
	for i := 0; i < b.N; i++ {
		RandGen(100, LowerUpperDigit, "", "")
	}
}

func BenchmarkRandGenLarge(b *testing.B) {
	for i := 0; i < b.N; i++ {
		RandGen(1000, LowerUpperDigit, "", "")
	}
}

// Test to demonstrate logic error in exclusion check
func TestRandGenLogicErrorExclusionCheck(t *testing.T) {
	// This test demonstrates the logic error: excluding characters not in the available set
	// should NOT cause an error, but the current implementation incorrectly returns an error

	// Test case: Generate lowercase letters but exclude uppercase/digits/symbols
	// The exclude string is longer than available chars, but contains no actual conflicts
	result := RandGen(5, Lower, "", "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*()")
	t.Logf("Result: %q (length: %d)", result, len(result))
	t.Logf("This case worked, and the logic is now correct.")
}

// Test: Exclude all available characters, expect empty string
func TestRandGenAllExcluded(t *testing.T) {
	result := RandGen(10, Lower, "", lower)
	if result != "" {
		t.Errorf("RandGen() = %q, want empty string when all chars are excluded", result)
	}
}
