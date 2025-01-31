package llm

import "net/http"

// HTTPRequest http requester interface.
type HTTPRequest interface {
	Do(req *http.Request) (*http.Response, error)
}
