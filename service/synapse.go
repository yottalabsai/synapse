package service

import (
	"context"
	"fmt"
	"synapse/log"

	synapse_grpc "github.com/yottalabsai/endorphin/pkg/services/synapse"
	"go.uber.org/zap"
)

// SynapseServer synapse service server
type SynapseServer struct {
	synapse_grpc.UnimplementedSynapseServiceServer
}

// NewSynapseServer create new synapse service server
func NewSynapseServer() *SynapseServer {
	return &SynapseServer{}
}

func (s *SynapseServer) SayHello(ctx context.Context, req *synapse_grpc.HelloRequest) (*synapse_grpc.HelloResponse, error) {
	log.Log.Info("SayHello", zap.Any("req", req))
	return &synapse_grpc.HelloResponse{
		Code:  0,
		Value: fmt.Sprintf("Hello %s!", req.Name),
	}, nil
}

func (s *SynapseServer) ReportAgentStatus(ctx context.Context, req *synapse_grpc.AgentStatusRequest) (*synapse_grpc.ReportAckResponse, error) {
	log.Log.Info("ReportAgentStatus", zap.Any("req", req))
	return &synapse_grpc.ReportAckResponse{
		Success: true,
		Message: "Status reported successfully",
	}, nil
}
