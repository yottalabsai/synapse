package controllers

import (
	"github.com/gin-gonic/gin"
	synapseGrpc "github.com/yottalabsai/endorphin/pkg/services/synapse"
	"go.uber.org/zap"
	"synapse/api/types"
	"synapse/common"
	service2 "synapse/connector/service"
	"synapse/log"
	"synapse/utils"
	"time"
)

type ImageController struct {
	server *service2.SynapseServer
}

func NewImageController(server *service2.SynapseServer) *ImageController {
	return &ImageController{server: server}
}

func (ctl *ImageController) Render(ctx *gin.Context) {
	endpointId := ctx.Param("endpointId")

	if endpointId == "" {
		common.JSON(ctx, common.HttpOk, common.ErrBadArgument)
		return
	}

	req := types.TextToImageRequest{}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		common.JSON(ctx, common.HttpOk, common.ErrBadArgument)
		return
	}
	if err := common.Validate(&req); err != nil {
		httpCode, body := common.GuessError(err, func(error) {
			log.Log.Errorf("validate failed: %v", err)
		})
		common.JSON(ctx, httpCode, body)
		return
	}

	// 执行inference
	ctl.DoRender(ctx, &req)
}

func (ctl *ImageController) DoRender(ctx *gin.Context, req *types.TextToImageRequest) {
	// filter ready client
	requestID := utils.GenerateRequestId()
	flag := false
	for clientID := range service2.GlobalStreamManager.GetStreams() {
		streamDetail := service2.GlobalStreamManager.GetStreams()[clientID]
		log.Log.Infow("[search] clients", zap.Any("clientInfo", streamDetail))
		if streamDetail.Ready && streamDetail.Model == req.Model {
			// create inference request message
			msg := &synapseGrpc.YottaLabsStream{
				MessageId: requestID,
				Timestamp: time.Now().Unix(),
				ClientId:  clientID,

				Payload: &synapseGrpc.YottaLabsStream_TextToImageMessage{
					TextToImageMessage: &synapseGrpc.TextToImageMessage{
						Prompt:            req.Prompt,
						NumInferenceSteps: req.NumInferenceSteps,
						GuidanceScale:     req.GuidanceScale,
						LoraWeight:        req.LoraWeight,
						Seed:              req.Seed,
						Width:             req.Width,
						Height:            req.Height,
						PagScale:          req.PagScale,
					},
				},
			}
			if err := service2.GlobalStreamManager.SendMessage(clientID, msg); err != nil {
				log.Log.Errorw("send message to client failed", zap.Error(err))
			} else {
				flag = true
				break
			}
		}
	}

	if !flag {
		ctx.JSON(common.HttpOk, common.ErrNoReadyClient)
		return
	}

	respChannel := service2.GlobalChannelManager.CreateChannel(requestID)
	defer service2.GlobalChannelManager.RemoveChannel(requestID)

	select {
	case result := <-respChannel.TextToImageResultChain:
		response := &types.TextToImageResponse{
			Created: time.Now().Unix(),
			Data:    result.TextToImageResult.Images,
		}

		ctx.JSON(common.HttpOk, response)
	case <-time.After(30 * time.Second):
		ctx.JSON(common.HttpOk, common.ErrTimeout)
	}

}
