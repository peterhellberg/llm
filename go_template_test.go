package llm

import "testing"

func TestGoTemplate(t *testing.T) {
	for _, tt := range []struct {
		content   string
		variables []string
		want      string
	}{
		{"{{.foo}} and {{.bar}}", []string{"foo", "ba"}, "a string and 123"},
	} {
		template := GoTemplate(tt.content, tt.variables)

		got, err := template.FormatString(map[string]any{
			"foo": "a string",
			"bar": 123,
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if got != tt.want {
			t.Fatalf("got = %q, want %q", got, tt.want)
		}
	}
}
