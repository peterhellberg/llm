package boolparser

import (
	"fmt"
	"testing"
)

func TestNew(t *testing.T) {
	parser := New()

	if got, want := parser.FormatInstructions(), parserFormatInstructions; got != want {
		t.Fatalf("parser.FormatInstructions() = %q, want %q", got, want)
	}

	if got, want := parser.Type(), parserType; got != want {
		t.Fatalf("parser.Type() = %q, want %q", got, want)
	}

	for _, tt := range []struct {
		str  string
		want bool
	}{
		{"TRUE", true},
		{"tRue", true},
		{"Yes", true},
		{"yes", true},
		{"FALSE", false},
		{"fAlse", false},
		{"No", false},
		{"no", false},
	} {
		t.Run(fmt.Sprintf("%v/%s", tt.want, tt.str), func(t *testing.T) {
			got, err := parser.parse(tt.str)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if got != tt.want {
				t.Fatalf("got = %v, want %v", got, tt.want)
			}
		})
	}
}
