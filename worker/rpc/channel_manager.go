package rpc

import (
	synapseGrpc "github.com/yottalabsai/endorphin/pkg/services/synapse"
	"sync"
)

// AgentChannel define the response channel
type AgentChannel struct {
	InferenceResultChan    chan *synapseGrpc.Message
	TextToImageResultChain chan *synapseGrpc.Message
	ErrorChan              chan error
}

// ChannelManager manage all request response channels
type ChannelManager struct {
	sync.RWMutex
	channels map[string]*AgentChannel
}

// GlobalChannelManager is a global request manager
var GlobalChannelManager = &ChannelManager{
	channels: make(map[string]*AgentChannel),
}

// CreateChannel create a new response channel for a new request
func (rm *ChannelManager) CreateChannel(requestID string) *AgentChannel {
	rm.Lock()
	defer rm.Unlock()

	ch := &AgentChannel{
		InferenceResultChan: make(chan *synapseGrpc.Message, 10), // 缓冲区大小可调整
		//TextToImageResultChain: make(chan *synapseGrpc.YottaLabsStream_TextToImageResult, 10),
		ErrorChan: make(chan error, 1),
	}

	rm.channels[requestID] = ch
	return ch
}

// GetChannel get the response channel for a request
func (rm *ChannelManager) GetChannel(requestID string) (*AgentChannel, bool) {
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
		close(ch.InferenceResultChan)
		close(ch.TextToImageResultChain)
		close(ch.ErrorChan)
		delete(rm.channels, requestID)
	}
}
