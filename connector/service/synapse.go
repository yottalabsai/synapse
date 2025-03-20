package service

import (
	"encoding/json"
	"fmt"
	"github.com/bwmarrin/snowflake"
	synapseGrpc "github.com/yottalabsai/endorphin/pkg/services/synapse"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"io"
	"synapse/common"
	"synapse/common/log"
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
			log.Log.Fatalw("Failed to create snowflake node", zap.Error(err))
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
	// Get metadata from the incoming context
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
			log.Log.Errorw("failed to receive", zap.Error(err))
			return err
		}

		// Handle the received message synchronously
		if err := handleMessage(stream, msg); err != nil {
			log.Log.Errorw("failed to handle message", zap.String("clientId", clientId), zap.Error(err))
			return err
		}
	}

}

func handleMessage(stream synapseGrpc.SynapseService_CallServer, msg *synapseGrpc.Message) error {
	// log.Log.Info("received stream message", zap.Any("message", msg))

	payload := msg.GetText()
	result, err := parseMessage[Example](payload)
	if err != nil {
		return err
	}

	switch v := any(result).(type) {
	case Ping:
		log.Log.Infow("Ping", zap.String("clientId", v.Name))
		// 	checkAgentHealth(msg.ClientId, msg.ModelType, stream)
		return stream.Send(nil)
	case Inference:
		log.Log.Info("Inference")

	case TextToImage:
		log.Log.Info("TextToImage")

	default:
		log.Log.Infow("UnknownResponse")
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

func parseMessage[T any](text string) (*T, error) {
	var result T
	err := json.Unmarshal([]byte(text), &result)
	if err != nil {
		log.Log.Errorw("parseMessage error", zap.Any("text", text), zap.Error(err))
		return nil, err
	}
	return &result, nil
}

type Example struct {
	Name  string `json:"name"`
	Value int    `json:"value"`
}

type Ping struct {
	Name  string `json:"name"`
	Value int    `json:"value"`
}

type Inference struct {
}

type TextToImage struct {
}
