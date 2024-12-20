package service

import (
	"sync"

	synapse_grpc "github.com/yottalabsai/endorphin/pkg/services/synapse"
)

type StreamManager struct {
	streams map[string]synapse_grpc.SynapseService_CallServer
	mu      sync.RWMutex
}

var GlobalStreamManager = &StreamManager{
	streams: make(map[string]synapse_grpc.SynapseService_CallServer),
}

func (m *StreamManager) AddStream(clientID string, stream synapse_grpc.SynapseService_CallServer) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.streams[clientID] = stream
}

func (m *StreamManager) SendMessage(clientID string, msg *synapse_grpc.StreamMessage) error {
	m.mu.RLock()
	stream, ok := m.streams[clientID]
	m.mu.RUnlock()

	if !ok {
		return nil
	}
	return stream.Send(msg)
}

// GetStreams returns a map of client streams
func (m *StreamManager) GetStreams() map[string]synapse_grpc.SynapseService_CallServer {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.streams
}

// RemoveStream remove a stream from the manager
func (m *StreamManager) RemoveStream(clientID string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.streams, clientID)
}
