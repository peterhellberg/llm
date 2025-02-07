package qdrant

import (
	"net/http"
	"net/url"
)

const defaultContentKey = "content"

type Options struct {
	Base                 *url.URL
	HTTP                 *http.Client
	APIKey               string
	ContentKey           string
	CollectionName       string
	CollectionVectorSize int
}

func DefaultOptions() Options {
	return Options{
		ContentKey: defaultContentKey,
	}
}
