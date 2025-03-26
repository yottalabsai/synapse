package service

import (
	"errors"
	"go.uber.org/zap"
	"synapse/common/log"
	"synapse/connector/constants"
	"sync"

	synapseGrpc "github.com/yottalabsai/endorphin/pkg/services/synapse"
)

type StreamManager struct {
	streamMap map[string]*StreamDetail
	mu        sync.RWMutex
}

type StreamDetail struct {
	stream        synapseGrpc.SynapseService_CallServer
	ClientId      string   `json:"clientId"`
	AgentTypes    []string `json:"agentTypes"`
	SupportModels []string `json:"supportModels"`
	Ready         bool     `json:"ready"`
}

var GlobalStreamManager = &StreamManager{
	streamMap: make(map[string]*StreamDetail),
}

func (m *StreamManager) AddStream(clientID string,
	rawAgentTypes []string,
	supportModels []string,
	stream synapseGrpc.SynapseService_CallServer) {
	m.mu.Lock()
	defer m.mu.Unlock()

	agentTypes, err := convertAgentTypes(rawAgentTypes)
	if err != nil {
		log.Log.Errorw("failed to convert agent types", zap.Error(err))
		return
	}

	m.streamMap[clientID] = &StreamDetail{
		stream:        stream,
		ClientId:      clientID,
		AgentTypes:    agentTypes,
		SupportModels: supportModels,
		Ready:         false,
	}
}

func (m *StreamManager) SendMessage(clientID string, msg *synapseGrpc.JsonRpcRequest) error {
	m.mu.RLock()
	streamDetail, ok := m.streamMap[clientID]
	m.mu.RUnlock()

	if !ok {
		return nil
	}
	return streamDetail.stream.Send(msg)
}

// GetStreams returns a map of ClientId streamMap
func (m *StreamManager) GetStreams() map[string]*StreamDetail {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.streamMap
}

// RemoveStream remove a stream from the manager
func (m *StreamManager) RemoveStream(clientID string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.streamMap, clientID)
}

// convertAgentTypes converts a list of agent types to a list of connector.AgentType
func convertAgentTypes(agentTypes []string) ([]string, error) {
	var converted []string
	for _, agentType := range agentTypes {
		// Add error handling logic if needed
		constants.IsValidAgentType(agentType)
		if agentType == "" {
			return nil, errors.New("invalid agent type")
		}
		converted = append(converted, agentType)
	}
	return converted, nil
}
