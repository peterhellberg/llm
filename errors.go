package llm

import (
	"errors"
	"fmt"
)

var (
	// ErrInvalidInputValues is returned if the input values to a chain is invalid.
	ErrInvalidInputValues = errors.New("invalid input values")
	// ErrMissingInputValues is returned when some expected input values keys to a chain is missing.
	ErrMissingInputValues = errors.New("missing key in input values")
	// ErrInputValuesWrongType is returned if an input value to a chain is of wrong type.
	ErrInputValuesWrongType = errors.New("input key is of wrong type")
	// ErrMissingMemoryKeyValues is returned when some expected input values keys to a chain is missing.
	ErrMissingMemoryKeyValues = errors.New("missing memory key in input values")
	// ErrMemoryValuesWrongType is returned if the memory value to a chain is of wrong type.
	ErrMemoryValuesWrongType = errors.New("memory key is of wrong type")
	// ErrInvalidOutputValues is returned when expected output keys to a chain does
	// not match the actual keys in the return output values map.
	ErrInvalidOutputValues = errors.New("missing key in output values")
	// ErrMultipleInputsInRun is returned in the run function if the chain expects more then one input values.
	ErrMultipleInputsInRun = errors.New("run not supported in chain with more then one expected input")
	// ErrMultipleOutputsInRun is returned in the run function if the chain expects more then one output values.
	ErrMultipleOutputsInRun = errors.New("run not supported in chain with more then one expected output")
	// ErrWrongOutputTypeInRun is returned in the run function if the chain returns a value that is not a string.
	ErrWrongOutputTypeInRun = errors.New("run not supported in chain that returns value that is not string")
	// ErrOutputNotStringInPredict is returned if a chain does not return a string in the predict function.
	ErrOutputNotStringInPredict = errors.New("predict is not supported with a chain that does not return a string")
	// ErrMultipleOutputsInPredict is returned if a chain has multiple return values in predict.
	ErrMultipleOutputsInPredict = errors.New("predict is not supported with a chain that returns multiple values")
	// ErrChainInitialization is returned if a chain is not initialized appropriately.
	ErrChainInitialization = errors.New("error initializing chain")
	// ErrMismatchMetadatasAndText is returned when the number of texts and metadatas
	// given to CreateDocuments does not match. The function will not error if the
	// length of the metadatas slice is zero.
	ErrMismatchMetadatasAndText = errors.New("number of texts and metadatas does not match")
	// ErrUnexpectedChatMessageType is returned when a chat message is of an unexpected type.
	ErrUnexpectedChatMessageType = errors.New("unexpected chat message type")
	// ErrInputVariableReserved is returned when there is a conflict with a reserved variable name.
	ErrInputVariableReserved = errors.New("conflict with reserved variable name")
	// ErrInvalidPartialVariableType is returned when the partial variable is not a string or a function.
	ErrInvalidPartialVariableType = errors.New("invalid partial variable type")
	// ErrNeedChatMessageList is returned when the variable is not a list of chat messages.
	ErrNeedChatMessageList = errors.New("variable should be a list of chat messages")
	// ErrInvalidTemplateFormat is the error when the template format is invalid and not supported.
	ErrInvalidTemplateFormat = errors.New("invalid template format")
	// ErrEmptyResponseFromModel is returned when there was an empty response from the model.
	ErrEmptyResponseFromProvider = fmt.Errorf("empty response from model")
)
