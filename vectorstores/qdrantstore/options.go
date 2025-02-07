package qdrantstore

import (
	"net/http"
	"net/url"

	"github.com/peterhellberg/llm"
)

// Option is a function that configures an Options.
type Option func(s *Store) error

// WithRawURL returns an Option for setting the Qdrant instance URL.
// Example: 'http://localhost:63333'. Required.
func WithRawURL(rawURL string) Option {
	return func(s *Store) error {
		var err error

		s.options.Base, err = url.Parse(rawURL)

		return err
	}
}

// WithHTTPClient Set custom http client.
func WithHTTPClient(client *http.Client) Option {
	return func(s *Store) error {
		s.options.HTTP = client

		return nil
	}
}

// WithAPIKey returns an Option for setting the API key to authenticate the connection. Optional.
func WithAPIKey(apiKey string) Option {
	return func(s *Store) error {
		s.options.APIKey = apiKey

		return nil
	}
}

// WithCollectionName returns an Option for setting the collection name. Required.
func WithCollectionName(collectionName string) Option {
	return func(s *Store) error {
		s.options.CollectionName = collectionName

		return nil
	}
}

// WithContent returns an Option for setting field name of the document content
// in the Qdrant payload. Optional. Defaults to "content".
func WithContentKey(contentKey string) Option {
	return func(s *Store) error {
		s.options.ContentKey = contentKey

		return nil
	}
}

// WithCollectionVectorSize returns an Option for setting the collection vector size.
func WithCollectionVectorSize(collectionVectorSize int) Option {
	return func(s *Store) error {
		s.options.CollectionVectorSize = collectionVectorSize

		return nil
	}
}

// WithEmbedder returns an Option for setting the embedder to be used when
// adding documents or doing similarity search. Required.
func WithEmbedder(embedder llm.Embedder) Option {
	return func(s *Store) error {
		s.embedder = embedder

		return nil
	}
}
