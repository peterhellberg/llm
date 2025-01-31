package llm

import (
	"strconv"
	"time"
)

// Getenv function for ENV variables
type Getenv func(string) string

// NewEnv creates an Env backed by provided Getenv function
// (such as os.Getenv)
func NewEnv(getenv Getenv) Env {
	return &env{getenv: getenv}
}

// Env is the ENV client interface
type Env interface {
	Bool(key string, fallback bool) bool
	Duration(key string, fallback time.Duration) time.Duration
	Float64(key string, fallback float64) float64
	Int(key string, fallback int) int
	String(key, fallback string) string
}

type env struct {
	getenv func(string) string
}

// Bool returns a bool from the env, or fallback variable
func (e *env) Bool(key string, fallback bool) bool {
	if b, err := strconv.ParseBool(e.getenv(key)); err == nil {
		return b
	}

	return fallback
}

// Duration returns a duration from the env, or fallback variable
func (e *env) Duration(key string, fallback time.Duration) time.Duration {
	if d, err := time.ParseDuration(e.getenv(key)); err == nil {
		return d
	}

	return fallback
}

// Float64 returns a float64 from the env, or a fallback variable
func (e *env) Float64(key string, fallback float64) float64 {
	if f, err := strconv.ParseFloat(e.getenv(key), 64); err == nil {
		return f
	}

	return fallback
}

// Int returns an int from the env, or fallback variable
func (e *env) Int(key string, fallback int) int {
	if i, err := strconv.Atoi(e.getenv(key)); err == nil {
		return i
	}

	return fallback
}

// String returns a string from the env, or fallback variable
func (e *env) String(key, fallback string) string {
	if v := e.getenv(key); v != "" {
		return v
	}

	return fallback
}
