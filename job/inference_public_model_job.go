package job

import (
	"context"
	synapseGrpc "github.com/yottalabsai/endorphin/pkg/services/synapse"
	"go.uber.org/zap"
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
		log.Log.Error("get public model list failed", zap.Error(err))
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
		log.Log.Infof("[1]已连接client信息: cilentID: %s, model: %s, ready: %t", streamDetail.ClientId, streamDetail.Model, streamDetail.Ready)
		if streamDetail.Ready {
			loadedModels[streamDetail.Model] = true
		}
	}

	for key := range modelInfoMap {
		modelInfo := modelInfoMap[key]
		// if model not loaded, send load model message to client
		log.Log.Infof("[2]公开model信息: modelID: %s, modelName: %s, ready: %t", modelInfo.ModelID, modelInfo.ModelName, modelInfo.Ready)
		if _, ok := loadedModels[modelInfo.ModelName]; !ok {
			for clientID := range service.GlobalStreamManager.GetStreams() {
				streamDetail := service.GlobalStreamManager.GetStreams()[clientID]
				// filter not ready client
				log.Log.Infof("[3]已连接client信息: cilentID: %s, model: %s, ready: %t", streamDetail.ClientId, streamDetail.Model, streamDetail.Ready)

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
						log.Log.Error("send load model message to client failed", zap.Error(err))
					} else {
						loadedModels[modelInfo.ModelName] = true
					}

					break
				}
			}
		}

	}

}
