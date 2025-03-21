package constants

const (
	Pod string = "pod"
	Api string = "api"
	Ray string = "ray"
)

var AgentTypes = map[string]struct{}{
	Pod: {},
	Api: {},
	Ray: {},
}

func IsValidAgentType(agentType string) bool {
	_, exists := AgentTypes[agentType]
	return exists
}
