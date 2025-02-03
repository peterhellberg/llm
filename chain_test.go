package llm_test

import (
	"testing"

	"github.com/peterhellberg/llm"
	"github.com/peterhellberg/llm/mock"
)

func TestNewChain(t *testing.T) {
	provider := mock.Provider{}
	template := llm.Template{}

	llm.NewChain(provider, template)
}
