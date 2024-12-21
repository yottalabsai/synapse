package service

import (
	"sync"

	synapseGrpc "github.com/yottalabsai/endorphin/pkg/services/synapse"
)

type StreamManager struct {
	streamMap map[string]*StreamDetail
	mu        sync.RWMutex
}

type StreamDetail struct {
	stream   synapseGrpc.SynapseService_CallServer
	ClientId string
	Model    string
	Ready    bool
}

var GlobalStreamManager = &StreamManager{
	streamMap: make(map[string]*StreamDetail),
}

func (m *StreamManager) AddStream(clientID string, stream synapseGrpc.SynapseService_CallServer) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.streamMap[clientID] = &StreamDetail{
		stream:   stream,
		ClientId: clientID,
		Model:    "",
		Ready:    false,
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
