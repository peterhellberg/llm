package agents

import "github.com/peterhellberg/llm"

type Options struct {
	prompt                  llm.Template
	memory                  llm.Memory
	hooks                   llm.Hooks
	errorHandler            *llm.ParserErrorHandler
	maxIterations           int
	returnIntermediateSteps bool
	outputKey               string
	promptPrefix            string
	formatInstructions      string
	promptSuffix            string
}

// Option is a function type that can be used to modify the creation of the agents and executors.
type Option func(*Options)

// WithMaxIterations is an option for setting the max number of iterations the executor will complete.
func WithMaxIterations(maxIterations int) Option {
	return func(co *Options) {
		co.maxIterations = maxIterations
	}
}

// WithOutputKey is an option for setting the output key of the agent.
func WithOutputKey(outputKey string) Option {
	return func(co *Options) {
		co.outputKey = outputKey
	}
}

// WithPromptPrefix is an option for setting the prefix of the prompt used by the agent.
func WithPromptPrefix(promptPrefix string) Option {
	return func(co *Options) {
		co.promptPrefix = promptPrefix
	}
}

// WithPromptFormatInstructions is an option for setting the format instructions of the prompt used by the agent.
func WithPromptFormatInstructions(formatInstructions string) Option {
	return func(co *Options) {
		co.formatInstructions = formatInstructions
	}
}

// WithPromptSuffix is an option for setting the suffix of the prompt used by the agent.
func WithPromptSuffix(suffix string) Option {
	return func(co *Options) {
		co.promptSuffix = suffix
	}
}

// WithPrompt is an option for setting the prompt the agent will use.
func WithPrompt(prompt llm.Template) Option {
	return func(co *Options) {
		co.prompt = prompt
	}
}

// WithReturnIntermediateSteps is an option for making the executor return the intermediate steps taken.
func WithReturnIntermediateSteps() Option {
	return func(co *Options) {
		co.returnIntermediateSteps = true
	}
}

// WithMemory is an option for setting the memory of the executor.
func WithMemory(m llm.Memory) Option {
	return func(co *Options) {
		co.memory = m
	}
}

// WithHooks is an option for setting a callback handler to an executor.
func WithHooks(hooks llm.Hooks) Option {
	return func(co *Options) {
		co.hooks = hooks
	}
}

// WithParserErrorHandler is an option for setting a parser error handler to an executor.
func WithParserErrorHandler(errorHandler *llm.ParserErrorHandler) Option {
	return func(co *Options) {
		co.errorHandler = errorHandler
	}
}
