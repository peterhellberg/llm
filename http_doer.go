package llm

import "net/http"

// HTTPDoer interface.
type HTTPDoer interface {
	Do(req *http.Request) (*http.Response, error)
}
