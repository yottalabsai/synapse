package service

import (
	"encoding/json"
	"fmt"
	pb "github.com/yottalabsai/endorphin/pkg/services/synapse"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/structpb"
	"io"
	"synapse/common/log"
	"synapse/connector/types"
)

var (
	ClientMap = make(map[string]pb.SynapseService_CallServer)
)

// SynapseServer synapse service server
type SynapseServer struct {
	pb.UnimplementedSynapseServiceServer
}

// NewSynapseServer create new synapse service server
func NewSynapseServer() *SynapseServer {
	return &SynapseServer{}
}

func (s *SynapseServer) Call(stream pb.SynapseService_CallServer) error {
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

func handleMessage(stream pb.SynapseService_CallServer, req *pb.JsonRpcRequest) error {
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
		result, err := parseObject[types.Common](req)
		if err != nil {
			log.Log.Errorw("failed to parse message", zap.String("messageId", messageId), zap.Error(err))
		}

		log.Log.Infow("Successfully parsed message", zap.String("messageId", messageId), zap.Any("result", result))

		pongRequest := &types.PongRequest{
			Common: types.Common{
				ClientId:    "123",
				Timestamp:   0, // Set appropriate timestamp
				MessageId:   messageId,
				MessageType: types.Pong,
			},
		}

		str, err := toJSON[types.PongRequest](pongRequest)
		if err != nil {
			log.Log.Errorw("failed to parse message", zap.String("messageId", messageId))
			return err
		}

		request := &pb.JsonRpcRequest{
			Jsonrpc: "2.0",
			Id:      req.GetId(),
			Method:  types.Pong,
			Params:  structpb.NewStringValue(str),
		}

		return stream.Send(request)
	case types.Inference:
		log.Log.Info("Inference")

	case types.TextToImage:
		log.Log.Info("TextToImage")

	default:
		log.Log.Infow("UnknownResponse")
	}

	return nil
}

func checkAgentHealth(clientId string, modeTypeStr string, stream pb.SynapseService_CallServer) {
	//modeType := getModelType(modeTypeStr)
	//_, ok := GlobalStreamManager.GetStreams()[clientId]
	//if !ok {
	//	GlobalStreamManager.AddStream(clientId, modeType, stream)
	//}
}

func parseObject[T any](req *pb.JsonRpcRequest) (*T, error) {
	var result T
	text := req.Params.GetStringValue()
	err := json.Unmarshal([]byte(text), &result)
	if err != nil {
		log.Log.Errorw("parseObject error", zap.Any("text", text), zap.Error(err))
		return nil, err
	}
	return &result, nil
}

func toJSON[T any](data *T) (string, error) {
	result, err := json.Marshal(data)
	if err != nil {
		log.Log.Errorw("toJSON error", zap.Any("data", data), zap.Error(err))
		return "", err
	}
	return string(result), nil
}

func getMessageId(req *pb.JsonRpcRequest) (string, error) {
	var msgId string
	var err error
	switch id := req.GetId().(type) {
	case *pb.JsonRpcRequest_StringId:
		msgId = id.StringId
	case *pb.JsonRpcRequest_NumberId:
		msgId = fmt.Sprintf("%d", id.NumberId)
	default:
		err = fmt.Errorf("unknown id type")
	}
	return msgId, err
}
