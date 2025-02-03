package llm

import (
	"context"
	"encoding/base64"
)

// TextPart creates TextContent from a given string.
func TextPart(s string) TextContent {
	return TextContent{Text: s}
}

// BinaryPart creates a new BinaryContent from the given MIME type (e.g.
// "image/png" and binary data).
func BinaryPart(mime string, data []byte) BinaryContent {
	return BinaryContent{
		MIMEType: mime,
		Data:     data,
	}
}

// ImageURLPart creates a new ImageURLContent from the given URL.
func ImageURLPart(url string) ImageURLContent {
	return ImageURLContent{
		URL: url,
	}
}

// ImageURLWithDetailPart creates a new ImageURLContent from the given URL and detail.
func ImageURLWithDetailPart(url string, detail string) ImageURLContent {
	return ImageURLContent{
		URL:    url,
		Detail: detail,
	}
}

// ContentPart is an interface all parts of content have to implement.
type ContentPart interface {
	isPart()
}

// TextContent is content with some text.
type TextContent struct {
	Text string
}

func (tc TextContent) String() string {
	return tc.Text
}

func (TextContent) isPart() {}

// ImageURLContent is content with an URL pointing to an image.
type ImageURLContent struct {
	URL    string `json:"url"`
	Detail string `json:"detail,omitempty"` // Detail is the detail of the image, e.g. "low", "high".
}

func (iuc ImageURLContent) String() string {
	return iuc.URL
}

func (ImageURLContent) isPart() {}

// BinaryContent is content holding some binary data with a MIME type.
type BinaryContent struct {
	MIMEType string
	Data     []byte
}

func (bc BinaryContent) String() string {
	base64Encoded := base64.StdEncoding.EncodeToString(bc.Data)
	return "data:" + bc.MIMEType + ";base64," + base64Encoded
}

func (BinaryContent) isPart() {}

// FunctionCall is the name and arguments of a function call.
type FunctionCall struct {
	// The name of the function to call.
	Name string `json:"name"`
	// The arguments to pass to the function, as a JSON string.
	Arguments string `json:"arguments"`
}

// ToolCall is a call to a tool (as requested by the model) that should be executed.
type ToolCall struct {
	// ID is the unique identifier of the tool call.
	ID string `json:"id"`
	// Type is the type of the tool call. Typically, this would be "function".
	Type string `json:"type"`
	// FunctionCall is the function call to be executed.
	FunctionCall *FunctionCall `json:"function,omitempty"`
}

func (ToolCall) isPart() {}

// ToolCallResponse is the response returned by a tool call.
type ToolCallResponse struct {
	// ToolCallID is the ID of the tool call this response is for.
	ToolCallID string `json:"tool_call_id"`
	// Name is the name of the tool that was called.
	Name string `json:"name"`
	// Content is the textual content of the response.
	Content string `json:"content"`
}

func (ToolCallResponse) isPart() {}

// ContentResponse is the response returned by a GenerateContent call.
// It can potentially return multiple content choices.
type ContentResponse struct {
	Choices []*ContentChoice
}

// ContentChoice is one of the response choices returned by GenerateContent
// calls.
type ContentChoice struct {
	// Content is the textual content of a response
	Content string

	// StopReason is the reason the model stopped generating output.
	StopReason string

	// GenerationInfo is arbitrary information the model adds to the response.
	GenerationInfo map[string]any

	// FuncCall is non-nil when the model asks to invoke a function/tool.
	// If a model invokes more than one function/tool, this field will only
	// contain the first one.
	FuncCall *FunctionCall

	// ToolCalls is a list of tool calls the model asks to invoke.
	ToolCalls []ToolCall
}

// ContentOption is a function that configures Options.
type ContentOption func(*ContentOptions)

// ContentOptions is a set of options for calling models.
// Not all models support all options.
type ContentOptions struct {
	// Model is the model to use.
	Model string `json:"model"`
	// CandidateCount is the number of response candidates to generate.
	CandidateCount int `json:"candidate_count"`
	// MaxTokens is the maximum number of tokens to generate.
	MaxTokens int `json:"max_tokens"`
	// Temperature is the temperature for sampling, between 0 and 1.
	Temperature float64 `json:"temperature"`
	// StopWords is a list of words to stop on.
	StopWords []string `json:"stop_words"`
	// StreamingFunc is a function to be called for each chunk of a streaming response.
	// Return an error to stop streaming early.
	StreamingFunc func(ctx context.Context, chunk []byte) error `json:"-"`
	// TopK is the number of tokens to consider for top-k sampling.
	TopK int `json:"top_k"`
	// TopP is the cumulative probability for top-p sampling.
	TopP float64 `json:"top_p"`
	// Seed is a seed for deterministic sampling.
	Seed int `json:"seed"`
	// MinLength is the minimum length of the generated text.
	MinLength int `json:"min_length"`
	// MaxLength is the maximum length of the generated text.
	MaxLength int `json:"max_length"`
	// N is how many chat completion choices to generate for each input message.
	N int `json:"n"`
	// RepetitionPenalty is the repetition penalty for sampling.
	RepetitionPenalty float64 `json:"repetition_penalty"`
	// FrequencyPenalty is the frequency penalty for sampling.
	FrequencyPenalty float64 `json:"frequency_penalty"`
	// PresencePenalty is the presence penalty for sampling.
	PresencePenalty float64 `json:"presence_penalty"`

	// JSONMode is a flag to enable JSON mode.
	JSONMode bool `json:"json"`

	// Tools is a list of tools to use. Each tool can be a specific tool or a function.
	Tools []Tool `json:"tools,omitempty"`
	// ToolChoice is the choice of tool to use, it can either be "none", "auto" (the default behavior),
	// or a specific tool as described in the ToolChoice type.
	ToolChoice any `json:"tool_choice"`

	// Metadata is a map of metadata to include in the request.
	// The meaning of this field is specific to the backend in use.
	Metadata map[string]any `json:"metadata,omitempty"`
}

