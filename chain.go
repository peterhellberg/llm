package llm

import "context"

var (
	_ Chain       = &chain{}
	_ ChainHooker = &chain{}
)

// Chain is the interface all chains must implement.
type Chain interface {
	// Call runs the logic of the chain and returns the output. This method should
	// not be called directly. Use rather the chains.Call, chains.Run or chains.Predict
	// functions that handles the memory and other aspects of the chain.
	Call(ctx context.Context, inputs map[string]any, options ...ChainOption) (map[string]any, error)
	// GetMemory gets the memory of the chain.
	Memory() Memory
	// InputKeys returns the input keys the chain expects.
	InputKeys() []string
	// OutputKeys returns the output keys the chain returns.
	OutputKeys() []string
}

const defaultOutputKey = "text"

type chain struct {
	Hooks ChainHooks

	prompter PromptFormatter
	provider Provider
	parser   Parser[any]
	memory   Memory

	outputKey string
}

// NewChain chain with a LLM provider and a prompt.
func NewChain(provider Provider, prompter PromptFormatter, opts ...ChainOption) Chain {
	opt := &ChainOptions{}

	for _, o := range opts {
		o(opt)
	}

	{
		if opt.Hooks == nil {
			opt.Hooks = EmptyHooks{}
		}

		if opt.Parser == nil {
			opt.Parser = EmptyParser{}
		}

		if opt.Memory == nil {
			opt.Memory = EmptyMemory{}
		}

		if opt.OutputKey == "" {
			opt.OutputKey = defaultOutputKey
		}
	}

	return &chain{
		provider: provider,
		prompter: prompter,

		parser:    opt.Parser,
		memory:    opt.Memory,
		outputKey: opt.OutputKey,

		Hooks: opt.Hooks,
	}
}

// Call formats the prompts with the input values, generates using the llm, and parses
// the output from the llm with the output parser. This function should not be called
// directly, use rather the Call or Run function if the prompt only requires one input value.
func (c *chain) Call(ctx context.Context, values map[string]any, options ...ChainOption) (map[string]any, error) {
	prompt, err := c.prompter.FormatPrompt(values)
	if err != nil {
		return nil, err
	}

	result, err := Call(ctx, c.provider, prompt.String(), chainToContentOptions(options...)...)
	if err != nil {
		return nil, err
	}

	finalOutput, err := c.parser.ParseWithPrompt(result, prompt)
	if err != nil {
		return nil, err
	}

	return map[string]any{
		c.outputKey: finalOutput,
	}, nil
}

// Memory returns the memory.
func (c *chain) Memory() Memory {
	return c.memory
}

// ChainHooks returns the hooks for the chain.
func (c *chain) ChainHooks() ChainHooks {
	return c.Hooks
}

// InputKeys returns the expected input keys.
func (c *chain) InputKeys() []string {
	return append([]string{}, c.prompter.InputVariables()...)
}

// OutputKeys returns the output keys the chain will return.
func (c *chain) OutputKeys() []string {
	return []string{
		c.outputKey,
	}
}

// ChainOption is a function that configures ChainOptions.
type ChainOption func(*ChainOptions)

// Options for a chain.
type ChainOptions struct {
	// Model is the model to use in an LLM call.
	Model    string
	modelSet bool

	// MaxTokens is the maximum number of tokens to generate to use in an LLM call.
	MaxTokens    int
	maxTokensSet bool

	// Temperature is the temperature for sampling to use in an LLM call, between 0 and 1.
	Temperature    float64
	temperatureSet bool

	// StopWords is a list of words to stop on to use in an LLM call.
	StopWords    []string
	stopWordsSet bool

	// StreamingFunc is a function to be called for each chunk of a streaming response.
	// Return an error to stop streaming early.
	StreamingFunc func(ctx context.Context, chunk []byte) error

	// TopK is the number of tokens to consider for top-k sampling in an LLM call.
	TopK    int
	topkSet bool

	// TopP is the cumulative probability for top-p sampling in an LLM call.
	TopP    float64
	toppSet bool

	// Seed is a seed for deterministic sampling in an LLM call.
	Seed    int
	seedSet bool

	// MinLength is the minimum length of the generated text in an LLM call.
	MinLength    int
	minLengthSet bool

	// MaxLength is the maximum length of the generated text in an LLM call.
	MaxLength    int
	maxLengthSet bool

	// RepetitionPenalty is the repetition penalty for sampling in an LLM call.
	RepetitionPenalty    float64
	repetitionPenaltySet bool

	// OutputKey to use by the Chain.
	OutputKey string

	// Hooks for the Chain
	Hooks ChainHooks

	// Parser to use by the Chain.
	Parser Parser[any]

	// Memory to use by the Chain.
	Memory Memory
}

// ChainWithModel is an option for LLM.Call.
func ChainWithModel(model string) ChainOption {
	return func(o *ChainOptions) {
		o.Model = model
		o.modelSet = true
	}
}

