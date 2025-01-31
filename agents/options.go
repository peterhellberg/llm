package agents

import (
	"github.com/peterhellberg/llm"
	"github.com/peterhellberg/llm/memory"
	"github.com/peterhellberg/llm/prompts"
)

type Options struct {
	prompt                  prompts.Template
	memory                  llm.Memory
	hooks                   llm.Hooks
	errorHandler            *llm.ParserErrorHandler
	maxIterations           int
	returnIntermediateSteps bool
	outputKey               string
	promptPrefix            string
	formatInstructions      string
	promptSuffix            string

	// openai
	systemMessage string
	extraMessages []llm.MessageFormatter
}

// Option is a function type that can be used to modify the creation of the agents
// and executors.
type Option func(*Options)

func executorDefaultOptions() Options {
	return Options{
		maxIterations: defaultMaxIterations,
		outputKey:     defaultOutputKey,
		memory:        memory.Empty{},
	}
}

func mrklDefaultOptions() Options {
	return Options{
		promptPrefix:       defaultMrklPrefix,
		formatInstructions: defaultMrklFormatInstructions,
		promptSuffix:       defaultMrklSuffix,
		outputKey:          defaultOutputKey,
	}
}

func conversationalDefaultOptions() Options {
	return Options{
		promptPrefix:       defaultConversationalPrefix,
		formatInstructions: defaultConversationalFormatInstructions,
		promptSuffix:       defaultConversationalSuffix,
		outputKey:          defaultOutputKey,
	}
}

func openAIFunctionsDefaultOptions() Options {
	return Options{
		systemMessage: "You are a helpful AI assistant.",
		outputKey:     defaultOutputKey,
	}
}

func (co Options) getMrklPrompt(tools []llm.AgentTool) prompts.Template {
	if co.prompt.Template != "" {
		return co.prompt
	}

	return createMRKLPrompt(
		tools,
		co.promptPrefix,
		co.formatInstructions,
		co.promptSuffix,
	)
}

func (co Options) getConversationalPrompt(tools []llm.AgentTool) prompts.Template {
	if co.prompt.Template != "" {
		return co.prompt
	}

	return createConversationalPrompt(
		tools,
		co.promptPrefix,
		co.formatInstructions,
		co.promptSuffix,
	)
}

// WithMaxIterations is an option for setting the max number of iterations the executor
// will complete.
func WithMaxIterations(iterations int) Option {
	return func(co *Options) {
		co.maxIterations = iterations
	}
}

// WithOutputKey is an option for setting the output key of the agent.
func WithOutputKey(outputKey string) Option {
	return func(co *Options) {
		co.outputKey = outputKey
	}
}

// WithPromptPrefix is an option for setting the prefix of the prompt used by the agent.
func WithPromptPrefix(prefix string) Option {
	return func(co *Options) {
		co.promptPrefix = prefix
	}
}

// WithPromptFormatInstructions is an option for setting the format instructions of the prompt
// used by the agent.
func WithPromptFormatInstructions(instructions string) Option {
	return func(co *Options) {
		co.formatInstructions = instructions
	}
}

// WithPromptSuffix is an option for setting the suffix of the prompt used by the agent.
func WithPromptSuffix(suffix string) Option {
	return func(co *Options) {
		co.promptSuffix = suffix
	}
}

// WithPrompt is an option for setting the prompt the agent will use.
func WithPrompt(prompt prompts.Template) Option {
	return func(co *Options) {
		co.prompt = prompt
	}
}

// WithReturnIntermediateSteps is an option for making the executor return the intermediate steps
// taken.
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

// WithHooks is an option for setting the hooks for an executor.
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

type OpenAIOption struct{}

func NewOpenAIOption() OpenAIOption {
	return OpenAIOption{}
}

func (o OpenAIOption) WithSystemMessage(msg string) Option {
	return func(co *Options) {
		co.systemMessage = msg
	}
}

func (o OpenAIOption) WithExtraMessages(extraMessages []llm.MessageFormatter) Option {
	return func(co *Options) {
		co.extraMessages = extraMessages
	}
}
