package openai

import "fmt"

var (
	ErrEmptyResponse              = fmt.Errorf("no response")
	ErrMissingToken               = fmt.Errorf("missing the OpenAI API key, set it in the OPENAI_API_KEY environment variable")
	ErrMissingAzureModel          = fmt.Errorf("model needs to be provided when using Azure API")
	ErrMissingAzureEmbeddingModel = fmt.Errorf("embeddings model needs to be provided when using Azure API")
	ErrUnexpectedResponseLength   = fmt.Errorf("unexpected length of response")
)
