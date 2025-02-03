package regexparser

import "testing"

func TestNew(t *testing.T) {
	parser := New(`(?P<first>\w+) (?P<second>\w+)`)

	if got, want := parser.FormatInstructions(), parserFormatInstructions; got != want {
		t.Fatalf("parser.FormatInstructions() = %q, want %q", got, want)
	}

	if got, want := parser.Type(), parserType; got != want {
		t.Fatalf("parser.Type() = %q, want %q", got, want)
	}

	matches, err := parser.parse("foo bar baz")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if got, want := len(matches), 2; got != want {
		t.Fatalf(`len(matches) = %d, want %d`, got, want)
	}

	if got, want := matches["first"], "foo"; got != want {
		t.Fatalf(`matches["first"] = %q, want %q`, got, want)
	}

	if got, want := matches["second"], "bar"; got != want {
		t.Fatalf(`matches["second"] = %q, want %q`, got, want)
	}
}
