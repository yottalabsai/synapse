package service

import (
	"context"
	"fmt"
	"io"
	"synapse/log"
	"time"

	synapse_grpc "github.com/yottalabsai/endorphin/pkg/services/synapse"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

const (
	schedulerId = "Scheduler-1"
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

func (s *SynapseServer) Call(stream synapse_grpc.SynapseService_CallServer) error {
	// get client metadata
	md, ok := metadata.FromIncomingContext(stream.Context())
	if !ok {
		return status.Error(codes.InvalidArgument, "missing client metadata")
	}

	// check client_id is required
	clientIds := md.Get("client_id")
	if len(clientIds) == 0 {
		return status.Error(codes.InvalidArgument, "client_id is required")
	}
	clientId := clientIds[0]

	// add stream to manager
	GlobalStreamManager.AddStream(clientId, stream)
	defer GlobalStreamManager.RemoveStream(clientId)

	// create a channel to notify goroutine to exit
	done := make(chan struct{})
	defer close(done)

	// start a dedicated goroutine to handle received messages
	go func() {
		for {
			select {
			case <-done:
				return
			default:
				msg, err := stream.Recv()
				if err == io.EOF {
					return
				}
				if err != nil {
					log.Log.Error("failed to receive", zap.Error(err))
					return
				}

				// handle message asynchronously
				go func(message *synapse_grpc.StreamMessage) {
					if err := handleMessage(stream, message); err != nil {
						log.Log.Error("failed to handle message",
							zap.String("clientId", clientId),
							zap.Error(err))
					}
				}(msg)
			}
		}
	}()

	// keep connection until client disconnects or an error occurs
	<-done
	return nil
}

func handleMessage(stream synapse_grpc.SynapseService_CallServer, msg *synapse_grpc.StreamMessage) error {
	log.Log.Info("received stream message", zap.Any("message", msg))

	resp := &synapse_grpc.StreamMessage{
		Base: &synapse_grpc.BaseMessage{
			MessageId: fmt.Sprintf("resp-%d", time.Now().UnixNano()),
			Timestamp: time.Now().Unix(),
			SenderId:  schedulerId,
		},
	}
	resp.Payload = &synapse_grpc.StreamMessage_Heartbeat{
		Heartbeat: &synapse_grpc.HeartbeatRequest{
			ClientId: schedulerId,
		},
	}

	switch payload := msg.GetPayload().(type) {
	case *synapse_grpc.StreamMessage_Heartbeat:
		log.Log.Info("Heartbeat", zap.Any("resp base", resp.GetBase()), zap.Any("payload", payload))
		resp.Metadata = map[string]string{
			"type": "StreamMessage_Heartbeat",
		}

	case *synapse_grpc.StreamMessage_RunModelResponse:
		log.Log.Info("RunModelResponse", zap.Any("resp base", resp.GetBase()), zap.Any("payload", payload))
		resp.Metadata = map[string]string{
			"type": "StreamMessage_RunModelResponse",
		}

	case *synapse_grpc.StreamMessage_InferenceResponse:
		resp.Metadata = map[string]string{
			"type": "StreamMessage_InferenceResponse",
		}
	}

	return stream.Send(resp)
}
