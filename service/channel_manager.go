package service

import (
	synapseGrpc "github.com/yottalabsai/endorphin/pkg/services/synapse"
	"sync"
)

// InferenceChannel define the response channel
type InferenceChannel struct {
	ResultChan chan *synapseGrpc.YottaLabsStream_InferenceResult
	ErrorChan  chan error
}

// ChannelManager manage all request response channels
type ChannelManager struct {
	sync.RWMutex
	channels map[string]*InferenceChannel
}

// GlobalChannelManager is a global request manager
var GlobalChannelManager = &ChannelManager{
	channels: make(map[string]*InferenceChannel),
}

// CreateChannel create a new response channel for a new request
func (rm *ChannelManager) CreateChannel(requestID string) *InferenceChannel {
	rm.Lock()
	defer rm.Unlock()

	ch := &InferenceChannel{
		ResultChan: make(chan *synapseGrpc.YottaLabsStream_InferenceResult, 10), // 缓冲区大小可调整
		ErrorChan:  make(chan error, 1),
	}

	rm.channels[requestID] = ch
	return ch
}

// GetChannel get the response channel for a request
func (rm *ChannelManager) GetChannel(requestID string) (*InferenceChannel, bool) {
	rm.RLock()
	defer rm.RUnlock()

	ch, exists := rm.channels[requestID]
	return ch, exists
}

// RemoveChannel remove and close the response channel
func (rm *ChannelManager) RemoveChannel(requestID string) {
	rm.Lock()
	defer rm.Unlock()

	if ch, exists := rm.channels[requestID]; exists {
		close(ch.ResultChan)
		close(ch.ErrorChan)
		delete(rm.channels, requestID)
	}
}
