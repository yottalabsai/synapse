package job

import (
	"context"
	synapseGrpc "github.com/yottalabsai/endorphin/pkg/services/synapse"
	"go.uber.org/zap"
	"synapse/common"
	"synapse/log"
	"synapse/rpc"
	"synapse/service"
	"synapse/utils"
	"time"
)

type InferencePublicModelJob struct {
	ctx    context.Context
	client *rpc.YottaSaaSClient
}

func NewInferencePublicModelJob(ctx context.Context, client *rpc.YottaSaaSClient) *InferencePublicModelJob {
	return &InferencePublicModelJob{ctx: ctx, client: client}
}

func (job *InferencePublicModelJob) Run() {
	// filter not ready client
	modelInfos, err := job.client.FindInferencePublicList(job.ctx)
	if err != nil {
		log.Log.Errorw("get public model list failed", zap.Error(err))
		return
	}

	// get public model
	modelInfoMap := make(map[string]*rpc.ModelInfo)
	for _, modelInfo := range *modelInfos {
		if _, ok := modelInfoMap[modelInfo.ModelName]; !ok {
			modelInfoMap[modelInfo.ModelName] = &modelInfo
		}
	}

	loadedModels := make(map[string]bool)
	// filter ready client
	for clientID := range service.GlobalStreamManager.GetStreams() {
		streamDetail := service.GlobalStreamManager.GetStreams()[clientID]
		log.Log.Infow("[1]connect clients", zap.Any("clientInfo", streamDetail))
		if streamDetail.Ready {
			loadedModels[streamDetail.Model] = true
		}
	}

	loadModels(loadedModels, modelInfoMap)

}

func loadModels(loadedModels map[string]bool, modelInfoMap map[string]*rpc.ModelInfo) {
	for key := range modelInfoMap {
		modelInfo := modelInfoMap[key]
		// if model not loaded, send load model message to client
		log.Log.Infow("[2]public models", zap.Any("modelInfo", modelInfo))
		if _, ok := loadedModels[modelInfo.ModelName]; !ok {
			for clientId := range service.GlobalStreamManager.GetStreams() {
				streamDetail := service.GlobalStreamManager.GetStreams()[clientId]
				log.Log.Infow("[3]running clients", zap.Any("clientInfo", streamDetail))
				if streamDetail.ModelType == modelInfo.ModelType {
					if loadModel(clientId, loadedModels, modelInfo, streamDetail) {
						break
					}
				}
			}
		}

	}
}

func loadModel(clientID string, loadedModels map[string]bool, modelInfo *rpc.ModelInfo, streamDetail *service.StreamDetail) bool {
	if !streamDetail.Ready {
		// create inference request message
		msg := &synapseGrpc.YottaLabsStream{
			MessageId: utils.GenerateRequestId(),
			Timestamp: time.Now().Unix(),
			ClientId:  clientID,
			Payload: &synapseGrpc.YottaLabsStream_RunModelMessage{
				RunModelMessage: &synapseGrpc.RunModelMessage{
					Model: modelInfo.ModelName,
				},
			},
		}
		if err := service.GlobalStreamManager.SendMessage(clientID, msg); err != nil {
			log.Log.Errorw("send load model message to client failed", zap.Error(err))
		} else {
			loadedModels[modelInfo.ModelName] = true
			return true
		}
		return false

	}
	return true
}

func checkModelType(modeTypeStr string, modeType common.ModelType) bool {
	modelType := common.ModelType(modeTypeStr)
	if modelType == modeType {
		return true
	}
	return false
}
