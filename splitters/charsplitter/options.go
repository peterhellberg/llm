package charsplitter

import "unicode/utf8"

// Options is a struct that contains options for a charsplitter.
type Options struct {
	ChunkSize     int
	ChunkOverlap  int
	Separators    []string
	KeepSeparator bool
	LenFunc       func(string) int
}

const (
	defaultTokenChunkSize    = 512
	defaultTokenChunkOverlap = 100
)

// defaultOptions returns the default options for charsplitter.
func defaultOptions() Options {
	return Options{
		ChunkSize:     defaultTokenChunkSize,
		ChunkOverlap:  defaultTokenChunkOverlap,
		Separators:    []string{"\n\n", "\n", " ", ""},
		KeepSeparator: false,
		LenFunc:       utf8.RuneCountInString,
	}
}

// Option is a function that can be used to set options for a charsplitter.
type Option func(*Options)

// WithChunkSize sets the chunk size for a charsplitter.
func WithChunkSize(chunkSize int) Option {
	return func(o *Options) {
		o.ChunkSize = chunkSize
	}
}

// WithChunkOverlap sets the chunk overlap for a charsplitter.
func WithChunkOverlap(chunkOverlap int) Option {
	return func(o *Options) {
		o.ChunkOverlap = chunkOverlap
	}
}

// WithSeparators sets the separators for a charsplitter.
func WithSeparators(separators ...string) Option {
	return func(o *Options) {
		o.Separators = separators
	}
}

// WithLenFunc sets the lenfunc for a charsplitter.
func WithLenFunc(lenFunc func(string) int) Option {
	return func(o *Options) {
		o.LenFunc = lenFunc
	}
}

// WithKeepSeparator sets whether the separators should be kept in the resulting
// split text or not. When it is set to True, the separators are included in the
// resulting split text. When it is set to False, the separators are not included
// in the resulting split text. The purpose of having this parameter is to provide
// flexibility in how text splitting is handled. Default to False if not specified.
func WithKeepSeparator(keepSeparator bool) Option {
	return func(o *Options) {
		o.KeepSeparator = keepSeparator
	}
}
