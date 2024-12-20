package service

import (
	"sync"

	synapse_grpc "github.com/yottalabsai/endorphin/pkg/services/synapse"
)

// ResponseChannel define the response channel
type ResponseChannel struct {
	ResultChan chan *synapse_grpc.StreamMessage_InferenceResponse
	ErrorChan  chan error
}

// RequestManager manage all request response channels
type RequestManager struct {
	sync.RWMutex
	channels map[string]*ResponseChannel
}

// GlobalRequestManager is a global request manager
var GlobalRequestManager = &RequestManager{
	channels: make(map[string]*ResponseChannel),
}

// CreateChannel create a new response channel for a new request
func (rm *RequestManager) CreateChannel(requestID string) *ResponseChannel {
	rm.Lock()
	defer rm.Unlock()

	ch := &ResponseChannel{
		ResultChan: make(chan *synapse_grpc.StreamMessage_InferenceResponse, 10), // 缓冲区大小可调整
		ErrorChan:  make(chan error, 1),
	}
	rm.channels[requestID] = ch
	return ch
}

// GetChannel get the response channel for a request
func (rm *RequestManager) GetChannel(requestID string) (*ResponseChannel, bool) {
	rm.RLock()
	defer rm.RUnlock()

	ch, exists := rm.channels[requestID]
	return ch, exists
}

// RemoveChannel remove and close the response channel
func (rm *RequestManager) RemoveChannel(requestID string) {
	rm.Lock()
	defer rm.Unlock()

	if ch, exists := rm.channels[requestID]; exists {
		close(ch.ResultChan)
		close(ch.ErrorChan)
		delete(rm.channels, requestID)
	}
}
