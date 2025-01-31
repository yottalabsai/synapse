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
	"synapse/common"
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
	if !ok || len(md["authorization"]) == 0 {
		return fmt.Errorf("missing authorization")
	}

	// check client_id is required
	clientIds := md.Get("clientid")

	if len(clientIds) == 0 {
		return status.Error(codes.InvalidArgument, "clientId is required")
	}
	clientId := clientIds[0]

	modelTypes := md.Get("modeltype")
	var modelType common.ModelType
	if len(modelTypes) == 0 {
		modelType = common.Inference
	} else {
		modelType = common.ModelType(modelTypes[0])
	}

	// add stream to manager
	GlobalStreamManager.AddStream(clientId, modelType, stream)
	defer GlobalStreamManager.RemoveStream(clientId)

	for {
		// Receive a message from the stream
		msg, err := stream.Recv()
		if err == io.EOF {
			// End of stream
			return nil
		}
		if err != nil {
			log.Log.Error("failed to receive", zap.Error(err))
			return err
		}

		// Handle the received message synchronously
		if err := handleMessage(stream, msg); err != nil {
			log.Log.Error("failed to handle message", zap.String("clientId", clientId), zap.Error(err))
			return err
		}
	}

}

func handleMessage(stream synapseGrpc.SynapseService_CallServer, msg *synapseGrpc.YottaLabsStream) error {
	// log.Log.Info("received stream message", zap.Any("message", msg))

	switch payload := msg.GetPayload().(type) {
	case *synapseGrpc.YottaLabsStream_Ping:
		// log.Log.Info("Ping", zap.String("clientId", msg.ClientId), zap.String("messageId", msg.MessageId))
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
		checkAgentHealth(msg.ClientId, msg.ModelType, stream)
		return stream.Send(resp)

	case *synapseGrpc.YottaLabsStream_RunModelResult:
		log.Log.Info("RunModelResponse", zap.String("clientId", msg.ClientId), zap.String("messageId", msg.MessageId))
		streamDetail := GlobalStreamManager.GetStreams()[msg.ClientId]
		if streamDetail != nil && !streamDetail.Ready {
			streamDetail.Ready = true
			streamDetail.Model = payload.RunModelResult.Model
			break
		}

	case *synapseGrpc.YottaLabsStream_InferenceResult:
		//log.Log.Debug("InferenceResponse", zap.String("clientId", msg.ClientId),
		//	zap.String("messageId", msg.MessageId),
		//	zap.String("content", payload.InferenceResult.Content),
		//)

		channel, ok := GlobalChannelManager.GetChannel(msg.MessageId)
		if !ok {
			log.Log.Error("InferenceResponse", zap.Any("messageId", msg.MessageId))
			return nil
		}
		channel.InferenceResultChan <- payload
	case *synapseGrpc.YottaLabsStream_TextToImageResult:
		channel, ok := GlobalChannelManager.GetChannel(msg.MessageId)
		if !ok {
			log.Log.Error("TextToImageResponse", zap.Any("messageId", msg.MessageId))
			return nil
		}
		channel.TextToImageResultChain <- payload

	default:
		log.Log.Info("UnknownResponse", zap.String("clientId", msg.ClientId), zap.String("messageId", msg.MessageId))
	}

	return nil
}

func checkAgentHealth(clientId string, modeTypeStr string, stream synapseGrpc.SynapseService_CallServer) {
	modeType := getModelType(modeTypeStr)
	_, ok := GlobalStreamManager.GetStreams()[clientId]
	if !ok {
		GlobalStreamManager.AddStream(clientId, modeType, stream)
	}
}

func getModelType(modeType string) common.ModelType {
	return common.ModelType(modeType)
}
