package service

import (
	"synapse/common"
	"sync"

	synapseGrpc "github.com/yottalabsai/endorphin/pkg/services/synapse"
)

type StreamManager struct {
	streamMap map[string]*StreamDetail
	mu        sync.RWMutex
}

type StreamDetail struct {
	stream    synapseGrpc.SynapseService_CallServer
	ClientId  string           `json:"clientId"`
	ModelType common.ModelType `json:"modelType"`
	Model     string           `json:"model"`
	Ready     bool             `json:"ready"`
}

var GlobalStreamManager = &StreamManager{
	streamMap: make(map[string]*StreamDetail),
}

func (m *StreamManager) AddStream(clientID string, modelType common.ModelType, stream synapseGrpc.SynapseService_CallServer) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.streamMap[clientID] = &StreamDetail{
		stream:    stream,
		ClientId:  clientID,
		ModelType: modelType,
		Model:     "",
		Ready:     false,
	}
}

func (m *StreamManager) SendMessage(clientID string, msg *synapseGrpc.YottaLabsStream) error {
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
