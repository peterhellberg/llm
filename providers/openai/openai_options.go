package openai

import (
	"github.com/peterhellberg/llm"
	"github.com/peterhellberg/llm/providers/openai/internal/openai"
)

const (
	tokenEnvVarName        = "OPENAI_API_KEY"
	modelEnvVarName        = "OPENAI_MODEL"
	baseURLEnvVarName      = "OPENAI_BASE_URL"
	baseAPIBaseEnvVarName  = "OPENAI_API_BASE"
	organizationEnvVarName = "OPENAI_ORGANIZATION"
)

type APIType openai.APIType

const (
	APITypeOpenAI  APIType = APIType(openai.APITypeOpenAI)
	APITypeAzure           = APIType(openai.APITypeAzure)
	APITypeAzureAD         = APIType(openai.APITypeAzureAD)
)

const DefaultAPIVersion = "2023-05-15"

type options struct {
	token        string
	model        string
	baseURL      string
	organization string
	apiType      APIType
	httpClient   openai.Doer

	responseFormat *ResponseFormat

	// required when APIType is APITypeAzure or APITypeAzureAD
	apiVersion     string
	embeddingModel string

	hooks llm.ProviderHooks
}

// Option is a functional option for the OpenAI client.
type Option func(*options)

// ResponseFormat is the response format for the OpenAI client.
type ResponseFormat = openai.ResponseFormat

// ResponseFormatJSONSchema is the JSON Schema response format in structured output.
type ResponseFormatJSONSchema = openai.ResponseFormatJSONSchema

// ResponseFormatJSONSchemaProperty is the JSON Schema property in structured output.
type ResponseFormatJSONSchemaProperty = openai.ResponseFormatJSONSchemaProperty

// ResponseFormatJSON is the JSON response format.
var ResponseFormatJSON = &ResponseFormat{Type: "json_object"}

// WithToken passes the OpenAI API token to the client. If not set, the token
// is read from the OPENAI_API_KEY environment variable.
func WithToken(token string) Option {
	return func(opts *options) {
		opts.token = token
	}
}

// WithModel passes the OpenAI model to the client. If not set, the model
// is read from the OPENAI_MODEL environment variable.
// Required when ApiType is Azure.
func WithModel(model string) Option {
	return func(opts *options) {
		opts.model = model
	}
}

// WithEmbeddingModel passes the OpenAI model to the client. Required when ApiType is Azure.
func WithEmbeddingModel(embeddingModel string) Option {
	return func(opts *options) {
		opts.embeddingModel = embeddingModel
	}
}

// WithBaseURL passes the OpenAI base url to the client. If not set, the base url
// is read from the OPENAI_BASE_URL environment variable. If still not set in ENV
// VAR OPENAI_BASE_URL, then the default value is https://api.openai.com/v1 is used.
func WithBaseURL(baseURL string) Option {
	return func(opts *options) {
		opts.baseURL = baseURL
	}
}

// WithOrganization passes the OpenAI organization to the client. If not set, the organization is read from the OPENAI_ORGANIZATION.
func WithOrganization(organization string) Option {
	return func(opts *options) {
		opts.organization = organization
	}
}

// WithAPIType passes the api type to the client. If not set, the default value is APITypeOpenAI.
func WithAPIType(apiType APIType) Option {
	return func(opts *options) {
		opts.apiType = apiType
	}
}

// WithAPIVersion passes the api version to the client. If not set, the default value is DefaultAPIVersion.
func WithAPIVersion(apiVersion string) Option {
	return func(opts *options) {
		opts.apiVersion = apiVersion
	}
}

// WithHTTPClient allows setting a custom HTTP client. If not set, the default value is http.DefaultClient.
func WithHTTPClient(client openai.Doer) Option {
	return func(opts *options) {
		opts.httpClient = client
	}
}

// WithHooks allows setting a custom Callback Handler.
func WithHooks(hooks llm.ProviderHooks) Option {
	return func(opts *options) {
		opts.hooks = hooks
	}
}

// WithResponseFormat allows setting a custom response format.
func WithResponseFormat(responseFormat *ResponseFormat) Option {
	return func(opts *options) {
		opts.responseFormat = responseFormat
	}
}
