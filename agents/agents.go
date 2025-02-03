package agents

const defaultMaxIterations = 5

// AgentType is a string type representing the type of agent to create.
type AgentType string

const (
	// ZeroShotReactDescription is an AgentType constant that represents
	// the "zeroShotReactDescription" agent type.
	ZeroShotReactDescription AgentType = "zeroShotReactDescription"

	// ConversationalReactDescription is an AgentType constant that represents
	// the "conversationalReactDescription" agent type.
	ConversationalReactDescription AgentType = "conversationalReactDescription"
)