// ChainWithMaxTokens is an option for LLM.Call.
func ChainWithMaxTokens(maxTokens int) ChainOption {
	return func(o *ChainOptions) {
		o.MaxTokens = maxTokens
		o.maxTokensSet = true
	}
}

// ChainWithTemperature is an option for LLM.Call.
func ChainWithTemperature(temperature float64) ChainOption {
	return func(o *ChainOptions) {
		o.Temperature = temperature
		o.temperatureSet = true
	}
}

// ChainWithStreamingFunc is an option for LLM.Call that allows streaming responses.
func ChainWithStreamingFunc(streamingFunc func(ctx context.Context, chunk []byte) error) ChainOption {
	return func(o *ChainOptions) {
		o.StreamingFunc = streamingFunc
	}
}

// ChainWithTopK will add an option to use top-k sampling for LLM.Call.
func ChainWithTopK(topK int) ChainOption {
	return func(o *ChainOptions) {
		o.TopK = topK
		o.topkSet = true
	}
}

// ChainWithTopP will add an option to use top-p sampling for LLM.Call.
func ChainWithTopP(topP float64) ChainOption {
	return func(o *ChainOptions) {
		o.TopP = topP
		o.toppSet = true
	}
}

// ChainWithSeed will add an option to use deterministic sampling for LLM.Call.
func ChainWithSeed(seed int) ChainOption {
	return func(o *ChainOptions) {
		o.Seed = seed
		o.seedSet = true
	}
}

// ChainWithMinLength will add an option to set the minimum length of the generated text for LLM.Call.
func ChainWithMinLength(minLength int) ChainOption {
	return func(o *ChainOptions) {
		o.MinLength = minLength
		o.minLengthSet = true
	}
}

// ChainWithMaxLength will add an option to set the maximum length of the generated text for LLM.Call.
func ChainWithMaxLength(maxLength int) ChainOption {
	return func(o *ChainOptions) {
		o.MaxLength = maxLength
		o.maxLengthSet = true
	}
}

// ChainWithRepetitionPenalty will add an option to set the repetition penalty for sampling.
func ChainWithRepetitionPenalty(repetitionPenalty float64) ChainOption {
	return func(o *ChainOptions) {
		o.RepetitionPenalty = repetitionPenalty
		o.repetitionPenaltySet = true
	}
}

// ChainWithStopWords is an option for setting the stop words for LLM.Call.
func ChainWithStopWords(stopWords []string) ChainOption {
	return func(o *ChainOptions) {
		o.StopWords = stopWords
		o.stopWordsSet = true
	}
}

// ChainWithOutputKey allows setting what output key should be used by the Chain. (Defaults to "text")
func ChainWithOutputKey(outputKey string) ChainOption {
	return func(o *ChainOptions) {
		o.OutputKey = outputKey
	}
}

// ChainWithHooks allows setting custom Hooks.
func ChainWithHooks(hooks ChainHooks) ChainOption {
	return func(o *ChainOptions) {
		o.Hooks = hooks
	}
}

// ChainWithParser allows setting what parser should be used by the Chain.
func ChainWithParser(parser Parser[any]) ChainOption {
	return func(o *ChainOptions) {
		o.Parser = parser
	}
}

// ChainWithMemory allows setting what memory should be used by the Chain.
func ChainWithMemory(memory Memory) ChainOption {
	return func(o *ChainOptions) {
		o.Memory = memory
	}
}

func chainToContentOptions(options ...ChainOption) []ContentOption {
	opts := &ChainOptions{}

	for _, option := range options {
		option(opts)
	}

	if opts.StreamingFunc == nil && opts.Hooks != nil {
		opts.StreamingFunc = func(ctx context.Context, chunk []byte) error {
			opts.Hooks.StreamingFunc(ctx, chunk)

			return nil
		}
	}

	var cos []ContentOption

	if opts.modelSet {
		cos = append(cos, WithModel(opts.Model))
	}

	if opts.maxTokensSet {
		cos = append(cos, WithMaxTokens(opts.MaxTokens))
	}

	if opts.temperatureSet {
		cos = append(cos, WithTemperature(opts.Temperature))
	}

	if opts.stopWordsSet {
		cos = append(cos, WithStopWords(opts.StopWords))
	}

	if opts.topkSet {
		cos = append(cos, WithTopK(opts.TopK))
	}

	if opts.toppSet {
		cos = append(cos, WithTopP(opts.TopP))
	}

	if opts.seedSet {
		cos = append(cos, WithSeed(opts.Seed))
	}

	if opts.minLengthSet {
		cos = append(cos, WithMinLength(opts.MinLength))
	}

	if opts.maxLengthSet {
		cos = append(cos, WithMaxLength(opts.MaxLength))
	}

	if opts.repetitionPenaltySet {
		cos = append(cos, WithRepetitionPenalty(opts.RepetitionPenalty))
	}

	cos = append(cos, WithStreamingFunc(opts.StreamingFunc))

	return cos
}