// Tool is a tool that can be used by the model.
type Tool struct {
	// Type is the type of the tool.
	Type string `json:"type"`
	// Function is the function to call.
	Function *FunctionDefinition `json:"function,omitempty"`
}

// FunctionDefinition is a definition of a function that can be called by the model.
type FunctionDefinition struct {
	// Name is the name of the function.
	Name string `json:"name"`
	// Description is a description of the function.
	Description string `json:"description"`
	// Parameters is a list of parameters for the function.
	Parameters any `json:"parameters,omitempty"`
}

// ToolChoice is a specific tool to use.
type ToolChoice struct {
	// Type is the type of the tool.
	Type string `json:"type"`
	// Function is the function to call (if the tool is a function).
	Function *FunctionReference `json:"function,omitempty"`
}

// FunctionReference is a reference to a function.
type FunctionReference struct {
	// Name is the name of the function.
	Name string `json:"name"`
}

// FunctionCallBehavior is the behavior to use when calling functions.
type FunctionCallBehavior string

const (
	// FunctionCallBehaviorNone will not call any functions.
	FunctionCallBehaviorNone FunctionCallBehavior = "none"
	// FunctionCallBehaviorAuto will call functions automatically.
	FunctionCallBehaviorAuto FunctionCallBehavior = "auto"
)

// WithModel specifies which model name to use.
func WithModel(model string) ContentOption {
	return func(o *ContentOptions) {
		o.Model = model
	}
}

// WithMaxTokens specifies the max number of tokens to generate.
func WithMaxTokens(maxTokens int) ContentOption {
	return func(o *ContentOptions) {
		o.MaxTokens = maxTokens
	}
}

// WithCandidateCount specifies the number of response candidates to generate.
func WithCandidateCount(c int) ContentOption {
	return func(o *ContentOptions) {
		o.CandidateCount = c
	}
}

// WithTemperature specifies the model temperature, a hyperparameter that
// regulates the randomness, or creativity, of the AI's responses.
func WithTemperature(temperature float64) ContentOption {
	return func(o *ContentOptions) {
		o.Temperature = temperature
	}
}

// WithStopWords specifies a list of words to stop generation on.
func WithStopWords(stopWords []string) ContentOption {
	return func(o *ContentOptions) {
		o.StopWords = stopWords
	}
}

// WithOptions specifies options.
func WithOptions(options ContentOptions) ContentOption {
	return func(o *ContentOptions) {
		(*o) = options
	}
}

// WithStreamingFunc specifies the streaming function to use.
func WithStreamingFunc(streamingFunc func(ctx context.Context, chunk []byte) error) ContentOption {
	return func(o *ContentOptions) {
		o.StreamingFunc = streamingFunc
	}
}

// WithTopK will add an option to use top-k sampling.
func WithTopK(topK int) ContentOption {
	return func(o *ContentOptions) {
		o.TopK = topK
	}
}

// WithTopP will add an option to use top-p sampling.
func WithTopP(topP float64) ContentOption {
	return func(o *ContentOptions) {
		o.TopP = topP
	}
}

// WithSeed will add an option to use deterministic sampling.
func WithSeed(seed int) ContentOption {
	return func(o *ContentOptions) {
		o.Seed = seed
	}
}

// WithMinLength will add an option to set the minimum length of the generated text.
func WithMinLength(minLength int) ContentOption {
	return func(o *ContentOptions) {
		o.MinLength = minLength
	}
}

// WithMaxLength will add an option to set the maximum length of the generated text.
func WithMaxLength(maxLength int) ContentOption {
	return func(o *ContentOptions) {
		o.MaxLength = maxLength
	}
}

// WithN will add an option to set how many chat completion choices to generate for each input message.
func WithN(n int) ContentOption {
	return func(o *ContentOptions) {
		o.N = n
	}
}

// WithRepetitionPenalty will add an option to set the repetition penalty for sampling.
func WithRepetitionPenalty(repetitionPenalty float64) ContentOption {
	return func(o *ContentOptions) {
		o.RepetitionPenalty = repetitionPenalty
	}
}

// WithFrequencyPenalty will add an option to set the frequency penalty for sampling.
func WithFrequencyPenalty(frequencyPenalty float64) ContentOption {
	return func(o *ContentOptions) {
		o.FrequencyPenalty = frequencyPenalty
	}
}

// WithPresencePenalty will add an option to set the presence penalty for sampling.
func WithPresencePenalty(presencePenalty float64) ContentOption {
	return func(o *ContentOptions) {
		o.PresencePenalty = presencePenalty
	}
}

// WithToolChoice will add an option to set the choice of tool to use.
// It can either be "none", "auto" (the default behavior), or a specific tool as described in the ToolChoice type.
func WithToolChoice(choice any) ContentOption {
	return func(o *ContentOptions) {
		o.ToolChoice = choice
	}
}

// WithTools will add an option to set the tools to use.
func WithTools(tools []Tool) ContentOption {
	return func(o *ContentOptions) {
		o.Tools = tools
	}
}

// WithJSONMode will add an option to set the response format to JSON.
// This is useful for models that return structured data.
func WithJSONMode() ContentOption {
	return func(o *ContentOptions) {
		o.JSONMode = true
	}
}

// WithMetadata will add an option to set metadata to include in the request.
// The meaning of this field is specific to the backend in use.
func WithMetadata(metadata map[string]any) ContentOption {
	return func(o *ContentOptions) {
		o.Metadata = metadata
	}
}
