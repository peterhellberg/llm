package llm

import "context"

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

	// Hooks for the Chain
	Hooks ChainHooks
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

// ChainWithHooks allows setting custom Hooks.
func ChainWithHooks(hooks ChainHooks) ChainOption {
	return func(o *ChainOptions) {
		o.Hooks = hooks
	}
}

func ChainToContentOptions(options ...ChainOption) []ContentOption {
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
