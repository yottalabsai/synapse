package service

import (
	"context"
	"fmt"
	"synapse/log"

	"github.com/yottalabsai/endorphin/pkg/services/synapse"
	"go.uber.org/zap"
)

// SynapseServer synapse service server
type SynapseServer struct {
	synapse.UnimplementedSynapseServiceServer
}

// NewSynapseServer create new synapse service server
func NewSynapseServer() *SynapseServer {
	return &SynapseServer{}
}

func (s *SynapseServer) SayHello(ctx context.Context, req *synapse.HelloRequest) (*synapse.HelloResponse, error) {
	log.Log.Info("SayHello", zap.Any("req", req))
	return &synapse.HelloResponse{
		Code:  0,
		Value: fmt.Sprintf("Hello %s!", req.Name),
	}, nil
}

func (s *SynapseServer) ReportAgentStatus(ctx context.Context, req *synapse.AgentStatusRequest) (*synapse.ReportAckResponse, error) {
	log.Log.Info("ReportAgentStatus", zap.Any("req", req))
	return &synapse.ReportAckResponse{
		Success: true,
		Message: "Status reported successfully",
	}, nil
}
