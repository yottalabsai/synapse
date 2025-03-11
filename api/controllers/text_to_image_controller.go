package controllers

import (
	"github.com/gin-gonic/gin"
	synapseGrpc "github.com/yottalabsai/endorphin/pkg/services/synapse"
	"go.uber.org/zap"
	"synapse/api/types"
	"synapse/common"
	"synapse/log"
	"synapse/service"
	"synapse/utils"
	"time"
)

type TextToImageController struct {
	server *service.SynapseServer
}

func NewTextToImageController(server *service.SynapseServer) *TextToImageController {
	return &TextToImageController{server: server}
}

func (ctl *TextToImageController) Render(ctx *gin.Context) {
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

func (ctl *TextToImageController) DoRender(ctx *gin.Context, req *types.TextToImageRequest) {
	// filter ready client
	requestID := utils.GenerateRequestId()
	flag := false
	for clientID := range service.GlobalStreamManager.GetStreams() {
		streamDetail := service.GlobalStreamManager.GetStreams()[clientID]
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
			if err := service.GlobalStreamManager.SendMessage(clientID, msg); err != nil {
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

	respChannel := service.GlobalChannelManager.CreateChannel(requestID)
	defer service.GlobalChannelManager.RemoveChannel(requestID)

	select {
	case result := <-respChannel.TextToImageResultChain:
		ctx.JSON(common.HttpOk, result.TextToImageResult.Images)
	case <-time.After(30 * time.Second):
		ctx.JSON(common.HttpOk, common.ErrTimeout)
	}

}
