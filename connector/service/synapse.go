package service

import (
	"encoding/json"
	"fmt"
	synapseGrpc "github.com/yottalabsai/endorphin/pkg/services/synapse"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"io"
	"synapse/common/log"
	"synapse/connector/types"
)

var (
	ClientMap = make(map[string]synapseGrpc.SynapseService_CallServer)
)

// SynapseServer synapse service server
type SynapseServer struct {
	synapseGrpc.UnimplementedSynapseServiceServer
}

// NewSynapseServer create new synapse service server
func NewSynapseServer() *SynapseServer {
	return &SynapseServer{}
}

func (s *SynapseServer) Call(stream synapseGrpc.SynapseService_CallServer) error {
	// Get metadata from the incoming context
	md, ok := metadata.FromIncomingContext(stream.Context())
	if !ok || len(md["authorization"]) == 0 {
		return fmt.Errorf("missing authorization")
	}

	// check client_id is required
	clientIds := md.Get("client_id")
	if len(clientIds) == 0 {
		return status.Error(codes.InvalidArgument, "client_id is required")
	}
	clientId := clientIds[0]
	rawAgentTypes := md.Get("agent_types")
	agentTypes, err := convertAgentTypes(rawAgentTypes)
	if err != nil {
		return status.Error(codes.InvalidArgument, err.Error())
	}

	supportModels := md.Get("support_models")

	// add stream to manager
	GlobalStreamManager.AddStream(clientId, agentTypes, supportModels, stream)
	defer GlobalStreamManager.RemoveStream(clientId)

	for {
		// Receive a message from the stream
		req, err := stream.Recv()
		if err == io.EOF {
			// End of stream
			return nil
		}
		if err != nil {
			log.Log.Errorw("failed to receive", zap.Error(err))
			return err
		}

		// Handle the received message synchronously
		if err := handleMessage(stream, req); err != nil {
			log.Log.Errorw("failed to handle message", zap.String("clientId", clientId), zap.Error(err))
			return err
		}
	}

}

func handleMessage(stream synapseGrpc.SynapseService_CallServer, req *synapseGrpc.JsonRpcRequest) error {
	// log.Log.Info("received stream message", zap.Any("message", req))
	messageId, err := getMessageId(req)
	if err != nil {
		return err
	}

	switch req.GetMethod() {
	case types.Error:
		log.Log.Errorw("Error", zap.String("messageId", messageId))

	case types.Ping:
		log.Log.Infow("Ping", zap.String("messageId", messageId))
		// 	checkAgentHealth(req.ClientId, req.ModelType, stream)
		result, err := parseMessage[types.Common](req)
		if err != nil {
			log.Log.Errorw("failed to parse message", zap.String("messageId", messageId), zap.Error(err))
		}

		log.Log.Infow("Successfully parsed message", zap.String("messageId", messageId), zap.Any("result", result))

		return stream.Send(nil)
	case types.Inference:
		log.Log.Info("Inference")

	case types.TextToImage:
		log.Log.Info("TextToImage")

	default:
		log.Log.Infow("UnknownResponse")
	}

	return nil
}

func checkAgentHealth(clientId string, modeTypeStr string, stream synapseGrpc.SynapseService_CallServer) {
	//modeType := getModelType(modeTypeStr)
	//_, ok := GlobalStreamManager.GetStreams()[clientId]
	//if !ok {
	//	GlobalStreamManager.AddStream(clientId, modeType, stream)
	//}
}

func parseMessage[T any](req *synapseGrpc.JsonRpcRequest) (*T, error) {
	var result T
	text := req.Params.GetStringValue()
	err := json.Unmarshal([]byte(text), &result)
	if err != nil {
		log.Log.Errorw("parseMessage error", zap.Any("text", text), zap.Error(err))
		return nil, err
	}
	return &result, nil
}

func getMessageId(req *synapseGrpc.JsonRpcRequest) (string, error) {
	var msgId string
	var err error
	switch id := req.GetId().(type) {
	case *synapseGrpc.JsonRpcRequest_StringId:
		msgId = id.StringId
	case *synapseGrpc.JsonRpcRequest_NumberId:
		msgId = fmt.Sprintf("%d", id.NumberId)
	default:
		err = fmt.Errorf("unknown id type")
	}
	return msgId, err
}
