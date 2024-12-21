package service

import (
	"fmt"
	"github.com/bwmarrin/snowflake"
	synapseGrpc "github.com/yottalabsai/endorphin/pkg/services/synapse"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"io"
	"synapse/log"
	"sync"
)

var (
	node      *snowflake.Node
	nodeOnce  sync.Once
	ClientMap = make(map[string]synapseGrpc.SynapseService_CallServer)
)

func GetSnowflakeNode() *snowflake.Node {
	nodeOnce.Do(func() {
		var err error
		node, err = snowflake.NewNode(1)
		if err != nil {
			log.Log.Fatal("Failed to create snowflake node", zap.Error(err))
		}
	})
	return node
}

// SynapseServer synapse service server
type SynapseServer struct {
	synapseGrpc.UnimplementedSynapseServiceServer
}

// NewSynapseServer create new synapse service server
func NewSynapseServer() *SynapseServer {
	return &SynapseServer{}
}

func (s *SynapseServer) Call(stream synapseGrpc.SynapseService_CallServer) error {
	// 获取元数据
	md, ok := metadata.FromIncomingContext(stream.Context())
	if !ok || len(md["Authorization"]) == 0 {
		return fmt.Errorf("missing authorization")
	}

	// check client_id is required
	clientIds := md.Get("clientId")
	if len(clientIds) == 0 {
		return status.Error(codes.InvalidArgument, "clientId is required")
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
				go func(message *synapseGrpc.YottaLabsStream) {
					if err := handleMessage(stream, message); err != nil {
						log.Log.Error("failed to handle message",
							zap.String("clientId", clientId),
							zap.Error(err))
					}
				}(msg)
			}
		}
	}()

	// keep connection until ClientId disconnects or an error occurs
	<-done
	return nil
}

func handleMessage(stream synapseGrpc.SynapseService_CallServer, msg *synapseGrpc.YottaLabsStream) error {
	log.Log.Info("received stream message", zap.Any("message", msg))

	switch payload := msg.GetPayload().(type) {
	case *synapseGrpc.YottaLabsStream_Ping:
		log.Log.Info("Ping", zap.String("clientId", msg.ClientId), zap.String("messageId", msg.MessageId), zap.Any("payload", payload))
		pong := &synapseGrpc.YottaLabsStream_Pong{
			Pong: &synapseGrpc.PongResult{
				Sequence: payload.Ping.Sequence,
			},
		}
		resp := &synapseGrpc.YottaLabsStream{
			ClientId:  msg.ClientId,
			MessageId: msg.MessageId,
			Payload:   pong,
		}
		return stream.Send(resp)

	case *synapseGrpc.YottaLabsStream_RunModelResult:
		log.Log.Info("RunModelResponse", zap.String("clientId", msg.ClientId), zap.String("messageId", msg.MessageId), zap.Any("payload", payload))
		streamDetail := GlobalStreamManager.GetStreams()[msg.ClientId]
		if streamDetail != nil && !streamDetail.Ready {
			streamDetail.Ready = true
			streamDetail.Model = payload.RunModelResult.Model
			break
		}

	case *synapseGrpc.YottaLabsStream_InferenceResult:
		log.Log.Info("InferenceResponse", zap.String("clientId", msg.ClientId), zap.String("messageId", msg.MessageId), zap.Any("payload", payload))

		inferenceId := payload.InferenceResult.Content
		channel, ok := GlobalChannelManager.GetChannel(inferenceId)
		if !ok {
			log.Log.Error("InferenceResponse", zap.Any("inferenceId", inferenceId), zap.Any("payload", payload))
			return nil
		}
		channel.ResultChan <- payload

	default:
		log.Log.Info("RunModelResponse", zap.String("clientId", msg.ClientId), zap.String("messageId", msg.MessageId), zap.Any("payload", payload))
	}

	return nil
}
